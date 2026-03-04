package avatarstorage

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/healthy-heroes/neskuchka/backend/app/storage/db"
)

func setupTestStorage(t *testing.T) *Storage {
	t.Helper()

	engine, err := db.NewSqliteEngine(":memory:", zerolog.Nop())
	require.NoError(t, err)

	t.Cleanup(func() {
		engine.Close()
	})

	return New(engine, zerolog.Nop())
}
