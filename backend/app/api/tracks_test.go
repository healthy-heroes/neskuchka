package api

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/healthy-heroes/neskuchka/backend/app/domain"
	"github.com/healthy-heroes/neskuchka/backend/app/internal/testutil"
)

// wire shapes of the tracks API, declared here on purpose: they pin the JSON
// contract the frontend relies on, so a rename in the handler schema fails a test
type trackResp struct {
	Track struct {
		ID   string
		Name string
	}
	IsOwner bool
}

type workoutResp struct {
	ID       string
	TrackID  string
	Date     string
	Notes    string
	Sections []domain.WorkoutSection
}

type workoutRespWrapper struct {
	Workout workoutResp
}

type workoutsRespWrapper struct {
	Workouts []workoutResp
}

// tracksFixture is an app with the main track and its owner already created
type tracksFixture struct {
	*TestApp

	Owner domain.User
	Track domain.Track
}

func setupTracks(t *testing.T) tracksFixture {
	t.Helper()

	app := NewTestApp(t)
	owner := createUser(t, app)

	track, err := app.DataStorage.CreateTrack(t.Context(), domain.Track{
		ID:          domain.NewTrackID(),
		Slug:        domain.TrackSlug("main"),
		Name:        "Main track",
		Description: "Track for tests",
		OwnerID:     owner.ID,
	})
	require.NoError(t, err)

	return tracksFixture{TestApp: app, Owner: owner, Track: track}
}

func createUser(t *testing.T, app *TestApp) domain.User {
	t.Helper()

	user, err := app.DataStorage.CreateUser(t.Context(), testutil.CreateUser())
	require.NoError(t, err)

	return user
}

// sections builds a workout body with the exercise slugs left empty:
// the domain clears them on write, so fixtures compare to responses as is
func sections(title string) []domain.WorkoutSection {
	return []domain.WorkoutSection{
		{
			Title:    title,
			Protocol: domain.Protocol{Type: domain.ProtocolTypeCustom, Title: "AMRAP 12"},
			Exercises: []domain.WorkoutExercise{
				{Description: "10 приседаний"},
				{Description: "15 отжиманий"},
			},
		},
	}
}

func seedWorkout(t *testing.T, app *TestApp, trackID domain.TrackID, date, notes string) domain.Workout {
	t.Helper()

	parsedDate, err := time.Parse(time.DateOnly, date)
	require.NoError(t, err)

	workout, err := app.DataStorage.CreateWorkout(t.Context(), domain.Workout{
		ID:       domain.NewWorkoutID(),
		TrackID:  trackID,
		Date:     parsedDate,
		Notes:    notes,
		Sections: sections("Разминка"),
	})
	require.NoError(t, err)

	return workout
}

func Test_ApiTracks_GetMainTrack(t *testing.T) {
	t.Run("should return the track to an anonymous user without the owner flag", func(t *testing.T) {
		f := setupTracks(t)

		resp := f.GET(t, "/api/v1/tracks/main")
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		data := ReadJSON[trackResp](t, resp)
		assert.Equal(t, string(f.Track.ID), data.Track.ID)
		assert.Equal(t, f.Track.Name, data.Track.Name)
		assert.False(t, data.IsOwner)
	})

	t.Run("should set the owner flag for the track owner", func(t *testing.T) {
		f := setupTracks(t)

		resp := f.GET(t, "/api/v1/tracks/main", WithCookie(f.LoginAs(t, f.Owner.ID)))
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		data := ReadJSON[trackResp](t, resp)
		assert.True(t, data.IsOwner)
	})

	t.Run("should not set the owner flag for another logged in user", func(t *testing.T) {
		f := setupTracks(t)
		stranger := createUser(t, f.TestApp)

		resp := f.GET(t, "/api/v1/tracks/main", WithCookie(f.LoginAs(t, stranger.ID)))
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		data := ReadJSON[trackResp](t, resp)
		assert.False(t, data.IsOwner)
	})
}

