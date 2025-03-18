package sqlite

import (
	"database/sql"

	"github.com/rs/zerolog/log"

	"github.com/healthy-heroes/neskuchka/backend/app/store"
)

type SqliteExerciseStore struct {
	db *sql.DB
}

func (s *SqliteExerciseStore) Create(exercise *store.Exercise) error {
	_, err := s.db.Exec(`
		INSERT INTO exercise (
			slug,
			name
		) VALUES (?, ?)`,
		exercise.Slug, exercise.Name)
	return err
}

func (s *SqliteExerciseStore) Get(slug store.ExerciseSlug) (*store.Exercise, error) {
	exercise := &store.Exercise{}
	err := s.db.QueryRow(`
		SELECT slug, name 
		FROM exercise 
		WHERE slug = ?
		`, slug).Scan(&exercise.Slug, &exercise.Name)

	if err != nil {
		return nil, err
	}
	return exercise, nil
}

func (s *SqliteExerciseStore) Find(criteria store.ExerciseGetCriteria) ([]*store.Exercise, error) {
	exercises := []*store.Exercise{}

	rows, err := s.db.Query(`
		SELECT slug, name 
		FROM exercise 
		WHERE id > ?
		ORDER BY slug
		LIMIT ?
		`, criteria.AfterID, criteria.Limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		exercise := &store.Exercise{}
		err := rows.Scan(&exercise.Slug, &exercise.Name)
		if err != nil {
			return nil, err
		}
		exercises = append(exercises, exercise)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return exercises, nil
}

func CreateExerciseTables(db *sql.DB) error {
	log.Debug().Msg("Creating exercise table")

	// Create exercises table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS exercise (
			slug TEXT NOT NULL UNIQUE,
			name TEXT NOT NULL
		)
	`)

	if err != nil {
		log.Error().Err(err).Msg("Failed to create exercise table")
		return err
	}

	log.Debug().Msg("Exercise table created")
	return nil
}
