package database

import (
	"testing"

	"github.com/rs/zerolog"
)

func setupTestDataStorage(t *testing.T) *DataStorage {
	t.Helper()
	return NewDataStorage(setupTestSqliteDB(t), zerolog.Nop())
}
