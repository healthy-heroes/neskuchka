package tracks

import (
	"net/http"

	"github.com/healthy-heroes/neskuchka/backend/app/api/httpx"
	"github.com/healthy-heroes/neskuchka/backend/app/domain"
	"github.com/healthy-heroes/neskuchka/backend/app/internal/session"
)

// GetMainTrack returns the main track and owner flag
func (s *Service) GetMainTrack(w http.ResponseWriter, r *http.Request) {
	logger := s.logger
	userID, _ := session.GetUserID(r)

	track, err := s.dataStore.GetMainTrack(r.Context())
	if err != nil {
		httpx.RenderDomainError(w, logger, err, "Failed to get main track")
		return
	}

	// The author is shown on the track page. A missing owner should not take the
	// whole page down, so fall back to an empty author instead of failing.
	author, err := s.dataStore.GetUser(r.Context(), track.OwnerID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get track owner")
		author = domain.User{}
	}

	httpx.Render(w, TrackSchema{
		Track:   MakeTrackInfo(track, author),
		IsOwner: track.IsOwner(domain.UserID(userID)),
	})
}

// getMainTrackLastWorkoutsCtrl returns the exercises for the main track
func (s *Service) GetMainTrackLastWorkouts(w http.ResponseWriter, r *http.Request) {
	logger := s.logger

	track, err := s.dataStore.GetMainTrack(r.Context())
	if err != nil {
		httpx.RenderDomainError(w, logger, err, "Failed to get main track")
		return
	}

	workouts, err := s.dataStore.FindWorkouts(r.Context(), track.ID, domain.WorkoutFindCriteria{Limit: 10})
	if err != nil {
		httpx.RenderError(w, logger, http.StatusInternalServerError, err, "Failed to get workouts")
		return
	}

	httpx.Render(w, WorkoutsSchema{
		Workouts: MakeWorkoutInfos(workouts),
	})
}

// UpdateMainTrack rewrites the main track texts on behalf of its owner
func (s *Service) UpdateMainTrack(w http.ResponseWriter, r *http.Request) {
	logger := s.logger
	userID := session.MustGetUserID(r)

	payload, ok := httpx.ParseAndValidateBody[UpdateTrackSchema](w, r, logger)
	if !ok {
		return
	}

	track, err := s.dataStore.GetMainTrack(r.Context())
	if err != nil {
		httpx.RenderDomainError(w, logger, err, "Failed to get main track")
		return
	}

	track, err = s.dataStore.UpdateTrack(r.Context(), domain.UserID(userID), payload.toDomain(track.ID))
	if err != nil {
		httpx.RenderDomainError(w, logger, err, "Failed to update main track")
		return
	}

	// Same shape as GetMainTrack, and for the same reason: the page redraws its
	// header from the response, author line included.
	author, err := s.dataStore.GetUser(r.Context(), track.OwnerID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get track owner")
		author = domain.User{}
	}

	httpx.Render(w, TrackSchema{
		Track:   MakeTrackInfo(track, author),
		IsOwner: true,
	})
}
