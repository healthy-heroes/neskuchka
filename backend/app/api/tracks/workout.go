package tracks

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/healthy-heroes/neskuchka/backend/app/api/httpx"
	"github.com/healthy-heroes/neskuchka/backend/app/domain"
	"github.com/healthy-heroes/neskuchka/backend/app/internal/session"
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

	// Anonymous readers have no user id, and that is fine: the domain only needs
	// it to decide whether an unpublished workout may be shown.
	userID, _ := session.GetUserID(r)

	id := chi.URLParam(r, "id")
	workout, err := s.dataStore.GetWorkout(r.Context(), domain.UserID(userID),
		domain.WorkoutRef{TrackID: track.ID, WorkoutID: domain.WorkoutID(id)})
	if err != nil {
		httpx.RenderDomainError(w, logger, err, "failed to get workout")
		return
	}
	httpx.Render(w, WorkoutSchema{Workout: MakeWorkoutInfo(workout)})
}

// UpdateWorkout updates a workout
func (s *Service) UpdateWorkout(w http.ResponseWriter, r *http.Request) {
	logger := s.logger
	userID := session.MustGetUserID(r)

	payload, ok := httpx.ParseBody[WorkoutInfo](w, r, logger)
	if !ok {
		return
	}

	workout, err := payload.toDomain()
	if err != nil {
		httpx.RenderError(w, logger, http.StatusBadRequest, err, "failed to parse workout")
		return
	}

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
	userID := session.MustGetUserID(r)

	payload, ok := httpx.ParseBody[WorkoutInfo](w, r, logger)
	if !ok {
		return
	}

	workout, err := payload.toDomain()
	if err != nil {
		httpx.RenderError(w, logger, http.StatusBadRequest, err, "failed to parse workout")
		return
	}

	workout, err = s.dataStore.CreateWorkout(r.Context(), domain.UserID(userID), workout)
	if err != nil {
		httpx.RenderDomainError(w, logger, err, "failed to create workout")
		return
	}

	httpx.Render(w, WorkoutSchema{
		Workout: MakeWorkoutInfo(workout),
	})
}

// managePageSize is how many rows the track management screen asks for at a time.
const managePageSize = 8

// GetMainTrackWorkouts returns a page of the whole track — drafts included —
// to its owner. Participants read the track through GetMainTrackLastWorkouts,
// which never shows them anything unpublished.
func (s *Service) GetMainTrackWorkouts(w http.ResponseWriter, r *http.Request) {
	logger := s.logger
	userID := session.MustGetUserID(r)

	track, err := s.dataStore.GetMainTrack(r.Context())
	if err != nil {
		httpx.RenderDomainError(w, logger, err, "Failed to get main track")
		return
	}

	after, err := parseCursor(r.URL.Query().Get("after"))
	if err != nil {
		httpx.RenderError(w, logger, http.StatusBadRequest, err, "Failed to parse cursor")
		return
	}

	page, err := s.dataStore.FindTrackWorkouts(r.Context(), domain.UserID(userID), track.ID,
		domain.WorkoutFindCriteria{
			Limit: httpx.QueryInt(r, "limit", managePageSize),
			After: after,
		},
	)
	if err != nil {
		httpx.RenderDomainError(w, logger, err, "Failed to get workouts")
		return
	}

	httpx.Render(w, WorkoutsPageSchema{
		Workouts:   MakeWorkoutInfos(page.Workouts),
		NextCursor: makeCursor(page.Next),
		Total:      page.Total,
		Planned:    page.Planned,
	})
}

// DeleteWorkout removes a workout from the main track
func (s *Service) DeleteWorkout(w http.ResponseWriter, r *http.Request) {
	logger := s.logger
	userID := session.MustGetUserID(r)

	track, err := s.dataStore.GetMainTrack(r.Context())
	if err != nil {
		httpx.RenderDomainError(w, logger, err, "failed to get main track")
		return
	}

	id := chi.URLParam(r, "id")
	ref := domain.WorkoutRef{TrackID: track.ID, WorkoutID: domain.WorkoutID(id)}

	if err := s.dataStore.DeleteWorkout(r.Context(), domain.UserID(userID), ref); err != nil {
		httpx.RenderDomainError(w, logger, err, "failed to delete workout")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
