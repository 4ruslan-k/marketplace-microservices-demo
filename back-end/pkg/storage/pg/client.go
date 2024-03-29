package storage

import (
	"context"
	"database/sql"

	"github.com/rs/zerolog"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

type Config struct {
	DSN string
}

func initializeClient(logger zerolog.Logger, dsn string, verbose bool) *bun.DB {
	connector := pgdriver.NewConnector(pgdriver.WithDSN(dsn))
	sqldb := sql.OpenDB(connector)
	db := bun.NewDB(sqldb, pgdialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(verbose),
	))
	_, err := db.ExecContext(context.Background(), "SELECT 1")
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to database")
	}
	logger.Info().Msgf("Connected to database, addr: %s, database: %s", connector.Config().Addr, connector.Config().Database)

	return db
}

func NewClient(logger zerolog.Logger, config Config) *bun.DB {
	dsn := config.DSN
	return initializeClient(logger, dsn, true)
}

func NewClientWithDSN(logger zerolog.Logger, dsn string, verbose bool) *bun.DB {
	return initializeClient(logger, dsn, verbose)

}
