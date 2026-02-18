package auth

import (
	"net/http"

	"github.com/healthy-heroes/neskuchka/backend/app/api/httpx"
	"github.com/healthy-heroes/neskuchka/backend/app/domain"
	"github.com/healthy-heroes/neskuchka/backend/app/internal/session"
)

func (s *Service) User(w http.ResponseWriter, r *http.Request) {
	id := domain.UserID(session.MustGetUserID(r))

	user, err := s.dataStore.GetUser(r.Context(), id)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to get user")
		httpx.RenderUnauthorized(w)
		return
	}

	httpx.Render(w, UserSchema{
		ID:   string(user.ID),
		Name: user.Name,
	})
}
