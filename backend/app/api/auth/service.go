package auth

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
	router.Route("/auth", func(r chi.Router) {
		r.Post("/register", s.registerUser)
	})
}
