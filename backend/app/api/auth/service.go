package auth

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/healthy-heroes/neskuchka/backend/app/internal/token"
	"github.com/healthy-heroes/neskuchka/backend/app/store/datastore"
)

// VerifyTokenService defines interface accessing tokens
type VerifyTokenService interface {
	Token(token.Claims) (string, error)
}

type Service struct {
	opts         Opts
	store        *datastore.DataStore
	tokenService VerifyTokenService
}

type Opts struct {
	Issuer string
	Secret string
}

func NewService(store *datastore.DataStore, opts Opts) *Service {
	s := &Service{
		opts:  opts,
		store: store,
	}

	s.tokenService = token.NewService(token.Opts{
		Issuer:         s.opts.Issuer,
		Secret:         s.opts.Secret,
		TokenDuration:  time.Minute * 15,
		CookieDuration: time.Hour * 24 * 7,
		SameSite:       http.SameSiteLaxMode,
	})

	return s
}

func (s *Service) MountHandlers(router chi.Router) {
	router.Route("/auth", func(r chi.Router) {
		r.Post("/register", s.registerUser)
	})
}
