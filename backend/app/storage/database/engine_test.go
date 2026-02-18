package database

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func setupTestSqliteDB(t *testing.T) *Engine {
	t.Helper()

	engine, err := NewSqliteEngine(":memory:", zerolog.Nop())
	require.NoError(t, err)

	t.Cleanup(func() {
		engine.Close()
	})

	return engine
}
