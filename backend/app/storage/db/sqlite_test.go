package db

import (
	"path/filepath"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestNewSqliteEngine_AppliesMigrations(t *testing.T) {
	engine, err := NewSqliteEngine(":memory:", zerolog.Nop())
	require.NoError(t, err)
	defer engine.Close()

	// goose ведёт версию схемы в собственной таблице
	var version int
	err = engine.Get(&version, "SELECT MAX(version_id) FROM goose_db_version")
	require.NoError(t, err)
	require.GreaterOrEqual(t, version, 1)

	var tables []string
	err = engine.Select(&tables,
		"SELECT name FROM sqlite_master WHERE type='table' AND name IN ('user','track','workout','avatar') ORDER BY name")
	require.NoError(t, err)
	require.Equal(t, []string{"avatar", "track", "user", "workout"}, tables)
}

func TestNewSqliteEngine_ReopenIsIdempotent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "app.db")

	engine, err := NewSqliteEngine(path, zerolog.Nop())
	require.NoError(t, err)
	_, err = engine.Exec("INSERT INTO user(id, email, name) VALUES('u1', 'u1@example.com', 'U1')")
	require.NoError(t, err)
	require.NoError(t, engine.Close())

	// повторное открытие = рестарт приложения: миграции уже применены
	engine, err = NewSqliteEngine(path, zerolog.Nop())
	require.NoError(t, err)
	defer engine.Close()

	var count int
	require.NoError(t, engine.Get(&count, "SELECT COUNT(*) FROM user"))
	require.Equal(t, 1, count)
}

func TestNewSqliteEngine_AdoptsPreGooseDatabase(t *testing.T) {
	path := filepath.Join(t.TempDir(), "app.db")

	// база, созданная старым createSchema: таблицы есть, goose_db_version нет
	raw, err := sqlx.Connect("sqlite", path)
	require.NoError(t, err)
	oldSchema := []string{
		`CREATE TABLE IF NOT EXISTS user (
			id TEXT PRIMARY KEY NOT NULL,
			email TEXT NOT NULL UNIQUE,
			name TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS track (
			id TEXT PRIMARY KEY NOT NULL,
			slug TEXT NOT NULL UNIQUE,
			name TEXT NOT NULL,
			description TEXT NOT NULL,
			owner_id TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS workout (
			id TEXT PRIMARY KEY NOT NULL,
			date TEXT NOT NULL,
			track_id TEXT NOT NULL,
			sections TEXT NOT NULL,
			notes TEXT,
			schema_version INTEGER NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS avatar (
			user_id TEXT PRIMARY KEY NOT NULL,
			mime_type TEXT NOT NULL,
			data BLOB NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
	}
	for _, ddl := range oldSchema {
		_, err = raw.Exec(ddl)
		require.NoError(t, err)
	}
	_, err = raw.Exec("INSERT INTO user(id, email, name) VALUES('u1', 'u1@example.com', 'U1')")
	require.NoError(t, err)
	require.NoError(t, raw.Close())

	engine, err := NewSqliteEngine(path, zerolog.Nop())
	require.NoError(t, err)
	defer engine.Close()

	// миграция 0001 прошла no-op'ом и зафиксировала версию
	var version int
	require.NoError(t, engine.Get(&version, "SELECT MAX(version_id) FROM goose_db_version"))
	require.Equal(t, 1, version)

	// данные пережили adoption
	var count int
	require.NoError(t, engine.Get(&count, "SELECT COUNT(*) FROM user"))
	require.Equal(t, 1, count)
}
