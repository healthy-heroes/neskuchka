package tracks

import (
	"net/http"

	"github.com/healthy-heroes/neskuchka/backend/app/api/httpx"
	"github.com/healthy-heroes/neskuchka/backend/app/store"
)

// GetMainTrack returns the main track and owner flag
func (s *Service) GetMainTrack(w http.ResponseWriter, r *http.Request) {
	logger := s.logger

	track, err := s.store.Track.GetMainTrack()
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
func (s *Service) GetMainTrackLastWorkouts(w http.ResponseWriter, _ *http.Request) {
	logger := s.logger

	track, err := s.store.Track.GetMainTrack()
	if err != nil {
		httpx.RenderError(w, logger, http.StatusInternalServerError, err, "Failed to get main track")
		return
	}

	workouts, err := s.store.Workout.Find(&store.WorkoutFindCriteria{
		TrackID: track.ID,
		Limit:   10,
	})

	if err != nil {
		httpx.RenderError(w, logger, http.StatusInternalServerError, err, "Failed to get workouts")
		return
	}

	slugs := store.ExtractSlugsFromWorkouts(workouts)
	exercises, err := s.store.Exercise.Find(&store.ExerciseFindCriteria{
		Slugs: slugs,
		Limit: len(slugs),
	})
	if err != nil {
		httpx.RenderError(w, logger, http.StatusInternalServerError, err, "Failed to get exercises")
		return
	}

	exercisesMap := make(map[store.ExerciseSlug]store.Exercise)
	for _, exercise := range exercises {
		exercisesMap[exercise.Slug] = *exercise
	}

	httpx.Render(w, WorkoutsSchema{
		Workouts: workouts,
	})
}
