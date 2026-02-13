package tracks

import (
	"net/http"

	"github.com/healthy-heroes/neskuchka/backend/app/api/httpx"
	"github.com/healthy-heroes/neskuchka/backend/app/domain"
)

// GetMainTrack returns the main track and owner flag
func (s *Service) GetMainTrack(w http.ResponseWriter, r *http.Request) {
	logger := s.logger

	track, err := s.dataStore.GetMainTrack(r.Context())
	if err != nil {
		httpx.RenderError(w, logger, http.StatusInternalServerError, err, "Failed to get main track")
		return
	}

	userID, _, _ := s.session.Get(r)
	owner := userID != "" && string(track.OwnerID) == userID

	httpx.Render(w, TrackSchema{
		Track:   TrackInfo{ID: string(track.ID), Name: track.Name},
		IsOwner: owner,
	})
}

// getMainTrackLastWorkoutsCtrl returns the exercises for the main track
func (s *Service) GetMainTrackLastWorkouts(w http.ResponseWriter, r *http.Request) {
	logger := s.logger

	track, err := s.dataStore.GetMainTrack(r.Context())
	if err != nil {
		httpx.RenderError(w, logger, http.StatusInternalServerError, err, "Failed to get main track")
		return
	}

	workouts, err := s.dataStore.FindWorkouts(r.Context(), domain.WorkoutFindCriteria{
		TrackID: domain.TrackID(track.ID),
		Limit:   10,
	})

	if err != nil {
		httpx.RenderError(w, logger, http.StatusInternalServerError, err, "Failed to get workouts")
		return
	}

	httpx.Render(w, WorkoutsSchema{
		Workouts: MakeWorkoutInfos(workouts),
	})
}
