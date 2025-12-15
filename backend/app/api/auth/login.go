package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/healthy-heroes/neskuchka/backend/app/api/httpx"
	"github.com/healthy-heroes/neskuchka/backend/app/internal/token"
)

const (
	confTokenTtlDuration = 30 * time.Minute
)

func (s *Service) login(w http.ResponseWriter, r *http.Request) {
	logger := s.logger

	loginData, ok := httpx.ParseAndValidateBody[LoginSchema](w, r, logger)
	if !ok {
		return
	}
	logger.Debug().Msgf("Received user data: %+v", loginData)

	jti, err := token.RandID()
	if err != nil {
		httpx.RenderError(w, logger, http.StatusInternalServerError, err, "Failed to generate JTI")
		return
	}

	claims := ConfirmationClaims{
		Data: loginData,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(confTokenTtlDuration)),
			Issuer:    s.opts.Issuer,
			ID:        jti,
		},
	}
	logger.Debug().Msgf("Make claims: %+v", claims)

	token, err := s.tokenService.Token(claims)
	if err != nil {
		httpx.RenderError(w, logger, http.StatusInternalServerError, err, "Failed to generate verification token")
		return
	}
	logger.Debug().Msgf("Token: %s", token)
}

func (s *Service) confirm(w http.ResponseWriter, r *http.Request) {
	logger := s.logger

	data, ok := httpx.ParseBody[ConfirmationSchema](w, r, logger)
	if !ok {
		return
	}

	var confClaims ConfirmationClaims
	err := s.tokenService.Parse(data.Token, &confClaims)
	if err != nil {
		httpx.RenderError(w, logger, http.StatusBadRequest, err, "Failed to parse token")
		return
	}

	// check if token is already used
	if _, ok := s.jtiCache.GetIfPresent(confClaims.ID); ok {
		httpx.RenderError(w, logger, http.StatusBadRequest, fmt.Errorf("token already used"), "Token already used")
		return
	}

	user, err := s.store.User.FindOrCreate(confClaims.Data.Email)
	if err != nil {
		httpx.RenderError(w, logger, http.StatusInternalServerError, err, "Failed to find or create user")
		return
	}

	// set token as used
	s.jtiCache.Set(confClaims.ID, "")
	s.jtiCache.SetExpiresAfter(confClaims.ID, confTokenTtlDuration+time.Minute*5) // add extra time

	err = s.sessionManager.Set(w, string(user.ID))
	if err != nil {
		httpx.RenderError(w, logger, http.StatusInternalServerError, err, "Failed to set session")
		return
	}

	httpx.Render(w, user)
}
