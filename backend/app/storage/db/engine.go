package db

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
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
