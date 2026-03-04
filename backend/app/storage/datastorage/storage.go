package datastorage

import (
	"github.com/rs/zerolog"

	"github.com/healthy-heroes/neskuchka/backend/app/storage/db"
)

type Storage struct {
	engine *db.Engine
	logger zerolog.Logger
}

func New(engine *db.Engine, logger zerolog.Logger) *Storage {
	return &Storage{
		engine: engine,
		logger: logger.With().Str("pkg", "datastorage").Logger(),
	}
}

