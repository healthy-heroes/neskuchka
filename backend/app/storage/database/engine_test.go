package database

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func setupTestSqliteDB(t *testing.T) *Engine {
	engine, err := NewSqliteEngine(":memory:", zerolog.Nop())
	assert.NoError(t, err)

	return engine
}
