package token

import (
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

func (js *Service) Parse(tokenString string, claims jwt.Claims) error {
	parser := jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	token, err := parser.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(js.Secret), nil
	})

	if err != nil {
		return fmt.Errorf("can't parse token: %w", err)
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

func (js *Service) Set(w http.ResponseWriter, claims jwt.Claims) error {
	return nil
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
