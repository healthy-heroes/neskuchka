package datastore

import (
	"github.com/rs/zerolog/log"

	"github.com/healthy-heroes/neskuchka/backend/app/store"
)

type ExerciseDBStore struct {
	*DataStore
}

func (ds *ExerciseDBStore) Create(exercise *store.Exercise) (*store.Exercise, error) {
	_, err := ds.DB.Exec(`INSERT INTO exercise (slug, name) VALUES (?, ?)`,
		exercise.Slug, exercise.Name)
	if err != nil {
		return nil, err
	}
	return exercise, nil
}

func (ds *ExerciseDBStore) Get(slug store.ExerciseSlug) (*store.Exercise, error) {
	exercise := &store.Exercise{}
	err := ds.DB.Get(exercise, `SELECT * FROM exercise WHERE slug = ?`, slug)

	log.Info().Msgf("Get exercise: %+v", exercise)

	if err != nil {
		return nil, err
	}
	return exercise, nil
}

func (ds *ExerciseDBStore) Find(criteria store.ExerciseFindCriteria) ([]*store.Exercise, error) {
	exercises := []*store.Exercise{}

	err := ds.DB.Select(&exercises,
		`SELECT * FROM exercise ORDER BY slug LIMIT ?`,
		criteria.Limit)

	return exercises, err
}

func (ds *ExerciseDBStore) InitTables() error {
	log.Debug().Msg("Creating exercise table")

	// Create exercises table
	_, err := ds.DB.Exec(`
		CREATE TABLE IF NOT EXISTS exercise (
			slug TEXT PRIMARY KEY NOT NULL,
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
