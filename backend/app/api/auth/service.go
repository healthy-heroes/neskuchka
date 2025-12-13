package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/maypok86/otter/v2"
	"github.com/rs/zerolog"

	"github.com/healthy-heroes/neskuchka/backend/app/internal/token"
	"github.com/healthy-heroes/neskuchka/backend/app/store/datastore"
)

const (
	JWTCookieName = "JWT"

	defaultSessionDuration = 7 * 24 * time.Hour
)

// VerifyTokenService defines interface accessing tokens
type VerifyTokenService interface {
	Token(jwt.Claims) (string, error)
	Parse(string, jwt.Claims) error
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

	SecureCookies   bool
	SessionDuration time.Duration
}

func NewService(store *datastore.DataStore, opts Opts) *Service {
	if opts.SessionDuration == 0 {
		opts.SessionDuration = defaultSessionDuration
	}

	s := &Service{
		opts:   opts,
		store:  store,
		logger: opts.Logger.With().Str("pkg", "auth").Logger(),
	}

	s.tokenService = token.NewService(token.Opts{
		Issuer: s.opts.Issuer,
		Secret: s.opts.Secret,
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

func (s *Service) setToken(w http.ResponseWriter, user UserSchema) error {
	now := time.Now().Unix()

	jti, err := token.RandID()
	if err != nil {
		return fmt.Errorf("failed to generate JTI: %w", err)
	}

	claims := UserClaims{
		Data: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			Issuer:    s.opts.Issuer,
			IssuedAt:  jwt.NewNumericDate(time.Unix(now, 0)),
			ExpiresAt: jwt.NewNumericDate(time.Unix(now, 0).Add(s.opts.SessionDuration)),
		},
	}

	tokenString, err := s.tokenService.Token(claims)
	if err != nil {
		return fmt.Errorf("failed to generate token: %w", err)
	}

	jwtCookie := http.Cookie{
		Name:     JWTCookieName,
		Value:    tokenString,
		HttpOnly: true,
		Path:     "/",
		Domain:   "",
		MaxAge:   int(s.opts.SessionDuration.Seconds()),
		Secure:   s.opts.SecureCookies,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &jwtCookie)

	return nil
}
