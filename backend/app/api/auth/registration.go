package auth

import (
	"net/http"
	"time"

	R "github.com/go-pkgz/rest"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"

	"github.com/healthy-heroes/neskuchka/backend/app/api/renderer"
	"github.com/healthy-heroes/neskuchka/backend/app/internal/token"
)

func (s *Service) registerUser(w http.ResponseWriter, r *http.Request) {
	logger := log.With().Str("method", "registerUser").Logger()

	user := &UserRegistrationSchema{}
	err := R.DecodeJSON(r, user)
	if err != nil {
		renderer.RenderError(w, logger, http.StatusBadRequest, err, "Failed to decode user data")
		return
	}

	err = user.Validate()
	if err != nil {
		renderer.RenderValidationError(w, logger, err)
		return
	}

	logger.Debug().Msgf("Received user data: %+v", user)

	jti, err := token.RandID()
	if err != nil {
		renderer.RenderError(w, logger, http.StatusInternalServerError, err, "Failed to generate jti")
		return
	}

	claims := RegistrationClaims{
		Data: user,
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
