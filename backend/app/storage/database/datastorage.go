package database

import (
	"github.com/rs/zerolog"
)

// DataStorage stores common data of the application
type DataStorage struct {
	engine *Engine
	logger zerolog.Logger
}

// NewDataStorage creates a new data storage
func NewDataStorage(engine *Engine, logger zerolog.Logger) *DataStorage {
	return &DataStorage{
		engine: engine,
		logger: logger.With().Str("pkg", "database").Logger(),
	}
}
