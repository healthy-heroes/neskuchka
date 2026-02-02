package tracks

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/healthy-heroes/neskuchka/backend/app/api/httpx"
	"github.com/healthy-heroes/neskuchka/backend/app/store"
)

// getWorkout returns a workout by id
func (s *Service) GetWorkout(w http.ResponseWriter, r *http.Request) {
	logger := s.logger

	id := chi.URLParam(r, "id")
	workout, err := s.store.Workout.Get(store.WorkoutID(id))
	if err != nil {
		httpx.RenderError(w, logger, http.StatusNotFound, err, "Workout not found")
		return
	}

	httpx.Render(w, WorkoutSchema{
		Workout: workout,
	})
}

// updateWorkout updates a workout
func (s *Service) UpdateWorkout(w http.ResponseWriter, r *http.Request) {
	logger := s.logger

	track, err := s.store.Track.GetMainTrack()
	if err != nil {
		httpx.RenderError(w, logger, http.StatusInternalServerError, err, "Failed to get main track")
		return
	}
	userID, _, _ := s.session.Get(r)
	if !track.IsOwner(store.UserID(userID)) {
		httpx.RenderError(w, logger, http.StatusForbidden, fmt.Errorf("user is not the owner of the track"), "User is not the owner of the track")
		return
	}

	id := chi.URLParam(r, "id")
	workout, err := s.store.Workout.Get(store.WorkoutID(id))
	if err != nil {
		httpx.RenderError(w, logger, http.StatusNotFound, err, "Workout not found")
		return
	}

	newWorkout, ok := httpx.ParseBody[*store.Workout](w, r, logger)
	if !ok {
		return
	}

	workout.Date = newWorkout.Date
	workout.Sections = newWorkout.Sections
	workout.Notes = newWorkout.Notes
	workout.ClearUnknownExercisesSlugs(map[store.ExerciseSlug]bool{})

	workout, err = s.store.Workout.Update(workout)
	if err != nil {
		httpx.RenderError(w, logger, http.StatusInternalServerError, err, "Failed to update workout")
		return
	}

	httpx.Render(w, WorkoutSchema{
		Workout: workout,
	})
}

// createWorkout creates a new workout
func (s *Service) CreateWorkout(w http.ResponseWriter, r *http.Request) {
	logger := s.logger

	track, err := s.store.Track.GetMainTrack()
	if err != nil {
		httpx.RenderError(w, logger, http.StatusInternalServerError, err, "Failed to get main track")
		return
	}
	userID, _, _ := s.session.Get(r)
	if !track.IsOwner(store.UserID(userID)) {
		httpx.RenderError(w, logger, http.StatusForbidden, fmt.Errorf("user is not the owner of the track"), "User is not the owner of the track")
		return
	}

	workout, ok := httpx.ParseBody[*store.Workout](w, r, logger)
	if !ok {
		return
	}

	workout.ID = store.CreateWorkoutId()
	workout.TrackID = track.ID
	workout.ClearUnknownExercisesSlugs(map[store.ExerciseSlug]bool{})

	workout, err = s.store.Workout.Create(workout)
	if err != nil {
		httpx.RenderError(w, logger, http.StatusInternalServerError, err, "Failed to create workout")
		return
	}

	httpx.Render(w, WorkoutSchema{
		Workout: workout,
	})
}
