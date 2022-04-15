package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/oklog/run"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/uptrace/bun"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	"github.com/hrabalvojta/dvdrental/internal/config"
	"github.com/hrabalvojta/dvdrental/pkg/films/endpoints"
	"github.com/hrabalvojta/dvdrental/pkg/films/psql"
	"github.com/hrabalvojta/dvdrental/pkg/films/psql/migrations"
	"github.com/hrabalvojta/dvdrental/pkg/films/service"
	"github.com/hrabalvojta/dvdrental/pkg/films/transport"
)

type appCtxKey struct{}

func ContextWithApp(ctx context.Context, app *App) context.Context {
	ctx = context.WithValue(ctx, appCtxKey{}, app)
	return ctx
}

type App struct {
	// Application core
	cfg    *config.Config
	ctx    context.Context
	logger log.Logger
	group  run.Group

	// Graceful stop
	stopping uint32
	stopCh   chan struct{}

	// Prometheus
	ints, chars metrics.Counter
	duration    metrics.Histogram

	// Lazy DB init
	dbOnce sync.Once
	db     *bun.DB
}

func NewApp() *App {
	app := &App{
		stopCh: make(chan struct{}),
	}

	app.ctx = ContextWithApp(context.Background(), app)

	return app
}

func (app *App) Context() context.Context {
	return app.ctx
}

func (app *App) Config() *config.Config {
	return app.cfg
}

func main() {

	// Create Application struct which contains all important variables
	app := NewApp()

	// Define our flags. Your service probably won't need to bind listeners for
	// *all* supported transports, or support both Zipkin and LightStep, and so
	// on, but we do it here for demonstration purposes.}
	//fs := flag.NewFlagSet("films", flag.ExitOnError)
	//var (
	//	httpAddr  = fs.String("http-addr", ":8081", "HTTP listen address")
	//	debugAddr = fs.String("debug.addr", ":8080", "Debug and metrics listen address")
	//)
	//fs.Usage = usageFor(fs, os.Args[0]+" [flags]")
	//fs.Parse(os.Args[1:])

	// Create a single logger, which we'll use and give to other components.
	//var logger log.Logger
	{
		app.logger = log.NewLogfmtLogger(os.Stderr)
		app.logger = log.With(app.logger, "ts", log.DefaultTimestampUTC)
		app.logger = log.With(app.logger, "caller", log.DefaultCaller)
	}

	{
		var err error
		app.cfg, err = config.InitConfig()
		if err != nil {
			app.logger.Log("env_config", "debug/env", "during", "Parse", "err", err)
			os.Exit(1)
		}
		app.logger = log.With(app.logger, "channel", app.cfg.Channel)
		app.logger.Log("env_config", "debug/env", "config", "env", "loaded", "success")
	}

	{
		var err error
		for true {
			app.db, err = psql.NewDB(app.cfg)
			if err == nil {
				break
			}
			app.logger.Log("db", "postgres", "state", "migration", "err", err)
			time.Sleep(time.Duration(app.cfg.Postgres_timeout) * time.Second)
		}
		err = psql.StartMigration(app.db, migrations.Migrations, app.ctx, app.logger)
		if err != nil {
			app.logger.Log("db", "postgres", "state", "migration", "err", err)
			os.Exit(1)
		}
	}

	// Create the (sparse) metrics we'll use in the service. They, too, are
	// dependencies that we pass to components that use them.
	{
		// Business-level metrics.
		app.ints = prometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "dvdrental",
			Subsystem: "films",
			Name:      "integers_summed",
			Help:      "Total count of integers summed via the Sum method.",
		}, []string{})
		app.chars = prometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "dvdrental",
			Subsystem: "films",
			Name:      "characters_concatenated",
			Help:      "Total count of characters concatenated via the Concat method.",
		}, []string{})
	}
	{
		// Endpoint-level metrics.
		app.duration = prometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "dvdrental",
			Subsystem: "films",
			Name:      "request_duration_seconds",
			Help:      "Request duration in seconds.",
		}, []string{"method", "success"})
	}
	http.DefaultServeMux.Handle("/metrics", promhttp.Handler())

	// Build the layers of the service "onion" from the inside out. First, the
	// business logic service; then, the set of endpoints that wrap the service;
	// and finally, a series of concrete transport adapters. The adapters, like
	// the HTTP handler or the gRPC server, are the bridge between Go kit and
	// the interfaces that the transports expect. Note that we're not binding
	// them to ports or anything yet; we'll do that next.
	var (
		service     = service.New(app.logger, app.ints, app.chars)
		endpoints   = endpoints.New(service, app.logger, app.duration)
		httpHandler = transport.NewHTTPHandler(endpoints, app.logger)
	)

	// Now we're to the part of the func main where we want to start actually
	// running things, like servers bound to listeners to receive connections.
	//
	// The method is the same for each component: add a new actor to the group
	// struct, which is a combination of 2 anonymous functions: the first
	// function actually runs the component, and the second function should
	// interrupt the first function and cause it to return. It's in these
	// functions that we actually bind the Go kit server/handler structs to the
	// concrete transports and run them.
	//
	// Putting each component into its own block is mostly for aesthetics: it
	// clearly demarcates the scope in which each listener/socket may be used.
	//var g run.Group
	{
		// The debug listener mounts the http.DefaultServeMux, and serves up
		// stuff like the Prometheus metrics route, the Go debug and profiling
		// routes, and so on.
		debugListener, err := net.Listen("tcp", app.cfg.Debug_addr)
		if err != nil {
			app.logger.Log("transport", "debug/HTTP", "during", "Listen", "err", err)
			os.Exit(1)
		}
		app.group.Add(func() error {
			app.logger.Log("transport", "debug/HTTP", "addr", app.cfg.Debug_addr)
			return http.Serve(debugListener, http.DefaultServeMux)
		}, func(error) {
			debugListener.Close()
		})
	}
	{
		// The HTTP listener mounts the Go kit HTTP handler we created.
		httpListener, err := net.Listen("tcp", app.cfg.Http_addr)
		if err != nil {
			app.logger.Log("transport", "HTTP", "during", "Listen", "err", err)
			os.Exit(1)
		}
		app.group.Add(func() error {
			app.logger.Log("transport", "HTTP", "addr", app.cfg.Http_addr)
			return http.Serve(httpListener, httpHandler)
		}, func(error) {
			httpListener.Close()
		})
	}
	{
		// This function just sits and waits for ctrl-C.
		cancelInterrupt := make(chan struct{})
		app.group.Add(func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			select {
			case sig := <-c:
				return fmt.Errorf("received signal %s", sig)
			case <-cancelInterrupt:
				return nil
			}
		}, func(error) {
			close(cancelInterrupt)
		})
	}
	app.logger.Log("exit", app.group.Run())
}

//func usageFor(fs *flag.FlagSet, short string) func() {
//	return func() {
//		fmt.Fprintf(os.Stderr, "USAGE\n")
//		fmt.Fprintf(os.Stderr, "  %s\n", short)
//		fmt.Fprintf(os.Stderr, "\n")
//		fmt.Fprintf(os.Stderr, "FLAGS\n")
//		w := tabwriter.NewWriter(os.Stderr, 0, 8, 2, ' ', 0)
//		fs.VisitAll(func(f *flag.Flag) {
//			fmt.Fprintf(w, "\t-%s\t%s\t%s\n", f.Name, f.DefValue, f.Usage)
//		})
//		w.Flush()
//		fmt.Fprintf(os.Stderr, "\n")
//	}
//}
