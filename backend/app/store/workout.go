package store

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type WorkoutID string
type WorkoutProtocol string

const (
	WorkoutProtocolDefault WorkoutProtocol = "default"
)

type Workout struct {
	ID WorkoutID

	Date    time.Time
	TrackID TrackID

	Sections []WorkoutSection
}

type WorkoutSection struct {
	Title     string
	Protocol  WorkoutProtocol
	Exercises []WorkoutExercise
}

type WorkoutExercise struct {
	ExerciseSlug ExerciseSlug

	Repetitions     int
	RepetitionsText string

	Weight     int
	WeightText string
}

func CreateWorkoutId() WorkoutID {
	id := ulid.Make()
	return WorkoutID(id.String())
}

type WorkoutStore interface {
	Create(workout *Workout) (*Workout, error)
	Get(id WorkoutID) (*Workout, error)
	GetList(criteria WorkoutFindCriteria) ([]*Workout, error)
}

type WorkoutFindCriteria struct {
	TrackID TrackID
	Limit   int
}
