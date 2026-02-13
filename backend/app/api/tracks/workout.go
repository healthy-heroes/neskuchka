package tracks

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/healthy-heroes/neskuchka/backend/app/api/httpx"
	"github.com/healthy-heroes/neskuchka/backend/app/domain"
)

// GetWorkout returns a workout by id
func (s *Service) GetWorkout(w http.ResponseWriter, r *http.Request) {
	logger := s.logger

	id := chi.URLParam(r, "id")
	workout, err := s.dataStore.GetWorkout(r.Context(), domain.WorkoutID(id))
	if err != nil {
		httpx.RenderError(w, logger, http.StatusNotFound, err, "Workout not found")
		return
	}

	httpx.Render(w, WorkoutSchema{
		Workout: MakeWorkoutInfo(workout),
	})
}

// UpdateWorkout updates a workout
func (s *Service) UpdateWorkout(w http.ResponseWriter, r *http.Request) {
	logger := s.logger

	track, err := s.dataStore.GetMainTrack(r.Context())
	if err != nil {
		httpx.RenderError(w, logger, http.StatusInternalServerError, err, "Failed to get main track")
		return
	}
	userID, _, _ := s.session.Get(r)
	if !track.IsOwner(domain.UserID(userID)) {
		httpx.RenderError(w, logger, http.StatusForbidden, fmt.Errorf("user is not the owner of the track"), "User is not the owner of the track")
		return
	}

	newWorkout, ok := httpx.ParseBody[WorkoutCreateSchema](w, r, logger)
	if !ok {
		return
	}

	workout, err := newWorkout.toDomain()
	if err != nil {
		httpx.RenderError(w, logger, http.StatusBadRequest, err, "Failed to parse workout")
		return
	}

	workoutID := domain.WorkoutID(chi.URLParam(r, "id"))
	workout, err = s.dataStore.UpdateWorkout(r.Context(), workoutID, workout)
	if err != nil {
		httpx.RenderError(w, logger, http.StatusInternalServerError, err, "Failed to update workout")
		return
	}

	httpx.Render(w, WorkoutSchema{
		Workout: MakeWorkoutInfo(workout),
	})
}

// CreateWorkout creates a new workout
func (s *Service) CreateWorkout(w http.ResponseWriter, r *http.Request) {
	logger := s.logger

	track, err := s.dataStore.GetMainTrack(r.Context())
	if err != nil {
		httpx.RenderError(w, logger, http.StatusInternalServerError, err, "Failed to get main track")
		return
	}
	userID, _, _ := s.session.Get(r)
	if !track.IsOwner(domain.UserID(userID)) {
		httpx.RenderError(w, logger, http.StatusForbidden, fmt.Errorf("user is not the owner of the track"), "User is not the owner of the track")
		return
	}

	newWorkout, ok := httpx.ParseBody[WorkoutCreateSchema](w, r, logger)
	if !ok {
		return
	}
	workout, err := newWorkout.toDomain()
	if err != nil {
		httpx.RenderError(w, logger, http.StatusBadRequest, err, "Failed to parse workout")
		return
	}

	workout, err = s.dataStore.CreateWorkout(r.Context(), track.ID, workout)
	if err != nil {
		httpx.RenderError(w, logger, http.StatusInternalServerError, err, "Failed to create workout")
		return
	}

	httpx.Render(w, WorkoutSchema{
		Workout: MakeWorkoutInfo(workout),
	})
}
