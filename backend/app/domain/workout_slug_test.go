package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWorkoutSlug(t *testing.T) {
	t.Run("should be the date and the tail of the id", func(t *testing.T) {
		w := Workout{
			ID:   WorkoutID("019a2b3c-4d5e-7f60-8a9b-0c1d2e3f4a5b"),
			Date: time.Date(2026, 9, 21, 0, 0, 0, 0, time.UTC),
		}

		assert.Equal(t, WorkoutSlug("2026-09-21-3f4a5b"), w.Slug())
	})

	t.Run("should survive a workout without an id", func(t *testing.T) {
		// Not a workout anyone can link to, but a renderer handed one must not
		// panic over it
		w := Workout{Date: time.Date(2026, 9, 21, 0, 0, 0, 0, time.UTC)}

		assert.Equal(t, WorkoutSlug("2026-09-21-"), w.Slug())
	})
}

func TestParseWorkoutSlug(t *testing.T) {
	t.Run("should decode the date and the id tail against the track", func(t *testing.T) {
		tid := NewTrackID()

		ref, err := ParseWorkoutSlug(tid, WorkoutSlug("2026-09-21-3f4a5b"))

		assert.NoError(t, err)
		assert.Equal(t, WorkoutSlugRef{
			TrackID: tid,
			Date:    time.Date(2026, 9, 21, 0, 0, 0, 0, time.UTC),
			IDTail:  "3f4a5b",
		}, ref)
	})

	t.Run("should round-trip what Slug produces", func(t *testing.T) {
		w := Workout{ID: NewWorkoutID(), TrackID: NewTrackID(), Date: day(2)}

		ref, err := ParseWorkoutSlug(w.TrackID, w.Slug())

		assert.NoError(t, err)
		assert.Equal(t, w.TrackID, ref.TrackID)
		assert.Equal(t, w.Date, ref.Date)
		assert.Equal(t, string(w.ID)[len(w.ID)-slugTailLen:], ref.IDTail)
	})

	t.Run("should reject anything else as not found", func(t *testing.T) {
		for name, slug := range map[string]WorkoutSlug{
			"empty":         "",
			"bare uuid":     WorkoutSlug(NewWorkoutID()),
			"date only":     "2026-09-21",
			"short tail":    "2026-09-21-3f4a5",
			"long tail":     "2026-09-21-3f4a5bc",
			"upper tail":    "2026-09-21-3F4A5B",
			"non-hex tail":  "2026-09-21-zzzzzz",
			"bad date":      "2026-13-41-3f4a5b",
			"bad separator": "2026-09-21_3f4a5b",
		} {
			_, err := ParseWorkoutSlug(NewTrackID(), slug)

			// Not a validation error: a malformed slug names nothing, the same as
			// a well-formed one nobody has
			assert.ErrorIs(t, err, ErrNotFound, name)
		}
	})
}
