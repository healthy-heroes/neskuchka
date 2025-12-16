package auth

import (
	"net/http"

	"github.com/healthy-heroes/neskuchka/backend/app/api/httpx"
	"github.com/healthy-heroes/neskuchka/backend/app/internal/session"
	"github.com/healthy-heroes/neskuchka/backend/app/store"
)

func (s *Service) User(w http.ResponseWriter, r *http.Request) {
	userID := session.MustGetUserID(r)

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
