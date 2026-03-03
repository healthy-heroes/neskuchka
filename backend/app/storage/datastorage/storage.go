package datastorage

import (
	"database/sql"
	"errors"

	"github.com/rs/zerolog"

	"github.com/healthy-heroes/neskuchka/backend/app/domain"
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

func (s *Storage) Close() {
	s.engine.Close()
}

func handleSqlError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return domain.ErrNotFound
	}

	return err
}
