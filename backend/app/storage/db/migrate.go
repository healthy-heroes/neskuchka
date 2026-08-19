package db

import (
	"context"
	"embed"
	"fmt"
	"io/fs"

	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// migrate brings the schema up to date by applying the embedded migrations.
func (e *Engine) migrate(logger zerolog.Logger) error {
	dir, err := fs.Sub(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("failed to open migrations dir: %w", err)
	}

	provider, err := goose.NewProvider(goose.DialectSQLite3, e.DB.DB, dir)
	if err != nil {
		return fmt.Errorf("failed to create migration provider: %w", err)
	}

	results, err := provider.Up(context.Background())
	if err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	if len(results) > 0 {
		logger.Info().Msgf("applied %d migrations", len(results))
	}

	return nil
}
