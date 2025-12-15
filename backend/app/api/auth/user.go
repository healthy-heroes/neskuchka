package auth

import (
	"net/http"

	"github.com/healthy-heroes/neskuchka/backend/app/api/httpx"
	"github.com/healthy-heroes/neskuchka/backend/app/store"
)

func (s *Service) user(w http.ResponseWriter, r *http.Request) {
	userID, _, err := s.sessionManager.Get(r)
	if err != nil {
		httpx.RenderError(w, s.logger, http.StatusUnauthorized, err, "Unauthorized")
		return
	}

	user, err := s.store.User.Get(store.UserID(userID))
	if err != nil {
		httpx.RenderError(w, s.logger, http.StatusInternalServerError, err, "Failed to get user")
		return
	}

	httpx.Render(w, UserSchema{
		ID:   user.ID,
		Name: user.Name,
	})
}
