package auth

import (
	"net/http"

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
}

type SessionManager interface {
	Set(w http.ResponseWriter, userID string) error
	Clear(w http.ResponseWriter)
}

type Service struct {
	opts   Opts
	logger zerolog.Logger

	store          *datastore.DataStore
	tokenService   VerifyTokenService
	sessionManager SessionManager

	jtiCache *otter.Cache[string, string]
}

type Opts struct {
	Issuer string
	Secret string
	Logger zerolog.Logger
}

func NewService(store *datastore.DataStore, session SessionManager, opts Opts) *Service {
	s := &Service{
		opts:   opts,
		logger: opts.Logger.With().Str("pkg", "auth").Logger(),

		store:          store,
		sessionManager: session,

		tokenService: token.NewService(token.Opts{
			Issuer: opts.Issuer,
			Secret: opts.Secret,
		}),

		jtiCache: otter.Must(&otter.Options[string, string]{
			MaximumSize:     10_000,
			InitialCapacity: 1_000,
		}),
	}

	return s
}
