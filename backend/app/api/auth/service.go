package auth

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/maypok86/otter/v2"
	"github.com/rs/zerolog"

	"github.com/healthy-heroes/neskuchka/backend/app/domain"
	"github.com/healthy-heroes/neskuchka/backend/app/internal/token"
)

// VerifyTokenService defines interface accessing tokens
type VerifyTokenService interface {
	Token(jwt.Claims) (string, error)
	Parse(string, jwt.Claims) error
}

// SessionManager defines interface managing sessions
type SessionManager interface {
	Set(w http.ResponseWriter, userID string) error
	Clear(w http.ResponseWriter)
}

// EmailSender defines interface sending emails
type EmailSender interface {
	Send(to, subject, text string) error
}

// EmailTemplater defines interface templating emails
type EmailTemplater interface {
	AuthLink(token string) (string, error)
}

type Service struct {
	opts   Opts
	logger zerolog.Logger

	dataStore      *domain.Store
	tokenService   VerifyTokenService
	sessionManager SessionManager

	emailSender    EmailSender
	emailTemplater EmailTemplater

	jtiCache *otter.Cache[string, string]
}

type Opts struct {
	Issuer string
	Secret string

	EmailSender    EmailSender
	EmailTemplater EmailTemplater

	Logger zerolog.Logger
}

func NewService(dataStore *domain.Store, session SessionManager, opts Opts) *Service {
	s := &Service{
		opts:   opts,
		logger: opts.Logger.With().Str("pkg", "auth").Logger(),

		dataStore:      dataStore,
		sessionManager: session,

		tokenService: token.NewService(token.Opts{
			Issuer: opts.Issuer,
			Secret: opts.Secret,
		}),

		emailSender:    opts.EmailSender,
		emailTemplater: opts.EmailTemplater,

		jtiCache: otter.Must(&otter.Options[string, string]{
			MaximumSize:     10_000,
			InitialCapacity: 1_000,
		}),
	}

	return s
}
