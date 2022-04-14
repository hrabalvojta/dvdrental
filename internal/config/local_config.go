package config

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	// Application
	App_env    string `env:"app_env,required" envDefault:"DEVELOPMENT"`
	Token      string `env:"token"`
	Debug_addr string `env:"debug_addr,required,notEmpty" envDefault:":8080"`
	Http_addr  string `env:"http_addr,required,notEmpty" envDefault:":8081"`
	Channel    string `env:"gokit-films" envDefault:"gokit-films"`

	// DATABASE PSQL
	Postgres_dsn     string `env:"postgres_database,required,notEmpty" envDefault:"postgres://postgres:secret@localhost:5432/dvdrental?sslmode=disable"`
	Postgres_timeout int    `env:"postgres_timeout" envDefault:"1"`
}

func InitConfig() (Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
