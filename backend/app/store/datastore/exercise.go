package datastore

import (
	"github.com/jmoiron/sqlx"
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

func (ds *ExerciseDBStore) Find(criteria *store.ExerciseFindCriteria) ([]*store.Exercise, error) {
	query := "SELECT * FROM exercise"
	if len(criteria.Slugs) > 0 {
		query += " WHERE slug IN (:slugs)"
	}
	query += " ORDER BY slug LIMIT :limit"

	query, args, err := sqlx.Named(query, criteria)
	if err != nil {
		log.Error().Err(err).Msg("Failed sqlx.Named")
		return nil, err
	}

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		log.Error().Err(err).Msg("Failed sqlx.In")
		return nil, err
	}

	query = ds.DB.Rebind(query)

	rows, err := ds.DB.Queryx(query, args...)
	if err != nil {
		log.Error().Err(err).Msg("Failed prepare query")
		return nil, err
	}
	defer rows.Close()

	exercises := make([]*store.Exercise, 0)
	for rows.Next() {
		var exercise store.Exercise
		err = rows.StructScan(&exercise)
		if err != nil {
			log.Error().Err(err).Msg("Failed scan exercise")
			return nil, err
		}
		exercises = append(exercises, &exercise)
	}

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
