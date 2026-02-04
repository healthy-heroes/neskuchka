package domain

import "context"

type dataStorage interface {
	GetUser(context.Context, UserID) (User, error)
	GetUserByEmail(context.Context, Email) (User, error)
	CreateUser(context.Context, User) (User, error)
	UpdateUser(context.Context, User) (User, error)

	GetTrackBySlug(context.Context, TrackSlug) (Track, error)

	GetWorkout(context.Context, WorkoutID) (Workout, error)
	FindWorkouts(context.Context, WorkoutFindCriteria) ([]Workout, error)
	CreateWorkout(context.Context, Workout) (Workout, error)
	UpdateWorkout(context.Context, Workout) (Workout, error)
}

// Store is a domain store
// it contains all domain services
type Store struct {
	dataStorage dataStorage
}

// Opts contains options for the store
type Opts struct {
	DataStorage dataStorage
}

// NewStore creates a domain store
func NewStore(opts Opts) *Store {
	return &Store{
		dataStorage: opts.DataStorage,
	}
}
