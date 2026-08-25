package domain

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// day builds a workout date the way storage hands them back: midnight UTC,
// counted in days from today.
func day(offset int) time.Time {
	return dayOf(time.Now()).AddDate(0, 0, offset)
}

// workoutOn is a workout of the track dated relative to today
func workoutOn(trackID TrackID, offset int) Workout {
	w := createWorkout(trackID)
	w.Date = day(offset)

	return w
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
	setup := func(track Track, w Workout) *Store {
		return NewStore(Opts{
			Storage: &StorageStub{
				GetWorkoutFunc: func(ctx context.Context, wr WorkoutRef) (Workout, error) {
					return w, nil
				},
				GetTrackFunc: func(ctx context.Context, tid TrackID) (Track, error) {
					return track, nil
				},
			},
		})
	}

	t.Run("should return a published workout to anyone", func(t *testing.T) {
		track := createTrack()
		w := workoutOn(track.ID, 0)

		workout, err := setup(track, w).GetWorkout(context.Background(), UserID(""), w.Ref())

		assert.Nil(t, err)
		assert.Equal(t, w, workout)
	})

	t.Run("should return an unpublished workout to the track owner", func(t *testing.T) {
		track := createTrack()
		w := workoutOn(track.ID, 3)

		workout, err := setup(track, w).GetWorkout(context.Background(), track.OwnerID, w.Ref())

		assert.Nil(t, err)
		assert.Equal(t, w, workout)
	})

	t.Run("should hide an unpublished workout behind not found", func(t *testing.T) {
		track := createTrack()
		w := workoutOn(track.ID, 3)

		for name, uid := range map[string]UserID{
			"anonymous": UserID(""),
			"stranger":  NewUserID(),
		} {
			_, err := setup(track, w).GetWorkout(context.Background(), uid, w.Ref())

			// Not ErrForbidden: a 403 would confirm the id names something real
			assert.ErrorIs(t, err, ErrNotFound, name)
		}
	})
}

