package database

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"

	"github.com/healthy-heroes/neskuchka/backend/app/domain"
)

type Type string

const (
	Sqlite Type = "sqlite"
)

// Engine wraps sqlx.DB
type Engine struct {
	*sqlx.DB
}

type Opts struct {
	Logger zerolog.Logger
}

func NewEngine(dbUrl string, opts Opts) (*Engine, error) {
	return NewSqliteEngine(dbUrl, opts.Logger)
}

func handleSqlError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return domain.ErrNotFound
	}

	return err
}
