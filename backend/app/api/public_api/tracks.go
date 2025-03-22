package public_api

import (
	"fmt"
	"net/http"

	R "github.com/go-pkgz/rest"
	"github.com/healthy-heroes/neskuchka/backend/app/store"
	"github.com/rs/zerolog/log"
)

type TrackWorkoutsSchema struct {
	Workouts  []*store.Workout
	Exercises map[store.ExerciseSlug]store.Exercise
}

// getMainTrackLastExercisesCtrl returns the exercises for the main track
func (api *PublicAPI) getMainTrackLastExercisesCtrl(w http.ResponseWriter, _ *http.Request) {
	logger := log.With().Str("method", "getMainTrackLastExercisesCtrl").Logger()

	track, err := api.store.Track.GetMainTrack()
	if err != nil {
		logger.Error().Msgf("Failed to get main track: %s", err)
		R.RenderJSON(w, err)
		return
	}

	workouts, err := api.store.Workout.Find(&store.WorkoutFindCriteria{
		TrackID: track.ID,
		Limit:   10,
	})

	if err != nil {
		logger.Error().Msgf("Failed to get workouts: %s", err)
		R.RenderJSON(w, err)
		return
	}

	slugs := store.ExtractSlugsFromWorkouts(workouts)
	exercises, err := api.store.Exercise.Find(&store.ExerciseFindCriteria{
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

	response := TrackWorkoutsSchema{
		Workouts:  workouts,
		Exercises: exercisesMap,
	}

	fmt.Printf("response: %+v\n", response)

	R.RenderJSON(w, response)
}
