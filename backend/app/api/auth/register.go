package auth

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	R "github.com/go-pkgz/rest"
	"github.com/golang-jwt/jwt/v5"
	"github.com/healthy-heroes/neskuchka/backend/app/store"

	"github.com/healthy-heroes/neskuchka/backend/app/api/renderer"
	"github.com/healthy-heroes/neskuchka/backend/app/internal/token"
)

func (s *Service) register(w http.ResponseWriter, r *http.Request) {
	logger := s.logger.With().Str("method", "register").Logger()

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
		renderer.RenderError(w, logger, http.StatusInternalServerError, err, "Failed to generate JTI")
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
		Claims: claims,
	})
}

type TempResponse struct {
	Claims jwt.Claims `json:"claims"`
	Token  string     `json:"token"`
}

func (s *Service) confirmRegistration(w http.ResponseWriter, r *http.Request) {
	logger := s.logger.With().Str("method", "confirmRegistration").Logger()

	tokenString := r.URL.Query().Get("token")
	if tokenString == "" {
		renderer.RenderError(w, logger, http.StatusBadRequest, fmt.Errorf("missing token"), "Missing token")
		return
	}

	var regClaims RegistrationClaims
	err := s.tokenService.Parse(tokenString, &regClaims)
	if err != nil {
		renderer.RenderError(w, logger, http.StatusBadRequest, err, "Failed to parse token")
		return
	}

	// todo:check used jti

	user, err := s.store.User.Create(&store.User{
		ID:    store.CreateUserId(),
		Name:  regClaims.Data.Name,
		Email: regClaims.Data.Email,
	})
	if err != nil {
		renderer.RenderError(w, logger, http.StatusInternalServerError, err, "Failed to create user")
		return
	}

	// todo:save jti

	jti, err := token.RandID()
	if err != nil {
		renderer.RenderError(w, logger, http.StatusInternalServerError, err, "Failed to generate JTI")
		return
	}

	claims := AccessClaims{
		Data: &UserSchema{
			ID:   user.ID,
			Name: user.Name,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ID:     jti,
			Issuer: s.opts.Issuer,
		},
	}

	err = s.tokenService.Set(w, claims)
	if err != nil {
		renderer.RenderError(w, logger, http.StatusInternalServerError, err, "Failed to set token")
	}

	//todo: delete it
	renderer.Render(w, &TempResponse{
		Token:  tokenString,
		Claims: regClaims,
	})
}
