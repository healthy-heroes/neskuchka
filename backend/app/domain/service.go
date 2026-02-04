package domain

type Store struct {
	userRepo    UserRepo
	trackRepo   TrackRepo
	workoutRepo WorkoutRepo
}

type Opts struct {
	UserRepo    UserRepo
	TrackRepo   TrackRepo
	WorkoutRepo WorkoutRepo
}

// NewStore creates a domain service
// its entry point for all domain services
func NewStore(opts Opts) *Store {
	return &Store{
		userRepo:    opts.UserRepo,
		trackRepo:   opts.TrackRepo,
		workoutRepo: opts.WorkoutRepo,
	}
}
