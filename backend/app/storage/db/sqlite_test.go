package db

import (
	"testing"

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
