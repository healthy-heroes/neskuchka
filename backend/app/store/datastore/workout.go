package datastore

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/healthy-heroes/neskuchka/backend/app/store"
)

type WorkoutDBStore struct {
	*DataStore
}

type WorkoutDB struct {
	ID       string
	Date     time.Time
	TrackID  string `db:"track_id"`
	Sections []byte
}

func workoutDBFromStore(workout *store.Workout) (*WorkoutDB, error) {
	// Serialize sections to JSON
	sectionsJSON, err := json.Marshal(workout.Sections)
	if err != nil {
		log.Error().Err(err).Msg("Failed to serialize sections to JSON")
		return nil, err
	}

	return &WorkoutDB{
		ID:       string(workout.ID),
		Date:     workout.Date,
		TrackID:  string(workout.TrackID),
		Sections: sectionsJSON,
	}, nil
}

func (w *WorkoutDB) toStore() (*store.Workout, error) {
	workout := &store.Workout{
		ID:      store.WorkoutID(w.ID),
		Date:    w.Date,
		TrackID: store.TrackID(w.TrackID),
	}

	// Deserialize sections from JSON
	err := json.Unmarshal(w.Sections, &workout.Sections)
	if err != nil {
		log.Error().Err(err).Msg("Failed to deserialize sections from JSON")
		return nil, err
	}

	return workout, nil
}

func (ds *WorkoutDBStore) Create(workout *store.Workout) (*store.Workout, error) {
	dbWorkout, err := workoutDBFromStore(workout)
	if err != nil {
		return nil, err
	}

	_, err = ds.Exec(`INSERT INTO workout (id, date, track_id, sections) 
						VALUES (?, ?, ?, ?)`,
		dbWorkout.ID, dbWorkout.Date, dbWorkout.TrackID, dbWorkout.Sections)
	if err != nil {
		return nil, err
	}
	return workout, nil
}

func (ds *WorkoutDBStore) Get(id store.WorkoutID) (*store.Workout, error) {
	var dbWorkout WorkoutDB

	err := ds.DB.Get(&dbWorkout, `SELECT id, date, track_id, sections FROM workout WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}

	workout, err := dbWorkout.toStore()
	if err != nil {
		return nil, err
	}

	return workout, nil
}

func (ds *WorkoutDBStore) Find(criteria *store.WorkoutFindCriteria) ([]*store.Workout, error) {
	if criteria.TrackID == "" {
		return nil, errors.New("criteria.TrackID is required")
	}

	query := `SELECT id, date, track_id, sections FROM workout WHERE track_id = ? ORDER BY date DESC LIMIT ?`

	limit := criteria.Limit
	if limit <= 0 {
		limit = 50
	}

	dbWorkouts := make([]*WorkoutDB, 0)
	err := ds.Select(&dbWorkouts, query, criteria.TrackID, limit)
	if err != nil {
		return nil, err
	}

	workouts := make([]*store.Workout, 0)
	for _, dbWorkout := range dbWorkouts {
		workout, err := dbWorkout.toStore()
		if err != nil {
			return nil, err
		}
		workouts = append(workouts, workout)
	}

	return workouts, nil
}

// InitTables initializes the workout table
func (ds *WorkoutDBStore) InitTables() error {
	log.Debug().Msg("Creating workout table")

	_, err := ds.Exec(`
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
