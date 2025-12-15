package auth

import (
	"net/http"

	"github.com/healthy-heroes/neskuchka/backend/app/api/httpx"
)

func (s *Service) Logout(w http.ResponseWriter, r *http.Request) {
	s.sessionManager.Clear(w)

	httpx.Render(w, nil)
}
