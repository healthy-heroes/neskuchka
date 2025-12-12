package auth

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/maypok86/otter/v2"
	"github.com/rs/zerolog"

	"github.com/healthy-heroes/neskuchka/backend/app/internal/token"
	"github.com/healthy-heroes/neskuchka/backend/app/store/datastore"
)

// VerifyTokenService defines interface accessing tokens
type VerifyTokenService interface {
	Token(jwt.Claims) (string, error)
	Parse(string, jwt.Claims) error
	Set(w http.ResponseWriter, claims jwt.Claims) error
}

type Service struct {
	opts Opts

	store        *datastore.DataStore
	tokenService VerifyTokenService
	logger       zerolog.Logger

	jtiCache *otter.Cache[string, string]
}

type Opts struct {
	Issuer string
	Secret string
	Logger zerolog.Logger
}

func NewService(store *datastore.DataStore, opts Opts) *Service {
	s := &Service{
		opts:   opts,
		store:  store,
		logger: opts.Logger.With().Str("pkg", "auth").Logger(),
	}

	s.tokenService = token.NewService(token.Opts{
		Issuer:         s.opts.Issuer,
		Secret:         s.opts.Secret,
		TokenDuration:  time.Minute * 15,
		CookieDuration: time.Hour * 24 * 7,
		SameSite:       http.SameSiteLaxMode,
	})

	s.jtiCache = otter.Must(&otter.Options[string, string]{
		MaximumSize:     10_000,
		InitialCapacity: 1_000,
	})

	return s
}

func (s *Service) MountHandlers(router chi.Router) {
	router.Route("/auth", func(r chi.Router) {
		r.Post("/login", s.login)
		r.Post("/login/confirm", s.confirm)
	})
}
