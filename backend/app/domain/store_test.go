package domain

import "context"

// DataStorageStub is a stub for dataStorage interface for testing.
type DataStorageStub struct {
	GetUserFunc        func(context.Context, UserID) (User, error)
	GetUserByEmailFunc func(context.Context, Email) (User, error)
	CreateUserFunc     func(context.Context, User) (User, error)
	UpdateUserFunc     func(context.Context, User) (User, error)

	GetTrackBySlugFunc func(context.Context, TrackSlug) (Track, error)

	GetWorkoutFunc    func(context.Context, WorkoutID) (Workout, error)
	FindWorkoutsFunc  func(context.Context, WorkoutFindCriteria) ([]Workout, error)
	CreateWorkoutFunc func(context.Context, Workout) (Workout, error)
	UpdateWorkoutFunc func(context.Context, Workout) (Workout, error)
}

func (s *DataStorageStub) GetUser(ctx context.Context, id UserID) (User, error) {
	return s.GetUserFunc(ctx, id)
}

func (s *DataStorageStub) GetUserByEmail(ctx context.Context, email Email) (User, error) {
	return s.GetUserByEmailFunc(ctx, email)
}

func (s *DataStorageStub) CreateUser(ctx context.Context, user User) (User, error) {
	return s.CreateUserFunc(ctx, user)
}

func (s *DataStorageStub) UpdateUser(ctx context.Context, user User) (User, error) {
	return s.UpdateUserFunc(ctx, user)
}

func (s *DataStorageStub) GetTrackBySlug(ctx context.Context, slug TrackSlug) (Track, error) {
	return s.GetTrackBySlugFunc(ctx, slug)
}

func (s *DataStorageStub) GetWorkout(ctx context.Context, id WorkoutID) (Workout, error) {
	return s.GetWorkoutFunc(ctx, id)
}

func (s *DataStorageStub) FindWorkouts(ctx context.Context, criteria WorkoutFindCriteria) ([]Workout, error) {
	return s.FindWorkoutsFunc(ctx, criteria)
}

func (s *DataStorageStub) CreateWorkout(ctx context.Context, workout Workout) (Workout, error) {
	return s.CreateWorkoutFunc(ctx, workout)
}

func (s *DataStorageStub) UpdateWorkout(ctx context.Context, workout Workout) (Workout, error) {
	return s.UpdateWorkoutFunc(ctx, workout)
}
