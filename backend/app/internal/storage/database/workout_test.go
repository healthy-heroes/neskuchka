package database

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/healthy-heroes/neskuchka/backend/app/domain"
)

func workoutFromDB(t *testing.T, engine *Engine, id string) workoutRow {
	row := workoutRow{}
	err := engine.Get(&row, "SELECT * FROM workout WHERE id = ?", id)
	require.NoError(t, err)

	return row
}

func workoutDate(s string) time.Time {
	t, err := time.Parse(time.DateOnly, s)
	if err != nil {
		panic(err)
	}
	return t
}

func Test_Workout_Create(t *testing.T) {
	ds := setupTestDataStorage(t)
	defer ds.engine.Close()

	newWorkout := domain.Workout{
		ID:      domain.NewWorkoutID(),
		TrackID: domain.NewTrackID(),
		Date:    workoutDate("2025-02-06"),
		Notes:   "Test workout notes",
		Sections: []domain.WorkoutSection{
			{
				Title: "Warm up",
				Protocol: domain.Protocol{
					Type:        domain.ProtocolTypeCustom,
					Title:       "Custom",
					Description: "",
				},
				Exercises: []domain.WorkoutExercise{
					{ExerciseSlug: "squat", Description: "10 reps"},
				},
			},
		},
	}

	createdWorkout, err := ds.CreateWorkout(context.Background(), newWorkout)
	require.NoError(t, err)
	assert.Equal(t, newWorkout.ID, createdWorkout.ID)
	assert.Equal(t, newWorkout.TrackID, createdWorkout.TrackID)
	assert.Equal(t, newWorkout.Date, createdWorkout.Date)
	assert.Equal(t, newWorkout.Notes, createdWorkout.Notes)
	assert.Equal(t, newWorkout.Sections, createdWorkout.Sections)

	workoutByID, err := ds.GetWorkout(context.Background(), newWorkout.ID)
	require.NoError(t, err)
	assert.Equal(t, createdWorkout, workoutByID)

	// checks system fields
	row := workoutFromDB(t, ds.engine, string(newWorkout.ID))
	assert.NotZero(t, row.CreatedAt)
	assert.NotZero(t, row.UpdatedAt)
}

func Test_Workout_Get_NotFound(t *testing.T) {
	ds := setupTestDataStorage(t)
	defer ds.engine.Close()

	_, err := ds.GetWorkout(context.Background(), domain.WorkoutID("non-existent-id"))
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func Test_Workout_Update(t *testing.T) {
	ds := setupTestDataStorage(t)
	defer ds.engine.Close()

	existingWorkout := domain.Workout{
		ID:      domain.NewWorkoutID(),
		TrackID: domain.NewTrackID(),
		Date:    workoutDate("2025-02-01"),
		Notes:   "Original notes",
		Sections: []domain.WorkoutSection{
			{
				Title:    "Section 1",
				Protocol: domain.Protocol{Type: domain.ProtocolTypeCustom},
				Exercises: []domain.WorkoutExercise{
					{ExerciseSlug: "squat", Description: "5x5"},
				},
			},
		},
	}
	createdWorkout, err := ds.CreateWorkout(context.Background(), existingWorkout)
	require.NoError(t, err)

	createdRow := workoutFromDB(t, ds.engine, string(createdWorkout.ID))

	updateWorkout := domain.Workout{
		ID:      createdWorkout.ID,
		TrackID: createdWorkout.TrackID,
		Date:    workoutDate("2025-02-06"),
		Notes:   "Updated notes",
		Sections: []domain.WorkoutSection{
			{
				Title:    "Section 1",
				Protocol: domain.Protocol{Type: domain.ProtocolTypeCustom},
				Exercises: []domain.WorkoutExercise{
					{ExerciseSlug: "bench-press", Description: "3x10"},
				},
			},
		},
	}
	updated, err := ds.UpdateWorkout(context.Background(), updateWorkout)
	require.NoError(t, err)
	assert.Equal(t, updateWorkout.Date, updated.Date)
	assert.Equal(t, updateWorkout.Notes, updated.Notes)
	assert.Equal(t, updateWorkout.Sections, updated.Sections)

	row := workoutFromDB(t, ds.engine, string(updated.ID))
	assert.Equal(t, createdRow.CreatedAt, row.CreatedAt)
	assert.GreaterOrEqual(t, row.UpdatedAt, createdRow.UpdatedAt)
}

func Test_Workout_FindWorkouts(t *testing.T) {
	ds := setupTestDataStorage(t)
	defer ds.engine.Close()

	// create workouts
	trackID := domain.NewTrackID()
	dates := []string{"2025-02-01", "2025-02-05", "2025-02-04"}
	for _, date := range dates {
		_, err := ds.CreateWorkout(context.Background(), domain.Workout{
			ID:      domain.NewWorkoutID(),
			TrackID: trackID,
			Date:    workoutDate(date),
		})
		require.NoError(t, err)
	}

	_, err := ds.CreateWorkout(context.Background(), domain.Workout{
		ID:      domain.NewWorkoutID(),
		TrackID: domain.NewTrackID(),
		Date:    workoutDate("2026-02-02"),
		Notes:   "Other track",
	})
	require.NoError(t, err)

	// full list
	list, err := ds.FindWorkouts(context.Background(), domain.WorkoutFindCriteria{
		TrackID: trackID,
		Limit:   3,
	})
	require.NoError(t, err)
	require.Len(t, list, 3)

	assert.Equal(t, workoutDate("2025-02-05"), list[0].Date)
	assert.Equal(t, workoutDate("2025-02-04"), list[1].Date)
	assert.Equal(t, workoutDate("2025-02-01"), list[2].Date)

	// limited list
	limited, err := ds.FindWorkouts(context.Background(), domain.WorkoutFindCriteria{
		TrackID: trackID,
		Limit:   2,
	})
	require.NoError(t, err)
	require.Len(t, limited, 2)
	assert.Equal(t, workoutDate("2025-02-05"), limited[0].Date)
	assert.Equal(t, workoutDate("2025-02-04"), limited[1].Date)

	// empty list
	list, err = ds.FindWorkouts(context.Background(), domain.WorkoutFindCriteria{
		TrackID: domain.NewTrackID(),
		Limit:   10,
	})
	require.NoError(t, err)
	require.Len(t, list, 0)
}
