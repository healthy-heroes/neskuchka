package sqlite

import (
	"database/sql"
	"errors"

	"github.com/rs/zerolog/log"
	_ "modernc.org/sqlite"

	"github.com/healthy-heroes/neskuchka/backend/app/store"
)

// SqliteStore is implementation of Engine interface for SQLite
type SqliteStore struct {
	db       *sql.DB
	filepath string
}

// NewSqliteStore creates a new SqliteStore
func NewSqliteStore(path string) *SqliteStore {
	return &SqliteStore{
		filepath: path,
	}
}

// Connect connects to the SQLite database
func (s *SqliteStore) Connect() error {
	log.Info().Msgf("Connecting to SQLite database at %s", s.filepath)
	dbPath := s.filepath
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {

		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}

	s.db = db

	// Create tables if not exists
	log.Debug().Msg("Creating tables...")
	if err = s.createTables(); err != nil {
		log.Error().Err(err).Msg("Failed to create tables")

		return err
	}

	return nil
}

// createTables creates all tables in the database
func (s *SqliteStore) createTables() error {
	hasErrs := false

	if err := CreateExerciseTables(s.db); err != nil {
		hasErrs = true
	}

	if hasErrs {
		return errors.New("failed to create tables")
	}

	return nil
}

// Close closes the SQLite database
func (s *SqliteStore) Close() error {
	log.Info().Msg("Closing SQLite database")
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// Exercise creates a new ExerciseStore
func (s *SqliteStore) Exercise() store.ExerciseStore {
	return &SqliteExerciseStore{db: s.db}
}
