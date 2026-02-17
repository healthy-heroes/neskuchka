package domain

import (
	"context"
	"time"

	"github.com/healthy-heroes/neskuchka/backend/app/internal/uuid"
)

type WorkoutID string

// NewWorkoutID generates a new workout id
func NewWorkoutID() WorkoutID {
	return WorkoutID(uuid.New())
}

type WorkoutRef struct {
	TrackID   TrackID
	WorkoutID WorkoutID
}

// Workout is a workout aggregate
type Workout struct {
	ID      WorkoutID
	TrackID TrackID

	Date  time.Time
	Notes string

	Sections []WorkoutSection
}

func (w *Workout) Ref() WorkoutRef {
	return WorkoutRef{TrackID: w.TrackID, WorkoutID: w.ID}
}

// ApplyUpdate applies an update to a workout
func (w *Workout) ApplyUpdate(wu Workout) {
	w.Date = wu.Date
	w.Notes = wu.Notes
	w.Sections = wu.Sections

	// todo: exists bug with clearing slugs in wu.Sections
	w.clearSlugs()
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
	Limit int
}

// GetWorkout gets a workout by id
func (s *Store) GetWorkout(ctx context.Context, wr WorkoutRef) (Workout, error) {
	return s.dataStorage.GetWorkout(ctx, wr)
}

// CreateWorkout creates a new workout
// Generates a new workout id
func (s *Store) CreateWorkout(ctx context.Context, uid UserID, w Workout) (Workout, error) {
	t, err := s.dataStorage.GetTrack(ctx, w.TrackID)
	if err != nil {
		return Workout{}, err
	}

	// Permission check
	if !t.IsOwner(uid) {
		return Workout{}, ErrForbidden
	}

	w.ID = NewWorkoutID()
	w.clearSlugs()

	return s.dataStorage.CreateWorkout(ctx, w)
}

// UpdateWorkout updates a workout
// updates only safe fields, other should be ignored
// don't check if fields are empty; just update them.
func (s *Store) UpdateWorkout(ctx context.Context, uid UserID, wu Workout) (Workout, error) {
	t, err := s.dataStorage.GetTrack(ctx, wu.TrackID)
	if err != nil {
		return Workout{}, err
	}

	// Permission check
	if !t.IsOwner(uid) {
		return Workout{}, ErrForbidden
	}

	w, err := s.dataStorage.GetWorkout(ctx, wu.Ref())
	if err != nil {
		return Workout{}, err
	}

	w.ApplyUpdate(wu)

	return s.dataStorage.UpdateWorkout(ctx, w)
}

// FindWorkouts finds workouts by criteria
func (s *Store) FindWorkouts(ctx context.Context, tid TrackID, criteria WorkoutFindCriteria) ([]Workout, error) {
	if criteria.Limit <= 0 || criteria.Limit > 50 {
		criteria.Limit = 10
	}

	return s.dataStorage.FindWorkouts(ctx, tid, criteria)
}
