package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	_ "modernc.org/sqlite"
)

func NewSqliteEngine(fileSource string, logger zerolog.Logger) (*Engine, error) {
	logger.Info().Msgf("new connection to sqlite: %s", fileSource)
	db, err := sqlx.Connect("sqlite", fileSource)
	if err != nil {
		logger.Error().Err(err).Msg("failed to connect to sqlite")
		return nil, err
	}

	engine := &Engine{DB: db}

	if err := engine.setup(); err != nil {
		logger.Error().Err(err).Msg("failed to setup sqlite engine")
	}

	if err := engine.migrate(logger); err != nil {
		logger.Error().Err(err).Msg("failed to migrate sqlite schema")
		_ = engine.Close()
		return nil, err
	}

	return engine, nil
}

func (e *Engine) setup() error {
	pragmas := []string{
		"PRAGMA journal_mode=WAL",
		"PRAGMA busy_timeout=5000",
		"PRAGMA synchronous=NORMAL",
		"PRAGMA cache_size=1000",
		"PRAGMA foreign_keys=ON",
	}
	for _, pragma := range pragmas {
		if _, err := e.Exec(pragma); err != nil {
			_ = e.Close()
			return fmt.Errorf("failed to set pragma %q: %w", pragma, err)
		}
	}

	// limit connections for SQLite (single writer)
	e.SetMaxOpenConns(1)

	return nil
}
