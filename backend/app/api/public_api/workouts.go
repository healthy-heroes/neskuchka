package public_api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	R "github.com/go-pkgz/rest"
	"github.com/healthy-heroes/neskuchka/backend/app/store"
	"github.com/rs/zerolog/log"
)

type WorkoutSchema struct {
	Workout *store.Workout
}

// getWorkoutCtrl returns a workout by id
func (api *PublicAPI) getWorkoutCtrl(w http.ResponseWriter, r *http.Request) {
	logger := log.With().Str("method", "getWorkoutCtrl").Logger()

	id := chi.URLParam(r, "id")
	workout, err := api.store.Workout.Get(store.WorkoutID(id))
	if err != nil {
		logger.Error().Msgf("Failed to get workout by id: %s", err)
		R.SendErrorJSON(w, r, nil, http.StatusNotFound, err, "Workout not found")
		return
	}

	response := WorkoutSchema{
		Workout: workout,
	}

	R.RenderJSON(w, response)
}

// updateWorkoutCtrl updates a workout
func (api *PublicAPI) updateWorkoutCtrl(w http.ResponseWriter, r *http.Request) {
	logger := log.With().Str("method", "updateWorkoutCtrl").Logger()

	id := chi.URLParam(r, "id")
	workout, err := api.store.Workout.Get(store.WorkoutID(id))
	if err != nil {
		logger.Error().Msgf("Failed to get workout by id: %s", err)
		R.SendErrorJSON(w, r, nil, http.StatusNotFound, err, "Workout not found")
		return
	}

	newWorkout := &store.Workout{}
	err = R.DecodeJSON(r, newWorkout)
	if err != nil {
		logger.Error().Msgf("Failed to decode workout: %s", err)
		R.SendErrorJSON(w, r, nil, http.StatusBadRequest, err, "Failed to decode workout")
		return
	}

	workout.Date = newWorkout.Date
	workout.Sections = newWorkout.Sections
	store.ClearUnknownExercisesSlugs(workout, map[store.ExerciseSlug]bool{})

	workout, err = api.store.Workout.Update(workout)
	if err != nil {
		logger.Error().Msgf("Failed to update workout by id: %s", err)
		R.SendErrorJSON(w, r, nil, http.StatusInternalServerError, err, "Failed to update workout")
		return
	}

	response := WorkoutSchema{
		Workout: workout,
	}

	R.RenderJSON(w, response)
}

// createWorkoutCtrl creates a new workout
func (api *PublicAPI) createWorkoutCtrl(w http.ResponseWriter, r *http.Request) {
	logger := log.With().Str("method", "createWorkoutCtrl").Logger()

	workout := &store.Workout{}
	err := R.DecodeJSON(r, workout)
	if err != nil {
		logger.Error().Msgf("Failed to decode workout: %s", err)
		R.SendErrorJSON(w, r, nil, http.StatusBadRequest, err, "Failed to decode workout")
		return
	}

	track, err := api.store.Track.GetMainTrack()
	if err != nil {
		logger.Error().Msgf("Failed to get main track: %s", err)
		R.SendErrorJSON(w, r, nil, http.StatusInternalServerError, err, "Failed to get main track")
		return
	}

	workout.ID = store.CreateWorkoutId()
	workout.TrackID = track.ID
	store.ClearUnknownExercisesSlugs(workout, map[store.ExerciseSlug]bool{})

	workout, err = api.store.Workout.Create(workout)
	if err != nil {
		logger.Error().Msgf("Failed to create workout: %s", err)
		R.SendErrorJSON(w, r, nil, http.StatusInternalServerError, err, "Failed to create workout")
		return
	}

	response := WorkoutSchema{
		Workout: workout,
	}

	R.RenderJSON(w, response)
}
