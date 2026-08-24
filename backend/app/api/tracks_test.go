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

	IsPublished bool
	CanEdit     bool
}

type workoutRespWrapper struct {
	Workout workoutResp
}

type workoutsRespWrapper struct {
	Workouts []workoutResp
}

type workoutsPageResp struct {
	Workouts []workoutResp

	Total   int
	Planned int
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
				{Name: "Squat", Prescription: []string{"10"}},
				{Name: "Push-up", Prescription: []string{"15"}},
			},
		},
	}
}

// dateIn is a workout date relative to today, formatted the way the API takes
// and returns them.
//
// The edit window is measured from today, so anything a test writes or changes
// has to be dated from today too: a pinned calendar date passes, then quietly
// falls outside the window and takes the test with it.
func dateIn(days int) string {
	return time.Now().AddDate(0, 0, days).Format(time.DateOnly)
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
		Sections: sections("Warm-up"),
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

	t.Run("should return 404 when the main track is missing", func(t *testing.T) {
		app := NewTestApp(t)

		resp := app.GET(t, "/api/v1/tracks/main")
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
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

	t.Run("should return 404 when the main track is missing", func(t *testing.T) {
		app := NewTestApp(t)

		resp := app.GET(t, "/api/v1/tracks/main/last_workouts")
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func Test_ApiTracks_GetWorkout(t *testing.T) {
	t.Run("should return a workout by id", func(t *testing.T) {
		f := setupTracks(t)
		workout := seedWorkout(t, f.TestApp, f.Track.ID, dateIn(0), "easy day")

		resp := f.GET(t, "/api/v1/tracks/main/workouts/"+string(workout.ID))
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		data := ReadJSON[workoutRespWrapper](t, resp)
		assert.Equal(t, workoutResp{
			ID:       string(workout.ID),
			TrackID:  string(f.Track.ID),
			Date:     dateIn(0),
			Notes:    "easy day",
			Sections: sections("Warm-up"),

			IsPublished: true,
			CanEdit:     true,
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
			Date:     dateIn(1),
			Notes:    "new workout",
			Sections: sections("Main part"),
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
		assert.Equal(t, dateIn(1), data.Workout.Date)
		assert.Equal(t, "new workout", data.Workout.Notes)
		assert.Equal(t, sections("Main part"), data.Workout.Sections)

		stored, err := f.DataStorage.GetWorkout(t.Context(), domain.WorkoutRef{
			TrackID:   f.Track.ID,
			WorkoutID: domain.WorkoutID(data.Workout.ID),
		})
		require.NoError(t, err)
		assert.Equal(t, "new workout", stored.Notes)
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
		workout := seedWorkout(t, f.TestApp, f.Track.ID, dateIn(0), "before")

		resp := f.PUT(t, "/api/v1/tracks/main/workouts/"+string(workout.ID),
			WithCookie(f.LoginAs(t, f.Owner.ID)),
			WithJSON(workoutResp{
				ID:       string(workout.ID),
				TrackID:  string(f.Track.ID),
				Date:     dateIn(1),
				Notes:    "after",
				Sections: sections("Cool-down"),
			}),
		)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		data := ReadJSON[workoutRespWrapper](t, resp)
		assert.Equal(t, dateIn(1), data.Workout.Date)
		assert.Equal(t, "after", data.Workout.Notes)

		stored, err := f.DataStorage.GetWorkout(t.Context(), workout.Ref())
		require.NoError(t, err)
		assert.Equal(t, "after", stored.Notes)
		assert.Equal(t, sections("Cool-down"), stored.Sections)
	})

	t.Run("should return 401 for an anonymous user", func(t *testing.T) {
		f := setupTracks(t)
		workout := seedWorkout(t, f.TestApp, f.Track.ID, dateIn(0), "before")

		resp := f.PUT(t, "/api/v1/tracks/main/workouts/"+string(workout.ID),
			WithJSON(workoutResp{
				ID:      string(workout.ID),
				TrackID: string(f.Track.ID),
				Date:    dateIn(1),
				Notes:   "after",
			}),
		)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		stored, err := f.DataStorage.GetWorkout(t.Context(), workout.Ref())
		require.NoError(t, err)
		assert.Equal(t, "before", stored.Notes)
	})

	t.Run("should return 403 for a user who does not own the track", func(t *testing.T) {
		f := setupTracks(t)
		stranger := createUser(t, f.TestApp)
		workout := seedWorkout(t, f.TestApp, f.Track.ID, dateIn(0), "before")

		resp := f.PUT(t, "/api/v1/tracks/main/workouts/"+string(workout.ID),
			WithCookie(f.LoginAs(t, stranger.ID)),
			WithJSON(workoutResp{
				ID:      string(workout.ID),
				TrackID: string(f.Track.ID),
				Date:    dateIn(1),
				Notes:   "after",
			}),
		)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)

		stored, err := f.DataStorage.GetWorkout(t.Context(), workout.Ref())
		require.NoError(t, err)
		assert.Equal(t, "before", stored.Notes)
	})

	t.Run("should return 404 for an unknown workout", func(t *testing.T) {
		f := setupTracks(t)
		unknownID := domain.NewWorkoutID()

		resp := f.PUT(t, "/api/v1/tracks/main/workouts/"+string(unknownID),
			WithCookie(f.LoginAs(t, f.Owner.ID)),
			WithJSON(workoutResp{
				ID:      string(unknownID),
				TrackID: string(f.Track.ID),
				Date:    dateIn(1),
				Notes:   "after",
			}),
		)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func Test_ApiTracks_UpdateMainTrack(t *testing.T) {
	body := func(name, description string) map[string]string {
		return map[string]string{"Name": name, "Description": description}
	}

	t.Run("should rewrite the texts for the track owner", func(t *testing.T) {
		f := setupTracks(t)

		resp := f.PUT(t, "/api/v1/tracks/main",
			WithCookie(f.LoginAs(t, f.Owner.ID)),
			WithJSON(body("Renamed", "Rewritten")),
		)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		data := ReadJSON[trackResp](t, resp)
		assert.Equal(t, "Renamed", data.Track.Name)
		assert.True(t, data.IsOwner)

		stored, err := f.DataStorage.GetTrack(t.Context(), f.Track.ID)
		require.NoError(t, err)
		assert.Equal(t, "Renamed", stored.Name)
		assert.Equal(t, "Rewritten", stored.Description)

		// the slug is how the main track is found at all
		assert.Equal(t, f.Track.Slug, stored.Slug)
		assert.Equal(t, f.Owner.ID, stored.OwnerID)
	})

	t.Run("should return 401 for an anonymous user", func(t *testing.T) {
		f := setupTracks(t)

		resp := f.PUT(t, "/api/v1/tracks/main", WithJSON(body("Renamed", "Rewritten")))
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		stored, err := f.DataStorage.GetTrack(t.Context(), f.Track.ID)
		require.NoError(t, err)
		assert.Equal(t, f.Track.Name, stored.Name)
	})

	t.Run("should return 403 for a user who does not own the track", func(t *testing.T) {
		f := setupTracks(t)
		stranger := createUser(t, f.TestApp)

		resp := f.PUT(t, "/api/v1/tracks/main",
			WithCookie(f.LoginAs(t, stranger.ID)),
			WithJSON(body("Renamed", "Rewritten")),
		)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("should return 422 for an empty name", func(t *testing.T) {
		f := setupTracks(t)

		resp := f.PUT(t, "/api/v1/tracks/main",
			WithCookie(f.LoginAs(t, f.Owner.ID)),
			WithJSON(body("", "Rewritten")),
		)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})
}

func Test_ApiTracks_GetMainTrackWorkouts(t *testing.T) {
	t.Run("should return the whole track with its counters, drafts included", func(t *testing.T) {
		f := setupTracks(t)

		seedWorkout(t, f.TestApp, f.Track.ID, dateIn(4), "planned")
		seedWorkout(t, f.TestApp, f.Track.ID, dateIn(1), "planned too")
		seedWorkout(t, f.TestApp, f.Track.ID, dateIn(0), "today")
		seedWorkout(t, f.TestApp, f.Track.ID, dateIn(-5), "old")

		resp := f.GET(t, "/api/v1/tracks/main/workouts", WithCookie(f.LoginAs(t, f.Owner.ID)))
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		data := ReadJSON[workoutsPageResp](t, resp)
		require.Len(t, data.Workouts, 4)
		assert.Equal(t, 4, data.Total)
		assert.Equal(t, 2, data.Planned)

		// newest first, and the state of each row comes ready-made
		assert.Equal(t, dateIn(4), data.Workouts[0].Date)
		assert.False(t, data.Workouts[0].IsPublished)
		assert.True(t, data.Workouts[0].CanEdit)

		assert.True(t, data.Workouts[3].IsPublished)
		assert.False(t, data.Workouts[3].CanEdit)
	})

	t.Run("should page through the track", func(t *testing.T) {
		f := setupTracks(t)

		for day := 1; day <= 10; day++ {
			seedWorkout(t, f.TestApp, f.Track.ID, dateIn(-day), "")
		}

		cookie := f.LoginAs(t, f.Owner.ID)

		first := ReadJSON[workoutsPageResp](t,
			f.GET(t, "/api/v1/tracks/main/workouts", WithCookie(cookie)))
		require.Len(t, first.Workouts, 8)
		assert.Equal(t, 10, first.Total)
		assert.Equal(t, dateIn(-1), first.Workouts[0].Date)

		second := ReadJSON[workoutsPageResp](t,
			f.GET(t, "/api/v1/tracks/main/workouts?offset=8", WithCookie(cookie)))
		require.Len(t, second.Workouts, 2)
		assert.Equal(t, 10, second.Total)
		assert.Equal(t, dateIn(-9), second.Workouts[0].Date)
	})

	t.Run("should return an empty list, not null, for a track with no workouts", func(t *testing.T) {
		f := setupTracks(t)

		resp := f.GET(t, "/api/v1/tracks/main/workouts", WithCookie(f.LoginAs(t, f.Owner.ID)))
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		data := ReadJSON[workoutsPageResp](t, resp)
		assert.NotNil(t, data.Workouts)
		assert.Empty(t, data.Workouts)
		assert.Equal(t, 0, data.Total)
	})

	t.Run("should return 401 for an anonymous user", func(t *testing.T) {
		f := setupTracks(t)

		resp := f.GET(t, "/api/v1/tracks/main/workouts")
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 403 for a user who does not own the track", func(t *testing.T) {
		f := setupTracks(t)
		stranger := createUser(t, f.TestApp)

		resp := f.GET(t, "/api/v1/tracks/main/workouts", WithCookie(f.LoginAs(t, stranger.ID)))
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})
}

func Test_ApiTracks_DeleteWorkout(t *testing.T) {
	t.Run("should delete a workout inside the edit window", func(t *testing.T) {
		f := setupTracks(t)
		workout := seedWorkout(t, f.TestApp, f.Track.ID, dateIn(-1), "yesterday")

		resp := f.DELETE(t, "/api/v1/tracks/main/workouts/"+string(workout.ID),
			WithCookie(f.LoginAs(t, f.Owner.ID)))
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		_, err := f.DataStorage.GetWorkout(t.Context(), workout.Ref())
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})

	t.Run("should return 409 once the edit window has closed", func(t *testing.T) {
		f := setupTracks(t)
		workout := seedWorkout(t, f.TestApp, f.Track.ID, dateIn(-2), "settled")

		resp := f.DELETE(t, "/api/v1/tracks/main/workouts/"+string(workout.ID),
			WithCookie(f.LoginAs(t, f.Owner.ID)))
		assert.Equal(t, http.StatusConflict, resp.StatusCode)

		_, err := f.DataStorage.GetWorkout(t.Context(), workout.Ref())
		assert.NoError(t, err)
	})

	t.Run("should return 401 for an anonymous user", func(t *testing.T) {
		f := setupTracks(t)
		workout := seedWorkout(t, f.TestApp, f.Track.ID, dateIn(0), "today")

		resp := f.DELETE(t, "/api/v1/tracks/main/workouts/"+string(workout.ID))
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		_, err := f.DataStorage.GetWorkout(t.Context(), workout.Ref())
		assert.NoError(t, err)
	})

	t.Run("should return 403 for a user who does not own the track", func(t *testing.T) {
		f := setupTracks(t)
		stranger := createUser(t, f.TestApp)
		workout := seedWorkout(t, f.TestApp, f.Track.ID, dateIn(0), "today")

		resp := f.DELETE(t, "/api/v1/tracks/main/workouts/"+string(workout.ID),
			WithCookie(f.LoginAs(t, stranger.ID)))
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)

		_, err := f.DataStorage.GetWorkout(t.Context(), workout.Ref())
		assert.NoError(t, err)
	})

	t.Run("should return 404 for an unknown workout", func(t *testing.T) {
		f := setupTracks(t)

		resp := f.DELETE(t, "/api/v1/tracks/main/workouts/"+string(domain.NewWorkoutID()),
			WithCookie(f.LoginAs(t, f.Owner.ID)))
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func Test_ApiTracks_EditWindow(t *testing.T) {
	t.Run("should return 409 when creating a workout in the past", func(t *testing.T) {
		f := setupTracks(t)

		resp := f.POST(t, "/api/v1/tracks/main/workouts",
			WithCookie(f.LoginAs(t, f.Owner.ID)),
			WithJSON(workoutResp{
				TrackID:  string(f.Track.ID),
				Date:     dateIn(-1),
				Notes:    "backdated",
				Sections: sections("Main part"),
			}),
		)
		assert.Equal(t, http.StatusConflict, resp.StatusCode)
	})

	t.Run("should return 409 when updating a workout the window has closed on", func(t *testing.T) {
		f := setupTracks(t)
		workout := seedWorkout(t, f.TestApp, f.Track.ID, dateIn(-2), "before")

		resp := f.PUT(t, "/api/v1/tracks/main/workouts/"+string(workout.ID),
			WithCookie(f.LoginAs(t, f.Owner.ID)),
			WithJSON(workoutResp{
				ID:       string(workout.ID),
				TrackID:  string(f.Track.ID),
				Date:     dateIn(-2),
				Notes:    "after",
				Sections: sections("Main part"),
			}),
		)
		assert.Equal(t, http.StatusConflict, resp.StatusCode)

		stored, err := f.DataStorage.GetWorkout(t.Context(), workout.Ref())
		require.NoError(t, err)
		assert.Equal(t, "before", stored.Notes)
	})
}

func Test_ApiTracks_LastWorkoutsHideDrafts(t *testing.T) {
	t.Run("should not hand unpublished workouts to participants", func(t *testing.T) {
		f := setupTracks(t)

		seedWorkout(t, f.TestApp, f.Track.ID, dateIn(3), "draft")
		seedWorkout(t, f.TestApp, f.Track.ID, dateIn(0), "today")
		seedWorkout(t, f.TestApp, f.Track.ID, dateIn(-3), "old")

		// not even to the owner: this route is the participants' view of the track
		for name, opts := range map[string][]RequestOption{
			"anonymous": nil,
			"owner":     {WithCookie(f.LoginAs(t, f.Owner.ID))},
		} {
			resp := f.GET(t, "/api/v1/tracks/main/last_workouts", opts...)
			assert.Equal(t, http.StatusOK, resp.StatusCode, name)

			data := ReadJSON[workoutsRespWrapper](t, resp)
			require.Len(t, data.Workouts, 2, name)
			for _, workout := range data.Workouts {
				assert.NotEqual(t, "draft", workout.Notes, name)
				assert.True(t, workout.IsPublished, name)
			}
		}
	})
}
