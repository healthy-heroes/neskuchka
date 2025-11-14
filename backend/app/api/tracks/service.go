package tracks

import (
	"github.com/go-chi/chi/v5"
	"github.com/healthy-heroes/neskuchka/backend/app/store/datastore"
)

type Service struct {
	store *datastore.DataStore
}

func NewService(store *datastore.DataStore) *Service {
	return &Service{store}
}

func (s *Service) MountHandlers(router chi.Router) {
	// Concrete main track routes
	router.Route("/tracks/main", func(r chi.Router) {
		r.Get("/last_workouts", s.getMainTrackLastWorkoutsCtrl)

		r.Get("/workouts/{id}", s.getWorkoutCtrl)
		r.Post("/workouts", s.createWorkoutCtrl)
		r.Put("/workouts/{id}", s.updateWorkoutCtrl)
	})

}
