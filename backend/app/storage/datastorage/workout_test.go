package datastorage

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/healthy-heroes/neskuchka/backend/app/domain"
	"github.com/healthy-heroes/neskuchka/backend/app/storage/db"
)

func workoutFromDB(t *testing.T, engine *db.Engine, id string) workoutRow {
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
	ds := setupTestStorage(t)

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
					{ExerciseSlug: "squat", Name: "Squat", Prescription: []string{"10 reps"}},
				},
			},
		},
	}

	createdWorkout, err := ds.CreateWorkout(t.Context(), newWorkout)
	require.NoError(t, err)
	assert.Equal(t, newWorkout.ID, createdWorkout.ID)
	assert.Equal(t, newWorkout.TrackID, createdWorkout.TrackID)
	assert.Equal(t, newWorkout.Date, createdWorkout.Date)
	assert.Equal(t, newWorkout.Notes, createdWorkout.Notes)
	assert.Equal(t, newWorkout.Sections, createdWorkout.Sections)

	workoutByID, err := ds.GetWorkout(t.Context(), domain.WorkoutRef{
		TrackID:   newWorkout.TrackID,
		WorkoutID: newWorkout.ID,
	})
	require.NoError(t, err)
	assert.Equal(t, createdWorkout, workoutByID)

	// checks system fields
	row := workoutFromDB(t, ds.engine, string(newWorkout.ID))
	assert.NotZero(t, row.CreatedAt)
	assert.NotZero(t, row.UpdatedAt)
}

func Test_Workout_Get_NotFound(t *testing.T) {
	ds := setupTestStorage(t)
	defer ds.engine.Close()

	trackID := domain.NewTrackID()
	_, err := ds.CreateWorkout(t.Context(), domain.Workout{
		ID:      domain.NewWorkoutID(),
		TrackID: trackID,
		Date:    workoutDate("2025-02-06"),
		Notes:   "Test workout notes",
		Sections: []domain.WorkoutSection{
			{Title: "Warm up", Protocol: domain.Protocol{Type: domain.ProtocolTypeCustom, Title: "Custom", Description: ""}},
		},
	})
	require.NoError(t, err)

	_, err = ds.GetWorkout(t.Context(), domain.WorkoutRef{
		TrackID:   trackID,
		WorkoutID: domain.WorkoutID("non-existent-id"),
	})
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func Test_Workout_Update(t *testing.T) {
	ds := setupTestStorage(t)

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
					{ExerciseSlug: "squat", Name: "Squat", Prescription: []string{"5x5"}},
				},
			},
		},
	}
	createdWorkout, err := ds.CreateWorkout(t.Context(), existingWorkout)
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
					{ExerciseSlug: "bench-press", Name: "Bench press", Prescription: []string{"3x10"}},
				},
			},
		},
	}
	updated, err := ds.UpdateWorkout(t.Context(), updateWorkout)
	require.NoError(t, err)
	assert.Equal(t, updateWorkout.Date, updated.Date)
	assert.Equal(t, updateWorkout.Notes, updated.Notes)
	assert.Equal(t, updateWorkout.Sections, updated.Sections)

	row := workoutFromDB(t, ds.engine, string(updated.ID))
	assert.Equal(t, createdRow.CreatedAt, row.CreatedAt)
	assert.GreaterOrEqual(t, row.UpdatedAt, createdRow.UpdatedAt)
}

