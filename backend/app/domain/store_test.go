package domain

import (
	"context"
	"time"
)

// StorageStub is a stub for dataStorage interface for testing.
type StorageStub struct {
	GetUserFunc        func(context.Context, UserID) (User, error)
	GetUserByEmailFunc func(context.Context, Email) (User, error)
	CreateUserFunc     func(context.Context, User) (User, error)
	UpdateUserFunc     func(context.Context, User) (User, error)

	GetTrackFunc       func(context.Context, TrackID) (Track, error)
	GetTrackBySlugFunc func(context.Context, TrackSlug) (Track, error)
	UpdateTrackFunc    func(context.Context, Track) (Track, error)

	GetWorkoutFunc    func(context.Context, WorkoutRef) (Workout, error)
	FindWorkoutsFunc  func(context.Context, TrackID, WorkoutFindCriteria) ([]Workout, error)
	CreateWorkoutFunc func(context.Context, Workout) (Workout, error)
	UpdateWorkoutFunc func(context.Context, Workout) (Workout, error)
	DeleteWorkoutFunc func(context.Context, WorkoutRef) error
	CountWorkoutsFunc func(context.Context, TrackID, time.Time) (int, int, error)
}

func (s *StorageStub) GetUser(ctx context.Context, id UserID) (User, error) {
	return s.GetUserFunc(ctx, id)
}

func (s *StorageStub) GetUserByEmail(ctx context.Context, email Email) (User, error) {
	return s.GetUserByEmailFunc(ctx, email)
}

func (s *StorageStub) CreateUser(ctx context.Context, user User) (User, error) {
	return s.CreateUserFunc(ctx, user)
}

func (s *StorageStub) UpdateUser(ctx context.Context, user User) (User, error) {
	return s.UpdateUserFunc(ctx, user)
}

func (s *StorageStub) GetTrack(ctx context.Context, id TrackID) (Track, error) {
	return s.GetTrackFunc(ctx, id)
}

func (s *StorageStub) GetTrackBySlug(ctx context.Context, slug TrackSlug) (Track, error) {
	return s.GetTrackBySlugFunc(ctx, slug)
}

func (s *StorageStub) UpdateTrack(ctx context.Context, track Track) (Track, error) {
	return s.UpdateTrackFunc(ctx, track)
}

func (s *StorageStub) GetWorkout(ctx context.Context, wr WorkoutRef) (Workout, error) {
	return s.GetWorkoutFunc(ctx, wr)
}

func (s *StorageStub) FindWorkouts(ctx context.Context, tid TrackID, criteria WorkoutFindCriteria) ([]Workout, error) {
	return s.FindWorkoutsFunc(ctx, tid, criteria)
}

func (s *StorageStub) CreateWorkout(ctx context.Context, w Workout) (Workout, error) {
	return s.CreateWorkoutFunc(ctx, w)
}

func (s *StorageStub) UpdateWorkout(ctx context.Context, w Workout) (Workout, error) {
	return s.UpdateWorkoutFunc(ctx, w)
}

func (s *StorageStub) DeleteWorkout(ctx context.Context, wr WorkoutRef) error {
	return s.DeleteWorkoutFunc(ctx, wr)
}

func (s *StorageStub) CountWorkouts(ctx context.Context, tid TrackID, now time.Time) (int, int, error) {
	return s.CountWorkoutsFunc(ctx, tid, now)
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
		// dayOf, not time.Now(): a workout date is a calendar day, and the
		// window predicates compare against one. A zoned instant here makes the
		// test pass or fail by the hour it runs at.
		Date:  dayOf(time.Now()),
		Notes: "Test workout notes",
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