func Test_ApiTracks_GetMainTrackLastWorkouts(t *testing.T) {
	t.Run("should return an empty list, not null, when the track has no workouts", func(t *testing.T) {
		f := setupTracks(t)

		resp := f.GET(t, "/api/v1/tracks/main/last_workouts")
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		data := ReadJSON[workoutsRespWrapper](t, resp)
		assert.NotNil(t, data.Workouts)
		assert.Empty(t, data.Workouts)
	})

	t.Run("should return workouts newest first", func(t *testing.T) {
		f := setupTracks(t)

		seedWorkout(t, f.TestApp, f.Track.ID, "2026-01-10", "oldest")
		seedWorkout(t, f.TestApp, f.Track.ID, "2026-03-10", "newest")
		seedWorkout(t, f.TestApp, f.Track.ID, "2026-02-10", "middle")

		resp := f.GET(t, "/api/v1/tracks/main/last_workouts")
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		data := ReadJSON[workoutsRespWrapper](t, resp)
		require.Len(t, data.Workouts, 3)
		assert.Equal(t, []string{"2026-03-10", "2026-02-10", "2026-01-10"}, []string{
			data.Workouts[0].Date, data.Workouts[1].Date, data.Workouts[2].Date,
		})
	})

	t.Run("should return at most 10 workouts", func(t *testing.T) {
		f := setupTracks(t)

		for day := 1; day <= 12; day++ {
			seedWorkout(t, f.TestApp, f.Track.ID, fmt.Sprintf("2026-01-%02d", day), "")
		}

		resp := f.GET(t, "/api/v1/tracks/main/last_workouts")
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		data := ReadJSON[workoutsRespWrapper](t, resp)
		require.Len(t, data.Workouts, 10)
		assert.Equal(t, "2026-01-12", data.Workouts[0].Date)
		assert.Equal(t, "2026-01-03", data.Workouts[9].Date)
	})

	t.Run("should not return workouts of another track", func(t *testing.T) {
		f := setupTracks(t)

		otherTrack, err := f.DataStorage.CreateTrack(t.Context(), domain.Track{
			ID:      domain.NewTrackID(),
			Slug:    domain.TrackSlug("other"),
			Name:    "Other track",
			OwnerID: f.Owner.ID,
		})
		require.NoError(t, err)

		seedWorkout(t, f.TestApp, f.Track.ID, "2026-01-10", "mine")
		seedWorkout(t, f.TestApp, otherTrack.ID, "2026-01-11", "alien")

		resp := f.GET(t, "/api/v1/tracks/main/last_workouts")
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		data := ReadJSON[workoutsRespWrapper](t, resp)
		require.Len(t, data.Workouts, 1)
		assert.Equal(t, "mine", data.Workouts[0].Notes)
	})
}

