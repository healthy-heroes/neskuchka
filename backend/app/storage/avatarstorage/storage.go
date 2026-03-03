package avatarstorage

import (
	"github.com/healthy-heroes/neskuchka/backend/app/storage/db"
	"github.com/rs/zerolog"
)

type Storage struct {
	engine *db.Engine
	logger zerolog.Logger
}

func New(engine *db.Engine, logger zerolog.Logger) *Storage {
	return &Storage{
		engine: engine,
		logger: logger.With().Str("pkg", "avatarstorage").Logger(),
	}
}

