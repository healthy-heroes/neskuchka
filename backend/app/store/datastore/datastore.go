package datastore

import (
	"github.com/healthy-heroes/neskuchka/backend/app/store/db"
)

type DataStore struct {
	*db.DB

	Exercise *ExerciseDBStore
}

func NewDataStore(db *db.DB) *DataStore {
	dataStore := &DataStore{
		DB: db,
	}

	dataStore.Exercise = &ExerciseDBStore{
		DataStore: dataStore,
	}

	return dataStore
}
