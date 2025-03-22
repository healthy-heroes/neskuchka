package public_api

import (
	"github.com/go-chi/chi/v5"
	"github.com/healthy-heroes/neskuchka/backend/app/store/datastore"
)

type PublicAPI struct {
	store *datastore.DataStore
}

func NewPublicAPI(store *datastore.DataStore) *PublicAPI {
	return &PublicAPI{
		store: store,
	}
}

func (api *PublicAPI) InitRoutes(router chi.Router) {
	// middlewares

	router.Route("/exercises", func(router chi.Router) {
		router.Get("/{slug}", api.getExerciseCtrl)
	})

	router.Route("/tracks", func(router chi.Router) {
		router.Get("/main/last_workouts", api.getMainTrackLastExercisesCtrl)
	})
}
