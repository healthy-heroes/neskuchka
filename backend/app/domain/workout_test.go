package domain

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func wrf(tid TrackID, wid WorkoutID) WorkoutRef {
	return WorkoutRef{TrackID: tid, WorkoutID: wid}
}

func TestNewWorkoutID(t *testing.T) {
	t.Run("should generate a new workout id", func(t *testing.T) {
		workoutID := NewWorkoutID()
		assert.NotEmpty(t, workoutID)
	})
}

func TestClearSlugs(t *testing.T) {
	tests := []struct {
		name     string
		sections []WorkoutSection
		expected []WorkoutSection
	}{
		{
			name: "All exercises known",
			sections: []WorkoutSection{
				{
					Title: "Section 1",
					Exercises: []WorkoutExercise{
						{ExerciseSlug: "squat"},
						{ExerciseSlug: "bench-press"},
					},
				},
				{
					Title: "Section 2",
					Exercises: []WorkoutExercise{
						{ExerciseSlug: "deadlift"},
					},
				},
			},
			expected: []WorkoutSection{
				{
					Title: "Section 1",
					Exercises: []WorkoutExercise{
						{ExerciseSlug: ""},
						{ExerciseSlug: ""},
					},
				},
				{
					Title: "Section 2",
					Exercises: []WorkoutExercise{
						{ExerciseSlug: ""},
					},
				},
			},
		},
		{
			name: "No exercises in sections",
			sections: []WorkoutSection{
				{
					Title:     "Section 1",
					Exercises: []WorkoutExercise{},
				},
				{
					Title:     "Section 2",
					Exercises: []WorkoutExercise{},
				},
			},
			expected: []WorkoutSection{
				{
					Title:     "Section 1",
					Exercises: []WorkoutExercise{},
				},
				{
					Title:     "Section 2",
					Exercises: []WorkoutExercise{},
				},
			},
		},
		{
			name:     "Empty sections slice",
			sections: []WorkoutSection{},
			expected: []WorkoutSection{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workout := Workout{
				Sections: tt.sections,
			}

			workout.clearSlugs()
			assert.Equal(t, tt.expected, workout.Sections)
		})
	}
}

func TestGetWorkout(t *testing.T) {
	t.Run("should return workout", func(t *testing.T) {
		w := createWorkout(NewTrackID())
		service := NewStore(Opts{
			DataStorage: &DataStorageStub{
				GetWorkoutFunc: func(ctx context.Context, wr WorkoutRef) (Workout, error) {
					return w, nil
				},
			},
		})

		workout, err := service.GetWorkout(context.Background(), wrf(w.TrackID, w.ID))
		assert.Nil(t, err)
		assert.Equal(t, w, workout)
	})
}

func TestCreateWorkout(t *testing.T) {
	t.Run("should create workout", func(t *testing.T) {
		track := createTrack()

		service := NewStore(Opts{
			DataStorage: &DataStorageStub{
				GetTrackFunc: func(ctx context.Context, tid TrackID) (Track, error) {
					return track, nil
				},
				CreateWorkoutFunc: func(ctx context.Context, w Workout) (Workout, error) {
					return w, nil
				},
			},
		})

		newWorkout := createWorkout(track.ID)
		newWorkout.ID = ""
		workout, err := service.CreateWorkout(context.Background(), track.OwnerID, newWorkout)

		assert.NoError(t, err)
		assert.NotEmpty(t, workout.ID)

		assert.Equal(t, newWorkout.TrackID, workout.TrackID)
		assert.Equal(t, newWorkout.Date, workout.Date)
		assert.Equal(t, newWorkout.Notes, workout.Notes)
		assert.Equal(t, newWorkout.Sections, workout.Sections)

		// todo: check that slugs are cleared in workout.Sections but not in newWorkout.Sections
	})

	t.Run("should return error if create not owner of track", func(t *testing.T) {
		track1 := createTrack()
		track2 := createTrack()

		service := NewStore(Opts{
			DataStorage: &DataStorageStub{
				GetTrackFunc: func(ctx context.Context, tid TrackID) (Track, error) {
					if tid == track1.ID {
						return track1, nil
					}
					return track2, nil
				},
			},
		})

		newWorkout := createWorkout(track2.ID)
		_, err := service.CreateWorkout(context.Background(), track1.OwnerID, newWorkout)

		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrForbidden)
	})
}