func Test_ApiTracks_GetWorkout(t *testing.T) {
	t.Run("should return a workout by id", func(t *testing.T) {
		f := setupTracks(t)
		workout := seedWorkout(t, f.TestApp, f.Track.ID, "2026-01-10", "легкий день")

		resp := f.GET(t, "/api/v1/tracks/main/workouts/"+string(workout.ID))
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		data := ReadJSON[workoutRespWrapper](t, resp)
		assert.Equal(t, workoutResp{
			ID:       string(workout.ID),
			TrackID:  string(f.Track.ID),
			Date:     "2026-01-10",
			Notes:    "легкий день",
			Sections: sections("Разминка"),
		}, data.Workout)
	})

	t.Run("should return 404 for an unknown workout", func(t *testing.T) {
		f := setupTracks(t)

		resp := f.GET(t, "/api/v1/tracks/main/workouts/"+string(domain.NewWorkoutID()))
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should return 404 for a workout of another track", func(t *testing.T) {
		f := setupTracks(t)

		otherTrack, err := f.DataStorage.CreateTrack(t.Context(), domain.Track{
			ID:      domain.NewTrackID(),
			Slug:    domain.TrackSlug("other"),
			Name:    "Other track",
			OwnerID: f.Owner.ID,
		})
		require.NoError(t, err)

		alien := seedWorkout(t, f.TestApp, otherTrack.ID, "2026-01-10", "alien")

		resp := f.GET(t, "/api/v1/tracks/main/workouts/"+string(alien.ID))
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func Test_ApiTracks_CreateWorkout(t *testing.T) {
	newWorkoutBody := func(trackID domain.TrackID) workoutResp {
		return workoutResp{
			TrackID:  string(trackID),
			Date:     "2026-01-10",
			Notes:    "новая тренировка",
			Sections: sections("Основная часть"),
		}
	}

	t.Run("should create a workout for the track owner", func(t *testing.T) {
		f := setupTracks(t)

		resp := f.POST(t, "/api/v1/tracks/main/workouts",
			WithCookie(f.LoginAs(t, f.Owner.ID)),
			WithJSON(newWorkoutBody(f.Track.ID)),
		)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		data := ReadJSON[workoutRespWrapper](t, resp)
		assert.NotEmpty(t, data.Workout.ID)
		assert.Equal(t, "2026-01-10", data.Workout.Date)
		assert.Equal(t, "новая тренировка", data.Workout.Notes)
		assert.Equal(t, sections("Основная часть"), data.Workout.Sections)

		stored, err := f.DataStorage.GetWorkout(t.Context(), domain.WorkoutRef{
			TrackID:   f.Track.ID,
			WorkoutID: domain.WorkoutID(data.Workout.ID),
		})
		require.NoError(t, err)
		assert.Equal(t, "новая тренировка", stored.Notes)
	})

	t.Run("should ignore the id sent by the client and generate its own", func(t *testing.T) {
		f := setupTracks(t)

		body := newWorkoutBody(f.Track.ID)
		body.ID = "id-from-client"

		resp := f.POST(t, "/api/v1/tracks/main/workouts",
			WithCookie(f.LoginAs(t, f.Owner.ID)),
			WithJSON(body),
		)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		data := ReadJSON[workoutRespWrapper](t, resp)
		assert.NotEmpty(t, data.Workout.ID)
		assert.NotEqual(t, "id-from-client", data.Workout.ID)
	})

	t.Run("should return 401 for an anonymous user", func(t *testing.T) {
		f := setupTracks(t)

		resp := f.POST(t, "/api/v1/tracks/main/workouts", WithJSON(newWorkoutBody(f.Track.ID)))
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 403 for a user who does not own the track", func(t *testing.T) {
		f := setupTracks(t)
		stranger := createUser(t, f.TestApp)

		resp := f.POST(t, "/api/v1/tracks/main/workouts",
			WithCookie(f.LoginAs(t, stranger.ID)),
			WithJSON(newWorkoutBody(f.Track.ID)),
		)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("should return 400 for a malformed date", func(t *testing.T) {
		f := setupTracks(t)

		body := newWorkoutBody(f.Track.ID)
		body.Date = "10.01.2026"

		resp := f.POST(t, "/api/v1/tracks/main/workouts",
			WithCookie(f.LoginAs(t, f.Owner.ID)),
			WithJSON(body),
		)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func Test_ApiTracks_UpdateWorkout(t *testing.T) {
	t.Run("should update the workout of the track owner", func(t *testing.T) {
		f := setupTracks(t)
		workout := seedWorkout(t, f.TestApp, f.Track.ID, "2026-01-10", "было")

		resp := f.PUT(t, "/api/v1/tracks/main/workouts/"+string(workout.ID),
			WithCookie(f.LoginAs(t, f.Owner.ID)),
			WithJSON(workoutResp{
				ID:       string(workout.ID),
				TrackID:  string(f.Track.ID),
				Date:     "2026-01-11",
				Notes:    "стало",
				Sections: sections("Заминка"),
			}),
		)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		data := ReadJSON[workoutRespWrapper](t, resp)
		assert.Equal(t, "2026-01-11", data.Workout.Date)
		assert.Equal(t, "стало", data.Workout.Notes)

		stored, err := f.DataStorage.GetWorkout(t.Context(), workout.Ref())
		require.NoError(t, err)
		assert.Equal(t, "стало", stored.Notes)
		assert.Equal(t, sections("Заминка"), stored.Sections)
	})

	t.Run("should return 401 for an anonymous user", func(t *testing.T) {
		f := setupTracks(t)
		workout := seedWorkout(t, f.TestApp, f.Track.ID, "2026-01-10", "было")

		resp := f.PUT(t, "/api/v1/tracks/main/workouts/"+string(workout.ID),
			WithJSON(workoutResp{
				ID:      string(workout.ID),
				TrackID: string(f.Track.ID),
				Date:    "2026-01-11",
				Notes:   "стало",
			}),
		)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		stored, err := f.DataStorage.GetWorkout(t.Context(), workout.Ref())
		require.NoError(t, err)
		assert.Equal(t, "было", stored.Notes)
	})

	t.Run("should return 403 for a user who does not own the track", func(t *testing.T) {
		f := setupTracks(t)
		stranger := createUser(t, f.TestApp)
		workout := seedWorkout(t, f.TestApp, f.Track.ID, "2026-01-10", "было")

		resp := f.PUT(t, "/api/v1/tracks/main/workouts/"+string(workout.ID),
			WithCookie(f.LoginAs(t, stranger.ID)),
			WithJSON(workoutResp{
				ID:      string(workout.ID),
				TrackID: string(f.Track.ID),
				Date:    "2026-01-11",
				Notes:   "стало",
			}),
		)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)

		stored, err := f.DataStorage.GetWorkout(t.Context(), workout.Ref())
		require.NoError(t, err)
		assert.Equal(t, "было", stored.Notes)
	})

	t.Run("should return 404 for an unknown workout", func(t *testing.T) {
		f := setupTracks(t)
		unknownID := domain.NewWorkoutID()

		resp := f.PUT(t, "/api/v1/tracks/main/workouts/"+string(unknownID),
			WithCookie(f.LoginAs(t, f.Owner.ID)),
			WithJSON(workoutResp{
				ID:      string(unknownID),
				TrackID: string(f.Track.ID),
				Date:    "2026-01-11",
				Notes:   "стало",
			}),
		)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
