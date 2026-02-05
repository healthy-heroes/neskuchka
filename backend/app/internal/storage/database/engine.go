package database

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"

	"github.com/healthy-heroes/neskuchka/backend/app/store"
)

// Engine wraps sqlx.DB
type Engine struct {
	*sqlx.DB
}

type Opts struct {
}

func handleSqlError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return store.ErrNotFound
	}

	return err
}
