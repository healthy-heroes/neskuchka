package datastore

import (
	"database/sql"
	"errors"

	"github.com/healthy-heroes/neskuchka/backend/app/store"
	"github.com/healthy-heroes/neskuchka/backend/app/store/db"
)

type DataStore struct {
	*db.DB

	Exercise store.ExerciseStore
	User     store.UserStore
	Track    store.TrackStore
	Workout  store.WorkoutStore
}

func NewDataStore(db *db.DB) *DataStore {
	dataStore := &DataStore{
		DB: db,
	}

	dataStore.Exercise = &ExerciseDBStore{
		DataStore: dataStore,
	}

	dataStore.User = &UserDBStore{
		DataStore: dataStore,
	}

	dataStore.Track = &TrackDBStore{
		DataStore: dataStore,
	}

	dataStore.Workout = &WorkoutDBStore{
		DataStore: dataStore,
	}

	return dataStore
}

// handleFindError handles the error from the find operation and matches it to store errors
func handleFindError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return store.ErrNotFound
	}

	return err
}
