package auth

import (
	"errors"
	"net/http"
	"time"

	R "github.com/go-pkgz/rest"
	"github.com/golang-jwt/jwt/v5"
	"github.com/healthy-heroes/neskuchka/backend/app/store"
	"github.com/rs/zerolog/log"

	"github.com/healthy-heroes/neskuchka/backend/app/api/renderer"
	"github.com/healthy-heroes/neskuchka/backend/app/internal/token"
)

func (s *Service) registerUser(w http.ResponseWriter, r *http.Request) {
	logger := log.With().Str("method", "registerUser").Logger()

	// decode user and validation
	newUser := &UserRegistrationSchema{}
	err := R.DecodeJSON(r, newUser)
	if err != nil {
		renderer.RenderError(w, logger, http.StatusBadRequest, err, "Failed to decode user data")
		return
	}

	err = newUser.Validate()
	if err != nil {
		renderer.RenderValidationError(w, logger, err)
		return
	}

	logger.Debug().Msgf("Received user data: %+v", newUser)

	oldUser, err := s.store.User.FindByEmail(newUser.Email)
	if err != nil && !errors.Is(err, store.ErrNotFound) {
		renderer.RenderError(w, logger, http.StatusInternalServerError, err, "Failed to find user")
		return
	}

	if oldUser != nil {
		renderer.RenderError(w, logger, http.StatusConflict, err, "User already exists")
		return
	}

	jti, err := token.RandID()
	if err != nil {
		renderer.RenderError(w, logger, http.StatusInternalServerError, err, "Failed to generate JWT")
		return
	}

	claims := RegistrationClaims{
		Data: newUser,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Minute)),
			NotBefore: jwt.NewNumericDate(time.Now().Add(-1 * time.Minute)),
			Issuer:    s.opts.Issuer,
			ID:        jti,
		},
	}
	logger.Debug().Msgf("Make claims: %+v", claims)

	verifyToken, err := s.tokenService.Token(claims)
	if err != nil {
		renderer.RenderError(w, logger, http.StatusInternalServerError, err, "Failed to generate verification token")
		return
	}

	logger.Debug().Msgf("Token: %s", verifyToken)

	//todo: delete it
	renderer.Render(w, &TempResponse{
		Token:  verifyToken,
		Claims: &claims,
	})
}

type TempResponse struct {
	Claims jwt.Claims `json:"claims"`
	Token  string     `json:"token"`
}
