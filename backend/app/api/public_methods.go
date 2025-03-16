package api

import (
	"net/http"

	R "github.com/go-pkgz/rest"
)

type PublicMethods struct {
}

func (pm *PublicMethods) pingCtrl(w http.ResponseWriter, r *http.Request) {
	R.RenderJSON(w, "pong!")
}
