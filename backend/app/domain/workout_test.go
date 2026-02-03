package domain

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type WorkoutStoreStub struct {
	GetFunc    func(context.Context, WorkoutID) (Workout, error)
	FindFunc   func(context.Context, WorkoutFindCriteria) ([]Workout, error)
	CreateFunc func(context.Context, Workout) (Workout, error)
	UpdateFunc func(context.Context, Workout) (Workout, error)
}

func (s WorkoutStoreStub) Get(ctx context.Context, id WorkoutID) (Workout, error) {
	return s.GetFunc(ctx, id)
}

func (s WorkoutStoreStub) Find(ctx context.Context, criteria WorkoutFindCriteria) ([]Workout, error) {
	return s.FindFunc(ctx, criteria)
}

func (s WorkoutStoreStub) Create(ctx context.Context, workout Workout) (Workout, error) {
	return s.CreateFunc(ctx, workout)
}

func (s WorkoutStoreStub) Update(ctx context.Context, workout Workout) (Workout, error) {
	return s.UpdateFunc(ctx, workout)
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
		existingWorkout := Workout{
			ID: WorkoutID("1"),
		}

		service := NewService(Opts{
			WorkoutStore: WorkoutStoreStub{
				GetFunc: func(ctx context.Context, id WorkoutID) (Workout, error) {
					return existingWorkout, nil
				},
			},
		})
		workout, err := service.GetWorkout(context.Background(), WorkoutID("1"))
		assert.Nil(t, err)
		assert.Equal(t, existingWorkout, workout)
	})
}

func TestCreateWorkout(t *testing.T) {
	t.Run("should create workout", func(t *testing.T) {
		newWorkout := Workout{
			Date:    time.Now(),
			TrackID: TrackID("1"),
			Sections: []WorkoutSection{
				{
					Title: "Section 1",
					Exercises: []WorkoutExercise{
						{ExerciseSlug: "squat"},
					},
				},
			},
		}

		service := NewService(Opts{
			WorkoutStore: WorkoutStoreStub{
				CreateFunc: func(ctx context.Context, workout Workout) (Workout, error) {
					return workout, nil
				},
			},
		})

		workout, err := service.CreateWorkout(context.Background(), newWorkout)

		assert.Nil(t, err)
		assert.NotEmpty(t, workout.ID)
		assert.NotEqual(t, newWorkout.ID, workout.ID)
		assert.Equal(t, newWorkout.TrackID, workout.TrackID)
		assert.Equal(t, newWorkout.Date, workout.Date)

		// todo: it's wrong behavior, because affect newWorkout.Sections
		assert.Equal(t, newWorkout.Sections, workout.Sections)
	})
}

func TestUpdateWorkout(t *testing.T) {
	t.Run("should update workout", func(t *testing.T) {
		existingWorkout := Workout{
			ID:          WorkoutID("1"),
			TrackID:     TrackID("1"),
			Slug:        WorkoutSlug("workout-1"),
			Status:      WorkoutStatusPublished,
			PublishedAt: time.Now(),

			Date: time.Now(),
			Sections: []WorkoutSection{
				{
					Title: "Section 1",
					Exercises: []WorkoutExercise{
						{ExerciseSlug: "squat"},
					},
				},
			},
			Notes: "Test notes",
		}
		newWorkout := Workout{
			ID:          WorkoutID(existingWorkout.ID),
			TrackID:     TrackID("wrong-track-id"),
			Slug:        WorkoutSlug("wrong-slug"),
			Status:      WorkoutStatus("wrong-status"),
			PublishedAt: time.Now().Add(1 * time.Hour),

			Date: time.Now().Add(1 * time.Hour),
			Sections: []WorkoutSection{
				{
					Title: "Section 2",
					Exercises: []WorkoutExercise{
						{ExerciseSlug: "bench-press"},
					},
				},
			},
			Notes: "Test notes 2",
		}
		service := NewService(Opts{
			WorkoutStore: WorkoutStoreStub{
				GetFunc: func(ctx context.Context, id WorkoutID) (Workout, error) {
					return existingWorkout, nil
				},
				UpdateFunc: func(ctx context.Context, workout Workout) (Workout, error) {
					return workout, nil
				},
			},
		})
		workout, err := service.UpdateWorkout(context.Background(), newWorkout)

		assert.Nil(t, err)
		// Protected fields
		assert.Equal(t, existingWorkout.ID, workout.ID)
		assert.Equal(t, existingWorkout.TrackID, workout.TrackID)
		assert.Equal(t, existingWorkout.Slug, workout.Slug)
		assert.Equal(t, existingWorkout.Status, workout.Status)
		assert.Equal(t, existingWorkout.PublishedAt, workout.PublishedAt)

		// Changable fields
		assert.Equal(t, newWorkout.Date, workout.Date)
		assert.Equal(t, newWorkout.Notes, workout.Notes)
		// todo: it's wrong behavior, because affect newWorkout.Sections
		assert.Equal(t, newWorkout.Sections, workout.Sections)
	})

	t.Run("should return error if workout not found", func(t *testing.T) {
		service := NewService(Opts{
			WorkoutStore: WorkoutStoreStub{
				GetFunc: func(ctx context.Context, id WorkoutID) (Workout, error) {
					return Workout{}, ErrNotFound
				},
			},
		})
		_, err := service.UpdateWorkout(context.Background(), Workout{
			ID: WorkoutID("1"),
		})

		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
	})
}

func TestFindWorkouts(t *testing.T) {
	t.Run("should find workouts", func(t *testing.T) {
		usedCriteria := WorkoutFindCriteria{
			TrackID: TrackID("1"),
			Limit:   10,
		}

		workouts := []Workout{
			{
				ID: WorkoutID("1"),
			},
		}
		service := NewService(Opts{
			WorkoutStore: WorkoutStoreStub{
				FindFunc: func(ctx context.Context, criteria WorkoutFindCriteria) ([]Workout, error) {
					usedCriteria = criteria
					return workouts, nil
				},
			},
		})
		foundWorkouts, err := service.FindWorkouts(context.Background(), WorkoutFindCriteria{
			TrackID: TrackID("1"),
		})
		assert.Nil(t, err)
		assert.Equal(t, workouts, foundWorkouts)
		assert.Equal(t, usedCriteria, WorkoutFindCriteria{
			TrackID: TrackID("1"),
			Limit:   10,
		})
	})

	t.Run("should return error if track id is empty", func(t *testing.T) {
		service := NewService(Opts{
			WorkoutStore: WorkoutStoreStub{
				FindFunc: func(ctx context.Context, criteria WorkoutFindCriteria) ([]Workout, error) {
					return []Workout{}, nil
				},
			},
		})
		_, err := service.FindWorkouts(context.Background(), WorkoutFindCriteria{})
		assert.Error(t, err)
	})

	t.Run("should set default limit if limit is less than 0 or greater than 50", func(t *testing.T) {
		var usedLimit int
		service := NewService(Opts{
			WorkoutStore: WorkoutStoreStub{
				FindFunc: func(ctx context.Context, criteria WorkoutFindCriteria) ([]Workout, error) {
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
			_, err := service.FindWorkouts(context.Background(), WorkoutFindCriteria{
				TrackID: TrackID("1"),
				Limit:   limit,
			})
			assert.Nil(t, err)
			assert.Equal(t, expected, usedLimit, "limit %d", limit)
		}
	})
}
