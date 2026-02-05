package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
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
)

func NewSqliteEngine(fileSource string, logger zerolog.Logger) (*Engine, error) {
	logger.Info().Msgf("new connection to sqlite: %s", fileSource)
	db, err := sqlx.Connect("sqlite", fileSource)
	if err != nil {
		logger.Error().Err(err).Msg("failed to connect to sqlite")
		return nil, err
	}

	// todo useful pragma

	engine := &Engine{DB: db}

	if err := engine.createSqliteSchema(); err != nil {
		logger.Error().Err(err).Msg("failed to create sqlite schema")
		return nil, err
	}

	return engine, nil
}

func (e *Engine) createSqliteSchema() error {

	if _, err := e.Exec(userSchema); err != nil {
		return fmt.Errorf("failed to create user table: %w", err)
	}

	return nil
}