func TestUpdateWorkout(t *testing.T) {
	t.Run("should update workout", func(t *testing.T) {
		track := createTrack()
		workout := createWorkout(track.ID)
		service := NewStore(Opts{
			DataStorage: &DataStorageStub{
				GetTrackFunc: func(ctx context.Context, tid TrackID) (Track, error) {
					return track, nil
				},
				GetWorkoutFunc: func(ctx context.Context, wr WorkoutRef) (Workout, error) {
					return workout, nil
				},
				UpdateWorkoutFunc: func(ctx context.Context, w Workout) (Workout, error) {
					return w, nil
				},
			},
		})

		newWorkout := Workout{
			ID:      workout.ID,
			TrackID: track.ID,
			Date:    workout.Date.Add(1 * time.Hour),
			Notes:   workout.Notes + " updated",
			Sections: []WorkoutSection{
				{Title: "Section 1 updated", Exercises: []WorkoutExercise{{ExerciseSlug: "exercise-1 updated"}}},
			},
		}
		updated, err := service.UpdateWorkout(context.Background(), track.OwnerID, newWorkout)

		assert.Nil(t, err)
		// Protected fields
		assert.Equal(t, workout.ID, updated.ID)
		assert.Equal(t, workout.TrackID, updated.TrackID)

		// Changable fields
		assert.Equal(t, newWorkout.Date, updated.Date)
		assert.Equal(t, newWorkout.Notes, updated.Notes)
		assert.Equal(t, newWorkout.Sections, updated.Sections)

		// todo: check that slugs are cleared in updated.Sections but not in newWorkout.Sections
	})

	t.Run("should return error if workout not found", func(t *testing.T) {
		track := createTrack()
		service := NewStore(Opts{
			DataStorage: &DataStorageStub{
				GetTrackFunc: func(ctx context.Context, tid TrackID) (Track, error) {
					return track, nil
				},
				GetWorkoutFunc: func(ctx context.Context, wr WorkoutRef) (Workout, error) {
					return Workout{}, ErrNotFound
				},
			},
		})
		_, err := service.UpdateWorkout(context.Background(), track.OwnerID, createWorkout(track.ID))

		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("should return error if update not owner of track", func(t *testing.T) {
		track1 := createTrack()
		track2 := createTrack()
		workout := createWorkout(track1.ID)

		service := NewStore(Opts{
			DataStorage: &DataStorageStub{
				GetTrackFunc: func(ctx context.Context, tid TrackID) (Track, error) {
					if tid == track1.ID {
						return track1, nil
					}
					return track2, nil
				},
			},
		})
		_, err := service.UpdateWorkout(context.Background(), track2.OwnerID, workout)

		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrForbidden)
	})
}

func TestFindWorkouts(t *testing.T) {
	t.Run("should found workouts", func(t *testing.T) {

		var usedCriteria WorkoutFindCriteria
		track := createTrack()
		workouts := []Workout{
			createWorkout(track.ID),
			createWorkout(track.ID),
		}
		service := NewStore(Opts{
			DataStorage: &DataStorageStub{
				FindWorkoutsFunc: func(ctx context.Context, tid TrackID, criteria WorkoutFindCriteria) ([]Workout, error) {
					usedCriteria = criteria
					return workouts, nil
				},
			},
		})
		foundWorkouts, err := service.FindWorkouts(context.Background(), track.ID, WorkoutFindCriteria{Limit: 5})
		assert.NoError(t, err)
		assert.Equal(t, workouts, foundWorkouts)
		assert.Equal(t, WorkoutFindCriteria{Limit: 5}, usedCriteria)
	})

	t.Run("should set default limit if limit is less than 0 or greater than 50", func(t *testing.T) {
		var usedLimit int
		service := NewStore(Opts{
			DataStorage: &DataStorageStub{
				FindWorkoutsFunc: func(ctx context.Context, tid TrackID, criteria WorkoutFindCriteria) ([]Workout, error) {
					usedLimit = criteria.Limit
					return []Workout{}, nil
				},
			},
		})

		tcs := map[int]int{
			-1: 10,
			0:  10,
			1:  1,
			50: 50,
			51: 10,
		}
		for limit, expected := range tcs {
			_, err := service.FindWorkouts(context.Background(), NewTrackID(), WorkoutFindCriteria{
				Limit: limit,
			})
			assert.Nil(t, err)
			assert.Equal(t, expected, usedLimit, "limit %d", limit)
		}
	})
}
