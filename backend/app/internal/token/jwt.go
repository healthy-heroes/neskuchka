package token

import (
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	// default names for cookies and headers
	defaultJWTCookieName   = "JWT"
	defaultJWTCookieDomain = ""
	defaultJWTHeaderKey    = "X-JWT"
	defaultXSRFCookieName  = "XSRF-TOKEN"
	defaultXSRFHeaderKey   = "X-XSRF-TOKEN"

	defaultIssuer = "go-pkgz/auth"

	defaultTokenDuration  = time.Minute * 15
	defaultCookieDuration = time.Hour * 24 * 31

	defaultTokenQuery = "token"
)

type Service struct {
	Opts
}

type Opts struct {
	Issuer         string
	Secret         string
	TokenDuration  time.Duration
	CookieDuration time.Duration
	SameSite       http.SameSite
}

func NewService(opts Opts) *Service {
	return &Service{Opts: opts}
}

// Token makes token with claims
func (js *Service) Token(claims jwt.Claims) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	if js.Secret == "" {
		return "", fmt.Errorf("secret is empty")
	}

	token, err := t.SignedString([]byte(js.Secret))
	if err != nil {
		return "", fmt.Errorf("can't sign token: %w", err)
	}

	return token, nil
}

func RandID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("can't get random: %w", err)
	}
	s := sha1.New()
	if _, err := s.Write(b); err != nil {
		return "", fmt.Errorf("can't write randoms to sha1: %w", err)
	}
	return fmt.Sprintf("%x", s.Sum(nil)), nil
}
