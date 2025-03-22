package datastore

import (
	"testing"

	"github.com/healthy-heroes/neskuchka/backend/app/store"
	"github.com/healthy-heroes/neskuchka/backend/app/store/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *DataStore {
	testDB, err := db.NewSqlite(":memory:")
	require.NoError(t, err)

	ds := NewDataStore(testDB)
	err = ds.Exercise.InitTables()
	require.NoError(t, err)

	return ds
}

func TestExerciseDBStore_Create(t *testing.T) {
	ds := setupTestDB(t)
	defer ds.Close()

	exercise := &store.Exercise{
		Slug: "push-up",
		Name: "Push Up",
	}

	// Test creating a new exercise
	created, err := ds.Exercise.Create(exercise)
	require.NoError(t, err)
	assert.Equal(t, exercise, created)

	// Test duplicate creation should fail
	_, err = ds.Exercise.Create(exercise)
	assert.Error(t, err, "Should fail on duplicate slug")
}

func TestExerciseDBStore_Get(t *testing.T) {
	ds := setupTestDB(t)
	defer ds.Close()

	// Create test data
	exercise := &store.Exercise{
		Slug: "squat",
		Name: "Squat",
	}
	_, err := ds.Exercise.Create(exercise)
	require.NoError(t, err)

	// Test getting an existing exercise
	found, err := ds.Exercise.Get(exercise.Slug)
	require.NoError(t, err)
	assert.Equal(t, exercise, found)

	// Test getting a non-existent exercise
	_, err = ds.Exercise.Get(store.ExerciseSlug("non-existent"))
	assert.Error(t, err, "Should error when exercise not found")
}

func TestExerciseDBStore_Find(t *testing.T) {
	ds := setupTestDB(t)
	defer ds.Close()

	// Create test data
	exercises := []*store.Exercise{
		{Slug: "push-up", Name: "Push Up"},
		{Slug: "squat", Name: "Squat"},
		{Slug: "plank", Name: "Plank"},
		{Slug: "deadlift", Name: "Deadlift"},
	}

	for _, ex := range exercises {
		_, err := ds.Exercise.Create(ex)
		require.NoError(t, err)
	}

	// Test finding all exercises
	t.Run("Find all exercises", func(t *testing.T) {
		criteria := &store.ExerciseFindCriteria{
			Limit: 5,
		}

		results, err := ds.Exercise.Find(criteria)
		require.NoError(t, err)
		assert.Len(t, results, 4, "Should find all exercises")

		// Test with limit
		criteria.Limit = 2
		results, err = ds.Exercise.Find(criteria)
		require.NoError(t, err)
		assert.Len(t, results, 2, "Should limit results")
	})

	// Test finding exercises by slugs
	t.Run("Find exercises by slugs", func(t *testing.T) {
		criteria := &store.ExerciseFindCriteria{
			Slugs: []store.ExerciseSlug{
				store.ExerciseSlug("push-up"),
				store.ExerciseSlug("plank"),
			},
			Limit: 10,
		}

		results, err := ds.Exercise.Find(criteria)
		require.NoError(t, err)
		assert.Len(t, results, 2, "Should find exercises by slugs")

		// Verify the correct exercises were returned
		slugs := make(map[store.ExerciseSlug]bool)
		for _, ex := range results {
			slugs[ex.Slug] = true
		}
		assert.True(t, slugs["push-up"], "Should find push-up exercise")
		assert.True(t, slugs["plank"], "Should find plank exercise")
	})

	// Test with non-existent slugs
	t.Run("Find with non-existent slugs", func(t *testing.T) {
		criteria := &store.ExerciseFindCriteria{
			Slugs: []store.ExerciseSlug{
				store.ExerciseSlug("non-existent"),
			},
			Limit: 10,
		}

		results, err := ds.Exercise.Find(criteria)
		require.NoError(t, err)
		assert.Len(t, results, 0, "Should return empty slice for non-existent slugs")
	})
}
