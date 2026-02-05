package database

import (
	"testing"

	"github.com/rs/zerolog"
)

func setupTestDataStorage(t *testing.T, engine *Engine) *DataStorage {
	return NewDataStorage(engine, zerolog.Nop())
}
