package datastore

import (
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
