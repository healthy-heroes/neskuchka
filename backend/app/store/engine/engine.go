package engine

import (
	"github.com/healthy-heroes/neskuchka/backend/app/store"
)

type Engine interface {
	Connect() error
	Close() error

	// Exercise returns an ExerciseStore
	Exercise() store.ExerciseStore
}
