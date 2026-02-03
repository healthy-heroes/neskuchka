package domain

type Service struct {
	userStore    UserStore
	trackStore   TrackStore
	workoutStore WorkoutStore
}

type Opts struct {
	UserStore    UserStore
	TrackStore   TrackStore
	WorkoutStore WorkoutStore
}

// NewService creates a domain service
// its entry point for all domain services
func NewService(opts Opts) *Service {
	return &Service{
		userStore:    opts.UserStore,
		trackStore:   opts.TrackStore,
		workoutStore: opts.WorkoutStore,
	}
}
