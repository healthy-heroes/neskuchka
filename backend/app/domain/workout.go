package domain

import (
	"context"
	"errors"
	"time"

	"github.com/healthy-heroes/neskuchka/backend/app/internal/uuid"
)

type WorkoutID string

// NewWorkoutID generates a new workout id
func NewWorkoutID() WorkoutID {
	return WorkoutID(uuid.New())
}

// Workout is a workout aggregate
type Workout struct {
	ID      WorkoutID
	TrackID TrackID

	Date  time.Time
	Notes string

	Sections []WorkoutSection
}

// WorkoutSection is a section of a workout
type WorkoutSection struct {
	Title     string
	Protocol  Protocol
	Exercises []WorkoutExercise
}

// WorkoutExercise is an exercise in a workout section
type WorkoutExercise struct {
	ExerciseSlug ExerciseSlug

	Description string
}

// clearSlugs clears the exercise slugs
// NOTE: later we will use it to clear the unknown exercise slugs
func (w *Workout) clearSlugs() {
	for i := range w.Sections {
		for j := range w.Sections[i].Exercises {
			w.Sections[i].Exercises[j].ExerciseSlug = ""
		}
	}
}

// WorkoutFindCriteria is a criteria for finding workouts
type WorkoutFindCriteria struct {
	TrackID TrackID
	Limit   int
}

// GetWorkout gets a workout by id
func (s *Store) GetWorkout(ctx context.Context, id WorkoutID) (Workout, error) {
	return s.dataStorage.GetWorkout(ctx, id)
}

// CreateWorkout creates a new workout
// Generates a new workout id
// todo: clearing slugs should not affect incoming workout
func (s *Store) CreateWorkout(ctx context.Context, trackID TrackID, w Workout) (Workout, error) {
	w.ID = NewWorkoutID()
	w.TrackID = trackID
	w.clearSlugs()

	return s.dataStorage.CreateWorkout(ctx, w)
}

// UpdateWorkout updates a workout
// updates only safe fields, other should be ignored
// don't check if fields are empty; just update them.
// todo: generate slug
// todo: immutable date
// todo: clearing slugs should not affect incoming workout
func (s *Store) UpdateWorkout(ctx context.Context, id WorkoutID, wp Workout) (Workout, error) {
	workout, err := s.dataStorage.GetWorkout(ctx, id)
	if err != nil {
		return Workout{}, err
	}

	workout.Date = wp.Date
	workout.Sections = wp.Sections
	workout.Notes = wp.Notes
	workout.clearSlugs()

	return s.dataStorage.UpdateWorkout(ctx, workout)
}

// FindWorkouts finds workouts by criteria
func (s *Store) FindWorkouts(ctx context.Context, criteria WorkoutFindCriteria) ([]Workout, error) {
	if criteria.TrackID == "" {
		return nil, errors.New("TrackID is required")
	}

	if criteria.Limit <= 0 || criteria.Limit > 50 {
		criteria.Limit = 10
	}

	return s.dataStorage.FindWorkouts(ctx, criteria)
}
