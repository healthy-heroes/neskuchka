package db

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	_ "modernc.org/sqlite"
)

// Type is a type of database engine
type Type string

// enum of supported database engines
const (
	Sqlite   Type = "sqlite"
	Postgres Type = "postgres"
)

type DB struct {
	*sqlx.DB
	dbType Type
}

func NewDB(connURL string) (*DB, error) {
	log.Info().Msgf("New database engine, connection: %s", connURL)

	switch {
	case connURL == ":memory:":
		return NewSqlite(connURL)
	case strings.HasSuffix(connURL, ".db"):
		return NewSqlite(connURL)
	}

	return nil, fmt.Errorf("unknown database type in connection string")
}

func NewSqlite(connURL string) (*DB, error) {
	log.Info().Msg("New connection to sqlite")

	db, err := sqlx.Connect("sqlite", connURL)
	if err != nil {
		log.Error().Msgf("Failed to connect to sqlite: %s", err)
		return nil, err
	}

	return &DB{
		db,
		Sqlite,
	}, nil
}
