package database

import (
	"testing"

	"github.com/rs/zerolog"
)

func setupTestDataStorage(t *testing.T) *DataStorage {
	return NewDataStorage(setupTestSqliteDB(t), zerolog.Nop())
}
