package datastore

import (
	"encoding/json"
	"errors"

	"github.com/rs/zerolog/log"

	"github.com/healthy-heroes/neskuchka/backend/app/store"
)

type WorkoutDBStore struct {
	*DataStore
}

func (ds *WorkoutDBStore) Create(workout *store.Workout) (*store.Workout, error) {
	// Serialize sections to JSON
	sectionsJSON, err := json.Marshal(workout.Sections)
	if err != nil {
		return nil, err
	}

	_, err = ds.DB.Exec(`INSERT INTO workout (id, date, track_id, sections) 
						VALUES (?, ?, ?, ?)`,
		workout.ID, workout.Date, workout.TrackID, sectionsJSON)
	if err != nil {
		return nil, err
	}
	return workout, nil
}

func (ds *WorkoutDBStore) Get(id store.WorkoutID) (*store.Workout, error) {
	var workout store.Workout
	var sectionsJSON []byte

	err := ds.DB.QueryRow(`SELECT id, date, track_id, sections FROM workout WHERE id = ?`,
		id).Scan(&workout.ID, &workout.Date, &workout.TrackID, &sectionsJSON)
	if err != nil {
		return nil, err
	}

	// Deserialize sections from JSON
	err = json.Unmarshal(sectionsJSON, &workout.Sections)
	if err != nil {
		log.Error().Err(err).Msg("Failed to deserialize sections from JSON")
		return nil, err
	}

	log.Debug().Msgf("workout: %+v", workout)

	return &workout, nil
}

func (ds *WorkoutDBStore) Find(criteria *store.WorkoutFindCriteria) ([]*store.Workout, error) {
	if criteria.TrackID == "" {
		return nil, errors.New("track_id is required")
	}

	var (
		workouts []*store.Workout
		args     []interface{}
		query    string
	)

	// Build query based on criteria
	query = `
		SELECT id, date, track_id, sections 
		FROM workout WHERE track_id = ? ORDER BY date DESC LIMIT ?`
	args = append(args, criteria.TrackID)

	// Apply limit
	limit := criteria.Limit
	if limit <= 0 {
		limit = 10 // Default limit
	}
	args = append(args, limit)

	// Execute query
	rows, err := ds.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Process results
	for rows.Next() {
		var workout store.Workout
		var sectionsJSON []byte

		err := rows.Scan(&workout.ID, &workout.Date, &workout.TrackID, &sectionsJSON)
		if err != nil {
			return nil, err
		}

		// Deserialize sections from JSON
		err = json.Unmarshal(sectionsJSON, &workout.Sections)
		if err != nil {
			return nil, err
		}

		workouts = append(workouts, &workout)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return workouts, nil
}

func (ds *WorkoutDBStore) InitTables() error {
	log.Debug().Msg("Creating workout table")

	// Create workout table
	_, err := ds.DB.Exec(`
		CREATE TABLE IF NOT EXISTS workout (
			id TEXT PRIMARY KEY NOT NULL,
			date TIMESTAMP NOT NULL,
			track_id TEXT NOT NULL,
			sections TEXT NOT NULL
		)
	`)

	if err != nil {
		log.Error().Err(err).Msg("Failed to create workout table")
		return err
	}

	log.Debug().Msg("Workout table created")
	return nil
}
