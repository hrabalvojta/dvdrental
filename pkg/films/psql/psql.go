package psql

import (
	"context"
	"database/sql"

	"github.com/go-kit/kit/log"

	"github.com/hrabalvojta/dvdrental/internal/config"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/migrate"
)

func NewDB(cfg config.Config) *bun.DB {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(cfg.Postgres_dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())

	return db
}

func StartMigration(db *bun.DB, mig *migrate.Migrations, ctx context.Context, logger log.Logger) error {
	migrator := migrate.NewMigrator(db, mig)

	err := migrator.Init(ctx)
	if err != nil {
		return err
	}

	logger.Log("db", "postgres", "state", "migration", "status", "check")
	group, err := migrator.Migrate(ctx)
	if err != nil {
		return err
	}
	if group.ID == 0 {
		logger.Log("db", "postgres", "state", "migration", "status", "no_change")
		return nil
	}
	logger.Log("db", "postgres", "state", "migration", "updated_to", group)
	return nil
}
