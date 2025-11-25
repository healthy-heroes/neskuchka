package tracks

import (
	"net/http"

	R "github.com/go-pkgz/rest"
	"github.com/rs/zerolog/log"

	"github.com/healthy-heroes/neskuchka/backend/app/store"
)

// getMainTrackLastWorkoutsCtrl returns the exercises for the main track
func (s *Service) getMainTrackLastWorkouts(w http.ResponseWriter, _ *http.Request) {
	logger := log.With().Str("method", "getMainTrackLastWorkoutsCtrl").Logger()

	track, err := s.store.Track.GetMainTrack()
	if err != nil {
		logger.Error().Msgf("Failed to get main track: %s", err)
		R.RenderJSON(w, err)
		return
	}

	workouts, err := s.store.Workout.Find(&store.WorkoutFindCriteria{
		TrackID: track.ID,
		Limit:   10,
	})

	if err != nil {
		logger.Error().Msgf("Failed to get workouts: %s", err)
		R.RenderJSON(w, err)
		return
	}

	slugs := store.ExtractSlugsFromWorkouts(workouts)
	exercises, err := s.store.Exercise.Find(&store.ExerciseFindCriteria{
		Slugs: slugs,
		Limit: len(slugs),
	})
	if err != nil {
		logger.Error().Msgf("Failed to get exercises: %s", err)
		R.RenderJSON(w, err)
		return
	}

	exercisesMap := make(map[store.ExerciseSlug]store.Exercise)
	for _, exercise := range exercises {
		exercisesMap[exercise.Slug] = *exercise
	}

	response := WorkoutsSchema{
		Workouts: workouts,
	}

	R.RenderJSON(w, response)
}
