package domain

import (
	"context"
	"time"
)

type UrlSuffix string

type storage interface {
	GetUser(context.Context, UserID) (User, error)
	GetUserByEmail(context.Context, Email) (User, error)
	CreateUser(context.Context, User) (User, error)
	UpdateUser(context.Context, User) (User, error)

	GetTrack(context.Context, TrackID) (Track, error)
	GetTrackBySlug(context.Context, TrackSlug) (Track, error)
	UpdateTrack(context.Context, Track) (Track, error)

	GetWorkout(context.Context, WorkoutRef) (Workout, error)
	FindWorkouts(context.Context, TrackID, WorkoutFindCriteria, time.Time) ([]Workout, error)
	CreateWorkout(context.Context, Workout) (Workout, error)
	UpdateWorkout(context.Context, Workout) (Workout, error)
	DeleteWorkout(context.Context, WorkoutRef) error

	// CountWorkouts returns how many workouts the track holds in total and how
	// many of them are still ahead of the given moment.
	CountWorkouts(context.Context, TrackID, time.Time) (int, int, error)
}

// Store is a domain store
// it contains all domain services
type Store struct {
	storage storage
}

// Opts contains options for the store
type Opts struct {
	Storage storage
}

// NewStore creates a domain store
func NewStore(opts Opts) *Store {
	return &Store{
		storage: opts.Storage,
	}
}
