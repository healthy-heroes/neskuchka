package auth

import (
	"fmt"
	"net/http"
	"time"

	R "github.com/go-pkgz/rest"
	"github.com/golang-jwt/jwt/v5"

	"github.com/healthy-heroes/neskuchka/backend/app/api/renderer"
	"github.com/healthy-heroes/neskuchka/backend/app/internal/token"
)

const (
	confTokenTtlDuration = 30 * time.Minute
	userTokenTtlDuration = 7 * 24 * time.Hour
)

func (s *Service) login(w http.ResponseWriter, r *http.Request) {
	logger := s.logger

	loginData := &LoginSchema{}
	err := R.DecodeJSON(r, loginData)
	if err != nil {
		renderer.RenderError(w, logger, http.StatusBadRequest, err, "Failed to decode user data")
		return
	}

	err = loginData.Validate()
	if err != nil {
		renderer.RenderValidationError(w, logger, err)
		return
	}

	logger.Debug().Msgf("Received user data: %+v", loginData)

	jti, err := token.RandID()
	if err != nil {
		renderer.RenderError(w, logger, http.StatusInternalServerError, err, "Failed to generate JTI")
		return
	}

	claims := ConfirmationClaims{
		Data: loginData,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(confTokenTtlDuration)),
			NotBefore: jwt.NewNumericDate(time.Now().Add(-1 * time.Minute)),
			Issuer:    s.opts.Issuer,
			ID:        jti,
		},
	}

	logger.Debug().Msgf("Make claims: %+v", claims)

	token, err := s.tokenService.Token(claims)
	if err != nil {
		renderer.RenderError(w, logger, http.StatusInternalServerError, err, "Failed to generate verification token")
		return
	}

	logger.Debug().Msgf("Token: %s", token)

	//todo: delete it
	renderer.Render(w, &TempResponse{
		Token:  token,
		Claims: claims,
	})
}

type TempResponse struct {
	Claims jwt.Claims `json:"claims"`
	Token  string     `json:"token"`
}

func (s *Service) confirm(w http.ResponseWriter, r *http.Request) {
	logger := s.logger

	var data ConfirmationSchema
	err := R.DecodeJSON(r, &data)
	if err != nil {
		renderer.RenderError(w, logger, http.StatusBadRequest, err, "Failed to decode confirm data")
		return
	}

	var confClaims ConfirmationClaims
	err = s.tokenService.Parse(data.Token, &confClaims)
	if err != nil {
		renderer.RenderError(w, logger, http.StatusBadRequest, err, "Failed to parse token")
		return
	}

	if _, ok := s.jtiCache.GetIfPresent(confClaims.ID); ok {
		renderer.RenderError(w, logger, http.StatusBadRequest, fmt.Errorf("token already used"), "Token already used")
		return
	}

	// todo: use instead store
	user, err := s.store.User.FindOrCreate(confClaims.Data.Email)
	if err != nil {
		renderer.RenderError(w, logger, http.StatusInternalServerError, err, "Failed to find or create user")
		return
	}

	s.jtiCache.Set(confClaims.ID, "")
	s.jtiCache.SetExpiresAfter(confClaims.ID, confTokenTtlDuration+time.Minute*5) // add extra time

	jti, err := token.RandID()
	if err != nil {
		renderer.RenderError(w, logger, http.StatusInternalServerError, err, "Failed to generate JTI")
		return
	}

	claims := UserClaims{
		Data: &UserSchema{
			ID:   user.ID,
			Name: user.Name,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			Issuer:    s.opts.Issuer,
			NotBefore: jwt.NewNumericDate(time.Now().Add(-1 * time.Minute)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(userTokenTtlDuration)),
		},
	}

	err = s.tokenService.Set(w, claims)
	if err != nil {
		renderer.RenderError(w, logger, http.StatusInternalServerError, err, "Failed to set token")
	}

	//todo: delete it
	renderer.Render(w, &TempResponse{
		Claims: claims,
	})
}
