package db

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	_ "modernc.org/sqlite"
)

// Type is a type of database engine
type Type string

// enum of supported database engines
const (
	Sqlite Type = "sqlite"
	// Postgres Type = "postgres"
)

type DB struct {
	*sqlx.DB
	dbType Type
}

func NewSqlite(fileSource string) (*DB, error) {
	log.Info().Msg("New connection to sqlite")

	db, err := sqlx.Connect("sqlite", fileSource)
	if err != nil {
		log.Error().Msgf("Failed to connect to sqlite: %s", err)
		return nil, err
	}

	return &DB{
		db,
		Sqlite,
	}, nil
}
