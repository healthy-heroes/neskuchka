package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	R "github.com/go-pkgz/rest"
	"github.com/rs/zerolog/log"

	"github.com/healthy-heroes/neskuchka/backend/app/store"
	"github.com/healthy-heroes/neskuchka/backend/app/store/datastore"
)

type PublicMethods struct {
	store *datastore.DataStore
}

func (pm *PublicMethods) pingCtrl(w http.ResponseWriter, r *http.Request) {
	R.RenderJSON(w, "pong!")
}

func (pm *PublicMethods) getExerciseCtrl(w http.ResponseWriter, r *http.Request) {
	slugParam := chi.URLParam(r, "slug")
	log.Debug().Msgf("getExerciseCtrl:%s", slugParam)

	slug := store.ExerciseSlug(slugParam)

	exercise, err := pm.store.Exercise.Get(slug)
	if err != nil {
		log.Error().Msgf("getExerciseCtrl:%s", err)
		R.RenderJSON(w, err)
		return
	}

	log.Info().Msgf("exercise:%+v", exercise)

	R.RenderJSON(w, exercise)
}
