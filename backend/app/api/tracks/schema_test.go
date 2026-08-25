package tracks

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/healthy-heroes/neskuchka/backend/app/domain"
)

// Nil slices marshal to null, and clients iterate these lists without checking.
// Both cases are real: a skill section carries no exercises, a warm-up exercise
// carries no prescription.
func TestMakeWorkoutInfo_EmptyListsAreArraysNotNull(t *testing.T) {
	workout := domain.Workout{
		ID:      domain.WorkoutID("workout-1"),
		TrackID: domain.TrackID("track-1"),
		Date:    time.Date(2026, 8, 17, 0, 0, 0, 0, time.UTC),
		Sections: []domain.WorkoutSection{
			{Title: "Навык"},
			{
				Title:     "Разминка",
				Exercises: []domain.WorkoutExercise{{Name: "Краб тач"}},
			},
		},
	}

	info := MakeWorkoutInfo(workout, time.Now())

	require.Len(t, info.Sections, 2)
	assert.NotNil(t, info.Sections[0].Exercises, "section without exercises")
	assert.NotNil(t, info.Sections[1].Exercises[0].Prescription, "exercise without prescription")

	encoded, err := json.Marshal(info)
	require.NoError(t, err)
	assert.NotContains(t, string(encoded), "null", "no null in the payload")
	assert.True(t, strings.Contains(string(encoded), `"Exercises":[]`), "empty section is an array")
	assert.True(t, strings.Contains(string(encoded), `"Prescription":[]`), "empty prescription is an array")
}

// The source workout must not be touched: handlers render it after the domain
// call and may still use it.
func TestMakeWorkoutInfo_DoesNotMutateSource(t *testing.T) {
	workout := domain.Workout{
		Sections: []domain.WorkoutSection{
			{Title: "Навык"},
		},
	}

	MakeWorkoutInfo(workout, time.Now())

	assert.Nil(t, workout.Sections[0].Exercises)
}

func TestCursor(t *testing.T) {
	t.Run("should round-trip", func(t *testing.T) {
		want := domain.WorkoutCursor{
			Date: time.Date(2026, 8, 25, 0, 0, 0, 0, time.UTC),
			ID:   domain.NewWorkoutID(),
		}

		got, err := parseCursor(makeCursor(want))

		require.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("should treat an empty cursor as the top of the list", func(t *testing.T) {
		got, err := parseCursor("")

		require.NoError(t, err)
		assert.True(t, got.IsZero())
	})

	t.Run("should render a zero cursor as empty", func(t *testing.T) {
		assert.Empty(t, makeCursor(domain.WorkoutCursor{}))
	})

	t.Run("should refuse anything malformed", func(t *testing.T) {
		// each of these used to come back as an empty page, which a client
		// reads as the end of the list rather than as its own mistake
		for _, value := range []string{
			"nonsense",
			"2026-08-25",
			"2026-08-25_",
			"2026-08-25_not-a-uuid",
			"25.08.2026_" + string(domain.NewWorkoutID()),
			"_" + string(domain.NewWorkoutID()),
		} {
			_, err := parseCursor(value)
			assert.Error(t, err, value)
		}
	})
}