func Test_Workout_FindWorkouts(t *testing.T) {
	ds := setupTestStorage(t)

	// create workouts
	trackID := domain.NewTrackID()
	dates := []string{"2025-02-01", "2025-02-05", "2025-02-04"}
	for _, date := range dates {
		_, err := ds.CreateWorkout(t.Context(), domain.Workout{
			ID:      domain.NewWorkoutID(),
			TrackID: trackID,
			Date:    workoutDate(date),
		})
		require.NoError(t, err)
	}

	_, err := ds.CreateWorkout(t.Context(), domain.Workout{
		ID:      domain.NewWorkoutID(),
		TrackID: domain.NewTrackID(),
		Date:    workoutDate("2026-02-02"),
		Notes:   "Other track",
	})
	require.NoError(t, err)

	// full list
	list, err := ds.FindWorkouts(t.Context(), trackID, domain.WorkoutFindCriteria{
		Limit: 3,
	}, time.Now())
	require.NoError(t, err)
	require.Len(t, list, 3)

	assert.Equal(t, workoutDate("2025-02-05"), list[0].Date)
	assert.Equal(t, workoutDate("2025-02-04"), list[1].Date)
	assert.Equal(t, workoutDate("2025-02-01"), list[2].Date)

	// limited list
	limited, err := ds.FindWorkouts(t.Context(), trackID, domain.WorkoutFindCriteria{
		Limit: 2,
	}, time.Now())
	require.NoError(t, err)
	require.Len(t, limited, 2)
	assert.Equal(t, workoutDate("2025-02-05"), limited[0].Date)
	assert.Equal(t, workoutDate("2025-02-04"), limited[1].Date)

	// empty list
	list, err = ds.FindWorkouts(t.Context(), domain.NewTrackID(), domain.WorkoutFindCriteria{
		Limit: 10,
	}, time.Now())
	require.NoError(t, err)
	require.Len(t, list, 0)
}

func Test_Workout_FindWorkouts_Paging(t *testing.T) {
	ds := setupTestStorage(t)

	trackID := domain.NewTrackID()
	for _, date := range []string{"2025-02-01", "2025-02-02", "2025-02-03"} {
		_, err := ds.CreateWorkout(t.Context(), domain.Workout{
			ID:      domain.NewWorkoutID(),
			TrackID: trackID,
			Date:    workoutDate(date),
		})
		require.NoError(t, err)
	}

	first, err := ds.FindWorkouts(t.Context(), trackID, domain.WorkoutFindCriteria{Limit: 2}, time.Now())
	require.NoError(t, err)
	require.Len(t, first, 2)
	assert.Equal(t, workoutDate("2025-02-03"), first[0].Date)
	assert.Equal(t, workoutDate("2025-02-02"), first[1].Date)

	// the cursor names the last row read, and the next page starts past it
	after := domain.WorkoutCursor{Date: first[1].Date, ID: first[1].ID}
	second, err := ds.FindWorkouts(t.Context(), trackID, domain.WorkoutFindCriteria{Limit: 2, After: after}, time.Now())
	require.NoError(t, err)
	require.Len(t, second, 1)
	assert.Equal(t, workoutDate("2025-02-01"), second[0].Date)

	// past the end
	end := domain.WorkoutCursor{Date: second[0].Date, ID: second[0].ID}
	past, err := ds.FindWorkouts(t.Context(), trackID, domain.WorkoutFindCriteria{Limit: 10, After: end}, time.Now())
	require.NoError(t, err)
	assert.Len(t, past, 0)
}

// Two workouts on one date must page without repeating or skipping either:
// the cursor carries the id precisely so the date alone cannot tie.
func Test_Workout_FindWorkouts_PagingTiesOnDate(t *testing.T) {
	ds := setupTestStorage(t)

	trackID := domain.NewTrackID()
	for i := 0; i < 4; i++ {
		_, err := ds.CreateWorkout(t.Context(), domain.Workout{
			ID:      domain.NewWorkoutID(),
			TrackID: trackID,
			Date:    workoutDate("2025-02-01"),
		})
		require.NoError(t, err)
	}

	seen := map[domain.WorkoutID]bool{}
	cursor := domain.WorkoutCursor{}
	for range 4 {
		page, err := ds.FindWorkouts(t.Context(), trackID,
			domain.WorkoutFindCriteria{Limit: 1, After: cursor}, time.Now())
		require.NoError(t, err)
		require.Len(t, page, 1)

		require.False(t, seen[page[0].ID], "row returned twice")
		seen[page[0].ID] = true

		cursor = domain.WorkoutCursor{Date: page[0].Date, ID: page[0].ID}
	}

	assert.Len(t, seen, 4)
}

