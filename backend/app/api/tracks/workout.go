package tracks

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	R "github.com/go-pkgz/rest"
	"github.com/healthy-heroes/neskuchka/backend/app/store"
	"github.com/rs/zerolog/log"
)

// getWorkout returns a workout by id
func (s *Service) getWorkout(w http.ResponseWriter, r *http.Request) {
	logger := log.With().Str("method", "getWorkout").Logger()

	id := chi.URLParam(r, "id")
	workout, err := s.store.Workout.Get(store.WorkoutID(id))
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

// updateWorkout updates a workout
func (s *Service) updateWorkout(w http.ResponseWriter, r *http.Request) {
	logger := log.With().Str("method", "updateWorkout").Logger()

	id := chi.URLParam(r, "id")
	workout, err := s.store.Workout.Get(store.WorkoutID(id))
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
	workout.Notes = newWorkout.Notes
	store.ClearUnknownExercisesSlugs(workout, map[store.ExerciseSlug]bool{})

	workout, err = s.store.Workout.Update(workout)
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

// createWorkout creates a new workout
func (s *Service) createWorkout(w http.ResponseWriter, r *http.Request) {
	logger := log.With().Str("method", "createWorkout").Logger()

	workout := &store.Workout{}
	err := R.DecodeJSON(r, workout)
	if err != nil {
		logger.Error().Msgf("Failed to decode workout: %s", err)
		R.SendErrorJSON(w, r, nil, http.StatusBadRequest, err, "Failed to decode workout")
		return
	}

	track, err := s.store.Track.GetMainTrack()
	if err != nil {
		logger.Error().Msgf("Failed to get main track: %s", err)
		R.SendErrorJSON(w, r, nil, http.StatusInternalServerError, err, "Failed to get main track")
		return
	}

	workout.ID = store.CreateWorkoutId()
	workout.TrackID = track.ID
	store.ClearUnknownExercisesSlugs(workout, map[store.ExerciseSlug]bool{})

	workout, err = s.store.Workout.Create(workout)
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