func TestCreateWorkout(t *testing.T) {
	t.Run("should create workout", func(t *testing.T) {
		track := createTrack()

		service := NewStore(Opts{
			Storage: &StorageStub{
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
			Storage: &StorageStub{
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
			Storage: &StorageStub{
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
			Storage: &StorageStub{
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
			Storage: &StorageStub{
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
			Storage: &StorageStub{
				FindWorkoutsFunc: func(ctx context.Context, tid TrackID, criteria WorkoutFindCriteria) ([]Workout, error) {
					usedCriteria = criteria
					return workouts, nil
				},
			},
		})
		foundWorkouts, err := service.FindWorkouts(context.Background(), track.ID, WorkoutFindCriteria{Limit: 5})
		assert.NoError(t, err)
		assert.Equal(t, workouts, foundWorkouts)
		assert.Equal(t, WorkoutFindCriteria{Limit: 5, PublishedOnly: true}, usedCriteria)
	})

	t.Run("should force published only, whatever the caller asked for", func(t *testing.T) {
		var usedCriteria WorkoutFindCriteria
		track := createTrack()
		service := NewStore(Opts{
			Storage: &StorageStub{
				FindWorkoutsFunc: func(ctx context.Context, tid TrackID, criteria WorkoutFindCriteria) ([]Workout, error) {
					usedCriteria = criteria
					return []Workout{}, nil
				},
			},
		})

		_, err := service.FindWorkouts(context.Background(), track.ID,
			WorkoutFindCriteria{Limit: 5, PublishedOnly: false})

		assert.NoError(t, err)
		assert.True(t, usedCriteria.PublishedOnly)
	})

	t.Run("should fall back on a small limit and clamp a large one", func(t *testing.T) {
		var usedLimit int
		service := NewStore(Opts{
			Storage: &StorageStub{
				FindWorkoutsFunc: func(ctx context.Context, tid TrackID, criteria WorkoutFindCriteria) ([]Workout, error) {
					usedLimit = criteria.Limit
					return []Workout{}, nil
				},
			},
		})

		// too small falls back to the default, too large is clamped to the
		// ceiling — a client paging by its own limit must not step over rows
		tcs := map[int]int{
			-1:  10,
			0:   10,
			1:   1,
			50:  50,
			51:  50,
			500: 50,
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

func TestWorkout_IsPublished(t *testing.T) {
	now := time.Now()

	tcs := map[int]bool{
		1:  false, // tomorrow is still a draft
		0:  true,  // a workout shows up on its own day
		-1: true,
	}

	for offset, expected := range tcs {
		w := Workout{Date: day(offset)}
		assert.Equal(t, expected, w.IsPublished(now), "offset %d", offset)
	}
}

func TestWorkout_IsEditable(t *testing.T) {
	now := time.Now()

	tcs := map[int]bool{
		7:  true,
		0:  true,
		-1: true,  // the day of grace
		-2: false, // people have trained by it
	}

	for offset, expected := range tcs {
		w := Workout{Date: day(offset)}
		assert.Equal(t, expected, w.IsEditable(now), "offset %d", offset)
	}
}

func TestCreateWorkout_DateWindow(t *testing.T) {
	t.Run("should refuse a workout dated past the edit window", func(t *testing.T) {
		track := createTrack()
		service := NewStore(Opts{
			Storage: &StorageStub{
				GetTrackFunc: func(ctx context.Context, tid TrackID) (Track, error) {
					return track, nil
				},
				CreateWorkoutFunc: func(ctx context.Context, w Workout) (Workout, error) {
					t.Fatal("storage should not be reached")
					return Workout{}, nil
				},
			},
		})

		_, err := service.CreateWorkout(context.Background(), track.OwnerID, workoutOn(track.ID, -2))

		assert.ErrorIs(t, err, ErrLocked)
	})

	t.Run("should accept anything the window still covers", func(t *testing.T) {
		track := createTrack()
		service := NewStore(Opts{
			Storage: &StorageStub{
				GetTrackFunc: func(ctx context.Context, tid TrackID) (Track, error) {
					return track, nil
				},
				CreateWorkoutFunc: func(ctx context.Context, w Workout) (Workout, error) {
					return w, nil
				},
			},
		})

		for _, offset := range []int{-1, 0, 3} {
			_, err := service.CreateWorkout(context.Background(), track.OwnerID, workoutOn(track.ID, offset))
			assert.NoError(t, err, "offset %d", offset)
		}
	})
}

func TestUpdateWorkout_EditWindow(t *testing.T) {
	setup := func(stored Workout, track Track) *Store {
		return NewStore(Opts{
			Storage: &StorageStub{
				GetTrackFunc: func(ctx context.Context, tid TrackID) (Track, error) {
					return track, nil
				},
				GetWorkoutFunc: func(ctx context.Context, wr WorkoutRef) (Workout, error) {
					return stored, nil
				},
				UpdateWorkoutFunc: func(ctx context.Context, w Workout) (Workout, error) {
					return w, nil
				},
			},
		})
	}

	t.Run("should refuse a workout the window has closed on", func(t *testing.T) {
		track := createTrack()
		stored := workoutOn(track.ID, -2)

		_, err := setup(stored, track).UpdateWorkout(context.Background(), track.OwnerID, stored)

		assert.ErrorIs(t, err, ErrLocked)
	})

	t.Run("should refuse moving an editable workout out of the window", func(t *testing.T) {
		track := createTrack()
		stored := workoutOn(track.ID, 0)

		update := stored
		update.Date = day(-5)

		_, err := setup(stored, track).UpdateWorkout(context.Background(), track.OwnerID, update)

		assert.ErrorIs(t, err, ErrLocked)
	})

	t.Run("should allow yesterday", func(t *testing.T) {
		track := createTrack()
		stored := workoutOn(track.ID, -1)

		updated, err := setup(stored, track).UpdateWorkout(context.Background(), track.OwnerID, stored)

		assert.NoError(t, err)
		assert.Equal(t, stored.ID, updated.ID)
	})
}

func TestDeleteWorkout(t *testing.T) {
	setup := func(track Track, stored Workout, deleted *WorkoutRef) *Store {
		return NewStore(Opts{
			Storage: &StorageStub{
				GetTrackFunc: func(ctx context.Context, tid TrackID) (Track, error) {
					return track, nil
				},
				GetWorkoutFunc: func(ctx context.Context, wr WorkoutRef) (Workout, error) {
					return stored, nil
				},
				DeleteWorkoutFunc: func(ctx context.Context, wr WorkoutRef) error {
					*deleted = wr

					return nil
				},
			},
		})
	}

	t.Run("should delete a workout inside the window", func(t *testing.T) {
		track := createTrack()
		stored := workoutOn(track.ID, -1)
		var deleted WorkoutRef

		err := setup(track, stored, &deleted).DeleteWorkout(context.Background(), track.OwnerID, stored.Ref())

		assert.NoError(t, err)
		assert.Equal(t, stored.Ref(), deleted)
	})

	t.Run("should refuse a workout the window has closed on", func(t *testing.T) {
		track := createTrack()
		stored := workoutOn(track.ID, -2)
		var deleted WorkoutRef

		err := setup(track, stored, &deleted).DeleteWorkout(context.Background(), track.OwnerID, stored.Ref())

		assert.ErrorIs(t, err, ErrLocked)
		assert.Empty(t, deleted)
	})

	t.Run("should refuse a stranger", func(t *testing.T) {
		track := createTrack()
		stored := workoutOn(track.ID, 0)
		var deleted WorkoutRef

		err := setup(track, stored, &deleted).DeleteWorkout(context.Background(), NewUserID(), stored.Ref())

		assert.ErrorIs(t, err, ErrForbidden)
		assert.Empty(t, deleted)
	})
}

func TestFindTrackWorkouts(t *testing.T) {
	setup := func(track Track, workouts []Workout, used *WorkoutFindCriteria) *Store {
		return NewStore(Opts{
			Storage: &StorageStub{
				GetTrackFunc: func(ctx context.Context, tid TrackID) (Track, error) {
					return track, nil
				},
				FindWorkoutsFunc: func(ctx context.Context, tid TrackID, criteria WorkoutFindCriteria) ([]Workout, error) {
					*used = criteria

					return workouts, nil
				},
				CountWorkoutsFunc: func(ctx context.Context, tid TrackID, now time.Time) (int, int, error) {
					return 42, 3, nil
				},
			},
		})
	}

	t.Run("should return the page with its counters, drafts included", func(t *testing.T) {
		track := createTrack()
		workouts := []Workout{workoutOn(track.ID, 2), workoutOn(track.ID, -4)}
		cursor := WorkoutCursor{Date: day(5), ID: NewWorkoutID()}
		var used WorkoutFindCriteria

		page, err := setup(track, workouts, &used).FindTrackWorkouts(
			context.Background(), track.OwnerID, track.ID,
			WorkoutFindCriteria{Limit: 8, After: cursor},
		)

		assert.NoError(t, err)
		assert.Equal(t, workouts, page.Workouts)
		assert.Equal(t, 42, page.Total)
		assert.Equal(t, 3, page.Planned)

		assert.False(t, used.PublishedOnly)
		assert.Equal(t, cursor, used.After)

		// the storage is asked for one row past the page, to learn whether
		// another one exists without arithmetic on Total
		assert.Equal(t, 9, used.Limit)

		// fewer rows came back than were asked for, so this is the last page
		assert.True(t, page.Next.IsZero())
	})

	t.Run("should hand back a cursor when another page exists", func(t *testing.T) {
		track := createTrack()
		// three rows for a page of two: the extra one is the probe
		workouts := []Workout{workoutOn(track.ID, 2), workoutOn(track.ID, 0), workoutOn(track.ID, -4)}
		var used WorkoutFindCriteria

		page, err := setup(track, workouts, &used).FindTrackWorkouts(
			context.Background(), track.OwnerID, track.ID,
			WorkoutFindCriteria{Limit: 2},
		)

		assert.NoError(t, err)
		require.Len(t, page.Workouts, 2, "the probe row is not part of the page")
		assert.Equal(t, workouts[1].ID, page.Next.ID)
		assert.Equal(t, workouts[1].Date, page.Next.Date)
	})

	t.Run("should refuse a stranger", func(t *testing.T) {
		track := createTrack()
		var used WorkoutFindCriteria

		_, err := setup(track, nil, &used).FindTrackWorkouts(
			context.Background(), NewUserID(), track.ID, WorkoutFindCriteria{},
		)

		assert.ErrorIs(t, err, ErrForbidden)
	})
}
