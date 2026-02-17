package domain

import (
	"context"
	"time"
)

// DataStorageStub is a stub for dataStorage interface for testing.
type DataStorageStub struct {
	GetUserFunc        func(context.Context, UserID) (User, error)
	GetUserByEmailFunc func(context.Context, Email) (User, error)
	CreateUserFunc     func(context.Context, User) (User, error)
	UpdateUserFunc     func(context.Context, User) (User, error)

	GetTrackFunc       func(context.Context, TrackID) (Track, error)
	GetTrackBySlugFunc func(context.Context, TrackSlug) (Track, error)

	GetWorkoutFunc    func(context.Context, WorkoutRef) (Workout, error)
	FindWorkoutsFunc  func(context.Context, TrackID, WorkoutFindCriteria) ([]Workout, error)
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

func (s *DataStorageStub) GetTrack(ctx context.Context, id TrackID) (Track, error) {
	return s.GetTrackFunc(ctx, id)
}

func (s *DataStorageStub) GetTrackBySlug(ctx context.Context, slug TrackSlug) (Track, error) {
	return s.GetTrackBySlugFunc(ctx, slug)
}

func (s *DataStorageStub) GetWorkout(ctx context.Context, wr WorkoutRef) (Workout, error) {
	return s.GetWorkoutFunc(ctx, wr)
}

func (s *DataStorageStub) FindWorkouts(ctx context.Context, tid TrackID, criteria WorkoutFindCriteria) ([]Workout, error) {
	return s.FindWorkoutsFunc(ctx, tid, criteria)
}

func (s *DataStorageStub) CreateWorkout(ctx context.Context, w Workout) (Workout, error) {
	return s.CreateWorkoutFunc(ctx, w)
}

func (s *DataStorageStub) UpdateWorkout(ctx context.Context, w Workout) (Workout, error) {
	return s.UpdateWorkoutFunc(ctx, w)
}

func createTrack() Track {
	return Track{
		ID:      NewTrackID(),
		Slug:    TrackSlug("main-test"),
		OwnerID: NewUserID(),
	}
}

func createWorkout(trackID TrackID) Workout {
	return Workout{
		ID:      NewWorkoutID(),
		TrackID: trackID,
		Date:    time.Now(),
		Notes:   "Test workout notes",
		Sections: []WorkoutSection{
			{
				Title: "Section 1",
				Exercises: []WorkoutExercise{
					{ExerciseSlug: "exercise-1"},
					{ExerciseSlug: "exercise-2"},
				},
			},
			{
				Title: "Section 2",
				Exercises: []WorkoutExercise{
					{ExerciseSlug: "exercise-3"},
					{ExerciseSlug: "exercise-4"},
				},
			},
		},
	}
}
