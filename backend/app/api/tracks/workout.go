package tracks

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/healthy-heroes/neskuchka/backend/app/api/httpx"
	"github.com/healthy-heroes/neskuchka/backend/app/domain"
)

// GetWorkout returns a workout by id
func (s *Service) GetWorkout(w http.ResponseWriter, r *http.Request) {
	logger := s.logger

	// While we have only one track, we can use the main track
	// in future we will load track by slug from the path
	track, err := s.dataStore.GetMainTrack(r.Context())
	if err != nil {
		httpx.RenderDomainError(w, logger, err, "failed to get main track")
		return
	}

	id := chi.URLParam(r, "id")
	workout, err := s.dataStore.GetWorkout(r.Context(), domain.WorkoutRef{TrackID: track.ID, WorkoutID: domain.WorkoutID(id)})
	if err != nil {
		httpx.RenderDomainError(w, logger, err, "failed to get workout")
		return
	}
	httpx.Render(w, WorkoutSchema{Workout: MakeWorkoutInfo(workout)})
}

// UpdateWorkout updates a workout
func (s *Service) UpdateWorkout(w http.ResponseWriter, r *http.Request) {
	logger := s.logger

	payload, ok := httpx.ParseBody[WorkoutInfo](w, r, logger)
	if !ok {
		return
	}

	workout, err := payload.toDomain()
	if err != nil {
		httpx.RenderError(w, logger, http.StatusBadRequest, err, "failed to parse workout")
		return
	}

	userID, _, _ := s.session.Get(r)
	workout, err = s.dataStore.UpdateWorkout(r.Context(), domain.UserID(userID), workout)
	if err != nil {
		httpx.RenderDomainError(w, logger, err, "failed to update workout")
		return
	}

	httpx.Render(w, WorkoutSchema{
		Workout: MakeWorkoutInfo(workout),
	})
}

// CreateWorkout creates a new workout
func (s *Service) CreateWorkout(w http.ResponseWriter, r *http.Request) {
	logger := s.logger

	payload, ok := httpx.ParseBody[WorkoutInfo](w, r, logger)
	if !ok {
		return
	}

	workout, err := payload.toDomain()
	if err != nil {
		httpx.RenderError(w, logger, http.StatusBadRequest, err, "failed to parse workout")
		return
	}

	userID, _, _ := s.session.Get(r)
	workout, err = s.dataStore.CreateWorkout(r.Context(), domain.UserID(userID), workout)
	if err != nil {
		httpx.RenderDomainError(w, logger, err, "failed to create workout")
		return
	}

	httpx.Render(w, WorkoutSchema{
		Workout: MakeWorkoutInfo(workout),
	})
}
