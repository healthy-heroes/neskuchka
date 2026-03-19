package domain

import "context"

type UrlSuffix string

type storage interface {
	GetUser(context.Context, UserID) (User, error)
	GetUserByEmail(context.Context, Email) (User, error)
	CreateUser(context.Context, User) (User, error)
	UpdateUser(context.Context, User) (User, error)

	GetTrack(context.Context, TrackID) (Track, error)
	GetTrackBySlug(context.Context, TrackSlug) (Track, error)

	GetWorkout(context.Context, WorkoutRef) (Workout, error)
	FindWorkouts(context.Context, TrackID, WorkoutFindCriteria) ([]Workout, error)
	CreateWorkout(context.Context, Workout) (Workout, error)
	UpdateWorkout(context.Context, Workout) (Workout, error)
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
