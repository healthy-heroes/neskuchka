package public_api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	R "github.com/go-pkgz/rest"
	"github.com/rs/zerolog/log"

	"github.com/healthy-heroes/neskuchka/backend/app/store"
)

// getExerciseCtrl returns an exercise by slug
func (api *PublicAPI) getExerciseCtrl(w http.ResponseWriter, r *http.Request) {
	logger := log.With().Str("method", "getExerciseCtrl").Logger()

	slugParam := chi.URLParam(r, "slug")
	slug := store.ExerciseSlug(slugParam)

	exercise, err := api.store.Exercise.Get(slug)
	if err != nil {
		logger.Error().Msgf("Failed to get exercise by slug: %s", err)

		R.RenderJSON(w, err)
		return
	}

	R.RenderJSON(w, exercise)

}
