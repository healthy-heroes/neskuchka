package tracks

import (
	"net/http"

	"github.com/healthy-heroes/neskuchka/backend/app/internal/session"
	"github.com/healthy-heroes/neskuchka/backend/app/store/datastore"
	"github.com/rs/zerolog"
)

// SessionManager is the interface for getting current user ID and claims from session
type SessionManager interface {
	Get(r *http.Request) (string, *session.Claims, error)
}

// Service represents tracks endpoints
type Service struct {
	store   *datastore.DataStore
	session SessionManager
	logger  zerolog.Logger
}

// Opts contains options for the service
type Opts struct {
	Logger zerolog.Logger
}

func NewService(store *datastore.DataStore, session SessionManager, opts Opts) *Service {
	return &Service{
		store:   store,
		session: session,
		logger:  opts.Logger.With().Str("pkg", "tracks").Logger(),
	}
}
