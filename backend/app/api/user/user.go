package api_user

import (
	"net/http"

	"github.com/healthy-heroes/neskuchka/backend/app/api/httpx"
	"github.com/healthy-heroes/neskuchka/backend/app/domain"
	"github.com/healthy-heroes/neskuchka/backend/app/internal/session"
)

func (s *Service) Me(w http.ResponseWriter, r *http.Request) {
	id := domain.UserID(session.MustGetUserID(r))

	user, err := s.dataStore.GetUser(r.Context(), id)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to get user")
		httpx.RenderUnauthorized(w)
		return
	}

	response := UserSchema{
		ID:   string(user.ID),
		Name: user.Name,
	}

	exists, err := s.avatarStore.Exists(r.Context(), user.ID)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to check if avatar exists")
		exists = false
	}

	if exists {
		response.Avatar = s.avatarURLFunc(user.ID)
	}

	httpx.Render(w, response)
}

func (s *Service) GetSettings(w http.ResponseWriter, r *http.Request) {
	id := domain.UserID(session.MustGetUserID(r))

	user, err := s.dataStore.GetUser(r.Context(), id)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to get user")
		httpx.RenderUnauthorized(w)
		return
	}

	httpx.Render(w, MakeSettingsSchema(user))
}

func (s *Service) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	id := domain.UserID(session.MustGetUserID(r))

	body, ok := httpx.ParseAndValidateBody[UpdateSettingsSchema](w, r, s.logger)
	if !ok {
		return
	}

	user, err := s.dataStore.UpdateUser(r.Context(), domain.User{
		ID:   id,
		Name: body.Name,
	})
	if err != nil {
		httpx.RenderDomainError(w, s.logger, err, "failed to update user")
		return
	}

	httpx.Render(w, MakeSettingsSchema(user))
}
