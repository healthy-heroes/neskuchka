package testutil

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/healthy-heroes/neskuchka/backend/app/storage/db"
)

// NewEngine creates an in-memory sqlite engine closed when the test ends.
// It sits at the engine level on purpose: a helper handing out a ready
// *datastorage.Storage could not be used by the datastorage tests themselves,
// they are in-package and importing it back would be an import cycle.
func NewEngine(t *testing.T) *db.Engine {
	t.Helper()

	engine, err := db.NewSqliteEngine(":memory:", zerolog.Nop())
	require.NoError(t, err)

	t.Cleanup(func() {
		engine.Close()
	})

	return engine
}
