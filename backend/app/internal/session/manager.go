package session

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/healthy-heroes/neskuchka/backend/app/internal/token"
	"github.com/rs/zerolog"
)

const (
	defaultSessionCookieName = "JWT"
	defaultSessionDuration   = 7 * 24 * time.Hour
)

var ErrSessionNotFound = errors.New("session not found")

type TokenService interface {
	Token(claims jwt.Claims) (string, error)
	Parse(token string, claims jwt.Claims) error
}

type Claims struct {
	UserID string `json:"uid"`

	jwt.RegisteredClaims
}

// Manager is the session manager
type Manager struct {
	opts   Opts
	logger zerolog.Logger

	tokenService TokenService
}

// Opts are the options for the session manager
type Opts struct {
	Logger zerolog.Logger

	SecureCookies     bool
	SessionCookieName string
	SessionDuration   time.Duration

	Issuer string
	Secret string
}

// NewManager creates a new session manager
func NewManager(opts Opts) *Manager {
	if opts.SessionCookieName == "" {
		opts.SessionCookieName = defaultSessionCookieName
	}

	if opts.SessionDuration == 0 {
		opts.SessionDuration = defaultSessionDuration
	}

	m := &Manager{
		opts:   opts,
		logger: opts.Logger.With().Str("pkg", "session").Logger(),

		tokenService: token.NewService(token.Opts{
			Issuer: opts.Issuer,
			Secret: opts.Secret,
		}),
	}

	return m
}

// Set token to the cookie
func (m *Manager) Set(w http.ResponseWriter, userID string) error {
	jti, err := token.RandID()
	if err != nil {
		return fmt.Errorf("failed to generate JTI: %w", err)
	}

	now := time.Now()
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			Issuer:    m.opts.Issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.opts.SessionDuration)),
		},
	}

	tokenString, err := m.tokenService.Token(claims)
	if err != nil {
		return fmt.Errorf("failed to generate token: %w", err)
	}

	m.setCookie(w, tokenString, int(m.opts.SessionDuration.Seconds()))

	return nil
}

func (m *Manager) Clear(w http.ResponseWriter) {
	m.setCookie(w, "", -1)
}

func (m *Manager) setCookie(w http.ResponseWriter, value string, maxAge int) {
	http.SetCookie(w, &http.Cookie{
		Name:     m.opts.SessionCookieName,
		Value:    value,
		HttpOnly: true,
		Path:     "/",
		Domain:   "",
		MaxAge:   maxAge,
		Secure:   m.opts.SecureCookies,
		SameSite: http.SameSiteLaxMode,
	})
}

func (m *Manager) Get(r *http.Request) (string, *Claims, error) {
	tokenCookie, err := r.Cookie(m.opts.SessionCookieName)
	if err != nil {
		return "", nil, ErrSessionNotFound
	}

	claims := Claims{}
	err = m.tokenService.Parse(tokenCookie.Value, &claims)
	if err != nil {
		return "", nil, fmt.Errorf("failed to parse token: %w", err)
	}

	return claims.UserID, &claims, nil
}
