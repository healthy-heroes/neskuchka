package datastore

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/healthy-heroes/neskuchka/backend/app/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupWorkoutTestDB(t *testing.T) *DataStore {
	ds := setupTestDB(t)

	// Initialize tables with foreign key dependencies first
	err := ds.User.InitTables()
	require.NoError(t, err)

	err = ds.Exercise.InitTables()
	require.NoError(t, err)

	err = ds.Track.InitTables()
	require.NoError(t, err)

	err = ds.Workout.InitTables()
	require.NoError(t, err)

	return ds
}

func createTestTrack(t *testing.T, ds *DataStore) *store.Track {
	user := createTestUser(t, ds)

	track := &store.Track{
		ID:      store.CreateTrackId(),
		Name:    "Test Track",
		OwnerID: user.ID,
	}

	_, err := ds.Track.Create(track)
	require.NoError(t, err)

	return track
}

func createTestExercise(t *testing.T, ds *DataStore) *store.Exercise {
	exercise := &store.Exercise{
		Slug: store.ExerciseSlug(uuid.New().String()),
		Name: "Push Up",
	}

	_, err := ds.Exercise.Create(exercise)
	require.NoError(t, err)

	return exercise
}

func createTestWorkout(t *testing.T, ds *DataStore, track *store.Track) *store.Workout {
	exercise := createTestExercise(t, ds)

	now := time.Now()
	workout := &store.Workout{
		ID:      store.CreateWorkoutId(),
		Date:    now,
		TrackID: track.ID,
		Sections: []store.WorkoutSection{
			{
				Title:    "Warmup",
				Protocol: store.WorkoutProtocolDefault,
				Exercises: []store.WorkoutExercise{
					{
						ExerciseSlug:    exercise.Slug,
						Repetitions:     10,
						RepetitionsText: "10",
						Weight:          0,
						WeightText:      "body",
					},
				},
			},
		},
	}

	return workout
}

func TestWorkoutDBStore_Create(t *testing.T) {
	ds := setupWorkoutTestDB(t)
	defer ds.Close()

	track := createTestTrack(t, ds)
	workout := createTestWorkout(t, ds, track)

	// Test creating a new workout
	created, err := ds.Workout.Create(workout)
	require.NoError(t, err)
	assert.Equal(t, workout.ID, created.ID)
	assert.Equal(t, workout.TrackID, created.TrackID)
	assert.Equal(t, len(workout.Sections), len(created.Sections))
}

func TestWorkoutDBStore_Get(t *testing.T) {
	ds := setupWorkoutTestDB(t)
	defer ds.Close()

	track := createTestTrack(t, ds)
	workout := createTestWorkout(t, ds, track)

	_, err := ds.Workout.Create(workout)
	require.NoError(t, err)

	// Test getting an existing workout
	found, err := ds.Workout.Get(workout.ID)
	require.NoError(t, err)
	assert.Equal(t, workout.ID, found.ID)
	assert.Equal(t, workout.TrackID, found.TrackID)
	assert.Equal(t, len(workout.Sections), len(found.Sections))

	// Проверяем, что даты сохранены правильно
	assert.Equal(t, workout.Date.Unix(), found.Date.Unix())

	// Test section details
	assert.Equal(t, workout.Sections[0].Title, found.Sections[0].Title)
	assert.Equal(t, workout.Sections[0].Exercises[0].ExerciseSlug, found.Sections[0].Exercises[0].ExerciseSlug)

	// Test getting a non-existent workout
	nonExistentID := store.CreateWorkoutId()
	_, err = ds.Workout.Get(nonExistentID)
	assert.Error(t, err, "Should error when workout not found")
}

func TestWorkoutDBStore_GetList(t *testing.T) {
	ds := setupWorkoutTestDB(t)
	defer ds.Close()

	track := createTestTrack(t, ds)

	// Create multiple workouts
	now := time.Now()

	workout1 := createTestWorkout(t, ds, track)
	workout1.Date = now.Add(-24 * time.Hour) // Yesterday

	workout2 := createTestWorkout(t, ds, track)
	workout2.Date = now // Today

	workout3 := createTestWorkout(t, ds, track)
	workout3.Date = now.Add(-48 * time.Hour) // Day before yesterday

	_, err := ds.Workout.Create(workout1)
	require.NoError(t, err)

	_, err = ds.Workout.Create(workout2)
	require.NoError(t, err)

	_, err = ds.Workout.Create(workout3)
	require.NoError(t, err)

	// Test getting all workouts for a track
	criteria := store.WorkoutFindCriteria{
		TrackID: track.ID,
		Limit:   10,
	}

	workouts, err := ds.Workout.GetList(criteria)
	require.NoError(t, err)
	assert.Equal(t, 3, len(workouts))

	// First workout should be the most recent (today)
	assert.True(t, workouts[0].Date.After(workouts[1].Date))

	// Test with limit
	criteria.Limit = 2
	workouts, err = ds.Workout.GetList(criteria)
	require.NoError(t, err)
	assert.Equal(t, 2, len(workouts))
}
