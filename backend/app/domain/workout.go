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

// editWindowDays is how long past its date a workout stays editable. The screen
// says the rule out loud next to the buttons, so the number and the caption have
// to move together.
const editWindowDays = 1

// IsPublished reports whether participants can already see the workout.
// A workout shows up on its own day and not a moment earlier.
func (w *Workout) IsPublished(now time.Time) bool {
	return !w.Date.After(dayOf(now))
}

// IsEditable reports whether the owner may still change or drop the workout.
// Past that window the workout is what people trained by, and rewriting it
// would rewrite their history.
func (w *Workout) IsEditable(now time.Time) bool {
	return !w.Date.Before(dayOf(now).AddDate(0, 0, -editWindowDays))
}

// dayOf is the calendar day of t as workout dates are stored: midnight, UTC.
// Dates come from time.Parse(time.DateOnly, ...) and carry no zone of their own,
// so comparing them against a zoned clock has to go through here.
func dayOf(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
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

	// Prescription is what to do: "10", "3x1 min", "3x(1+2) @ 70%".
	// Several entries mean several sets shown on separate lines.
	// Free-form strings for now; later they are built from structured parameters.
	Prescription []string

	// Name is the exercise name, kept on the workout until the exercise
	// reference book exists and the name can be resolved through ExerciseSlug.
	Name string
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

const (
	defaultWorkoutLimit = 10
	maxWorkoutLimit     = 50
)

// WorkoutFindCriteria is a criteria for finding workouts
type WorkoutFindCriteria struct {
	Limit  int
	Offset int

	// PublishedOnly drops workouts dated ahead of today. Everything a
	// participant may read goes out with it set.
	PublishedOnly bool
}

// normalize keeps a caller from asking for an unbounded read
func (c *WorkoutFindCriteria) normalize() {
	if c.Limit <= 0 || c.Limit > maxWorkoutLimit {
		c.Limit = defaultWorkoutLimit
	}
	if c.Offset < 0 {
		c.Offset = 0
	}
}

// WorkoutPage is a slice of a track's workouts together with the counters the
// list header shows. Both come from one request, so they cannot disagree.
type WorkoutPage struct {
	Workouts []Workout

	Total   int
	Planned int
}

// GetWorkout gets a workout by id
func (s *Store) GetWorkout(ctx context.Context, wr WorkoutRef) (Workout, error) {
	return s.storage.GetWorkout(ctx, wr)
}

// CreateWorkout creates a new workout
// Generates a new workout id
func (s *Store) CreateWorkout(ctx context.Context, uid UserID, w Workout) (Workout, error) {
	t, err := s.storage.GetTrack(ctx, w.TrackID)
	if err != nil {
		return Workout{}, err
	}

	// Permission check
	if !t.IsOwner(uid) {
		return Workout{}, ErrForbidden
	}

	// A workout written into the past is born locked: nobody could train by it,
	// and the edit window is already closed on it.
	if w.Date.Before(dayOf(time.Now())) {
		return Workout{}, ErrLocked
	}

	w.ID = NewWorkoutID()
	w.clearSlugs()

	return s.storage.CreateWorkout(ctx, w)
}

// UpdateWorkout updates a workout
// updates only safe fields, other should be ignored
// don't check if fields are empty; just update them.
func (s *Store) UpdateWorkout(ctx context.Context, uid UserID, wu Workout) (Workout, error) {
	t, err := s.storage.GetTrack(ctx, wu.TrackID)
	if err != nil {
		return Workout{}, err
	}

	// Permission check
	if !t.IsOwner(uid) {
		return Workout{}, ErrForbidden
	}

	w, err := s.storage.GetWorkout(ctx, wu.Ref())
	if err != nil {
		return Workout{}, err
	}

	// Both the stored date and the one being set have to be inside the window:
	// otherwise an old workout could be dragged forward, or a fresh one parked
	// in the past where nothing can reach it again.
	now := time.Now()
	if !w.IsEditable(now) || !wu.IsEditable(now) {
		return Workout{}, ErrLocked
	}

	w.ApplyUpdate(wu)

	return s.storage.UpdateWorkout(ctx, w)
}

// DeleteWorkout removes a workout from its track
func (s *Store) DeleteWorkout(ctx context.Context, uid UserID, wr WorkoutRef) error {
	t, err := s.storage.GetTrack(ctx, wr.TrackID)
	if err != nil {
		return err
	}

	// Permission check
	if !t.IsOwner(uid) {
		return ErrForbidden
	}

	w, err := s.storage.GetWorkout(ctx, wr)
	if err != nil {
		return err
	}

	if !w.IsEditable(time.Now()) {
		return ErrLocked
	}

	return s.storage.DeleteWorkout(ctx, wr)
}

// FindWorkouts finds published workouts by criteria — the track as participants
// see it.
//
// PublishedOnly is set here rather than left to the caller on purpose: this is
// the one read that is not behind an ownership check, so there must be no way to
// reach it and forget the flag.
func (s *Store) FindWorkouts(ctx context.Context, tid TrackID, criteria WorkoutFindCriteria) ([]Workout, error) {
	criteria.normalize()
	criteria.PublishedOnly = true

	return s.storage.FindWorkouts(ctx, tid, criteria)
}

// FindTrackWorkouts returns a page of the whole track, drafts included, to its
// owner.
func (s *Store) FindTrackWorkouts(
	ctx context.Context, uid UserID, tid TrackID, criteria WorkoutFindCriteria,
) (WorkoutPage, error) {
	t, err := s.storage.GetTrack(ctx, tid)
	if err != nil {
		return WorkoutPage{}, err
	}

	// Permission check
	if !t.IsOwner(uid) {
		return WorkoutPage{}, ErrForbidden
	}

	criteria.normalize()

	workouts, err := s.storage.FindWorkouts(ctx, tid, criteria)
	if err != nil {
		return WorkoutPage{}, err
	}

	// One clock for the page and its counters, so "planned" means the same thing
	// in the rows and in the header above them.
	total, planned, err := s.storage.CountWorkouts(ctx, tid, time.Now())
	if err != nil {
		return WorkoutPage{}, err
	}

	return WorkoutPage{Workouts: workouts, Total: total, Planned: planned}, nil
}
