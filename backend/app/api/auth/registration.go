package auth

import (
	"net/http"

	R "github.com/go-pkgz/rest"
	"github.com/healthy-heroes/neskuchka/backend/app/api/renderer"
	"github.com/rs/zerolog/log"
)

func (s *Service) registerUser(w http.ResponseWriter, r *http.Request) {
	logger := log.With().Str("method", "registerUser").Logger()

	data := &UserRegistrationSchema{}
	err := R.DecodeJSON(r, data)
	if err != nil {
		renderer.RenderError(w, logger, http.StatusBadRequest, err, "Failed to decode user data")
		return
	}

	err = data.Validate()
	if err != nil {
		renderer.RenderValidationError(w, logger, err)
		return
	}

	logger.Debug().Msgf("Received user data: %+v", data)
}
