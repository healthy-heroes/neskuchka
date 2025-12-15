package auth

import (
	"net/http"

	"github.com/healthy-heroes/neskuchka/backend/app/api/httpx"
)

func (s *Service) user(w http.ResponseWriter, r *http.Request) {
	user, err := s.getUser(r)
	if err != nil {
		httpx.RenderError(w, s.logger, http.StatusUnauthorized, err, "Unauthorized")
		return
	}

	httpx.Render(w, user)
}
