package avatarstorage

import (
	"testing"

	"github.com/rs/zerolog"

	"github.com/healthy-heroes/neskuchka/backend/app/internal/testutil"
)

func setupTestStorage(t *testing.T) *Storage {
	t.Helper()

	return New(testutil.NewEngine(t), zerolog.Nop())
}
