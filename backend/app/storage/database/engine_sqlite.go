package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	_ "modernc.org/sqlite"
)

const (
	userSchema = `
		CREATE TABLE IF NOT EXISTS user (
			id TEXT PRIMARY KEY NOT NULL,
			email TEXT NOT NULL UNIQUE,
			name TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`
	trackSchema = `
		CREATE TABLE IF NOT EXISTS track (
			id TEXT PRIMARY KEY NOT NULL,
			slug TEXT NOT NULL UNIQUE,
			name TEXT NOT NULL,
			description TEXT NOT NULL,
			owner_id TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`

	workoutSchema = `
		CREATE TABLE IF NOT EXISTS workout (
			id TEXT PRIMARY KEY NOT NULL,
			date TEXT NOT NULL,
			track_id TEXT NOT NULL,
			sections TEXT NOT NULL,
			notes TEXT,
			schema_version INTEGER NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`
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

	if err := engine.createSchema(); err != nil {
		logger.Error().Err(err).Msg("failed to create sqlite schema")
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

func (e *Engine) createSchema() error {
	schemas := map[string]string{
		"user":    userSchema,
		"track":   trackSchema,
		"workout": workoutSchema,
	}

	for table, schema := range schemas {
		if _, err := e.Exec(schema); err != nil {
			return fmt.Errorf("failed to create %s table: %w", table, err)
		}
	}

	return nil
}
