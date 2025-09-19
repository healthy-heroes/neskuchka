package public_api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	R "github.com/go-pkgz/rest"
	"github.com/healthy-heroes/neskuchka/backend/app/store"
	"github.com/rs/zerolog/log"
)

type WorkoutSchema struct {
	Workout *store.Workout
}

// getWorkoutCtrl returns a workout by id
func (api *PublicAPI) getWorkoutCtrl(w http.ResponseWriter, r *http.Request) {
	logger := log.With().Str("method", "getWorkoutCtrl").Logger()

	id := chi.URLParam(r, "id")
	workout, err := api.store.Workout.Get(store.WorkoutID(id))
	if err != nil {
		logger.Error().Msgf("Failed to get workout by id: %s", err)
		R.RenderJSON(w, err)
		return
	}

	response := WorkoutSchema{
		Workout: workout,
	}

	R.RenderJSON(w, response)
}