func Test_Workout_FindWorkouts_PublishedOnly(t *testing.T) {
	ds := setupTestStorage(t)

	trackID := domain.NewTrackID()
	today := time.Now()
	dates := []time.Time{
		today.AddDate(0, 0, 3),
		today,
		today.AddDate(0, 0, -3),
	}
	for _, date := range dates {
		_, err := ds.CreateWorkout(t.Context(), domain.Workout{
			ID:      domain.NewWorkoutID(),
			TrackID: trackID,
			Date:    date,
		})
		require.NoError(t, err)
	}

	all, err := ds.FindWorkouts(t.Context(), trackID, domain.WorkoutFindCriteria{Limit: 10}, time.Now())
	require.NoError(t, err)
	assert.Len(t, all, 3)

	published, err := ds.FindWorkouts(t.Context(), trackID, domain.WorkoutFindCriteria{
		Limit:         10,
		PublishedOnly: true,
	}, time.Now())
	require.NoError(t, err)
	require.Len(t, published, 2)

	// today is published, the one three days out is not
	assert.Equal(t, today.Format(time.DateOnly), published[0].Date.Format(time.DateOnly))
}

func Test_Workout_CountWorkouts(t *testing.T) {
	ds := setupTestStorage(t)

	trackID := domain.NewTrackID()
	now := time.Now()
	for _, date := range []time.Time{now.AddDate(0, 0, 5), now.AddDate(0, 0, 1), now, now.AddDate(0, 0, -1)} {
		_, err := ds.CreateWorkout(t.Context(), domain.Workout{
			ID:      domain.NewWorkoutID(),
			TrackID: trackID,
			Date:    date,
		})
		require.NoError(t, err)
	}

	// a workout of another track must not be counted in
	_, err := ds.CreateWorkout(t.Context(), domain.Workout{
		ID:      domain.NewWorkoutID(),
		TrackID: domain.NewTrackID(),
		Date:    now,
	})
	require.NoError(t, err)

	total, planned, err := ds.CountWorkouts(t.Context(), trackID, now)
	require.NoError(t, err)

	assert.Equal(t, 4, total)
	assert.Equal(t, 2, planned)

	// an empty track counts to zero rather than failing
	total, planned, err = ds.CountWorkouts(t.Context(), domain.NewTrackID(), now)
	require.NoError(t, err)
	assert.Equal(t, 0, total)
	assert.Equal(t, 0, planned)
}

func Test_Workout_Delete(t *testing.T) {
	ds := setupTestStorage(t)

	trackID := domain.NewTrackID()
	created, err := ds.CreateWorkout(t.Context(), domain.Workout{
		ID:      domain.NewWorkoutID(),
		TrackID: trackID,
		Date:    workoutDate("2025-02-01"),
	})
	require.NoError(t, err)

	require.NoError(t, ds.DeleteWorkout(t.Context(), created.Ref()))

	_, err = ds.GetWorkout(t.Context(), created.Ref())
	assert.ErrorIs(t, err, domain.ErrNotFound)

	// deleting what is not there is a miss, not a no-op
	assert.ErrorIs(t, ds.DeleteWorkout(t.Context(), created.Ref()), domain.ErrNotFound)

	// and the workout of another track is not reachable by id alone
	other, err := ds.CreateWorkout(t.Context(), domain.Workout{
		ID:      domain.NewWorkoutID(),
		TrackID: domain.NewTrackID(),
		Date:    workoutDate("2025-02-01"),
	})
	require.NoError(t, err)

	err = ds.DeleteWorkout(t.Context(), domain.WorkoutRef{TrackID: trackID, WorkoutID: other.ID})
	assert.ErrorIs(t, err, domain.ErrNotFound)
}
