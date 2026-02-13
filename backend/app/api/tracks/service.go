package tracks

import (
	"net/http"

	"github.com/healthy-heroes/neskuchka/backend/app/domain"
	"github.com/healthy-heroes/neskuchka/backend/app/internal/session"
	"github.com/rs/zerolog"
)

// SessionManager is the interface for getting current user ID and claims from session
type SessionManager interface {
	Get(r *http.Request) (string, *session.Claims, error)
}

// Service represents tracks endpoints
type Service struct {
	dataStore *domain.Store
	session   SessionManager
	logger    zerolog.Logger
}

// Opts contains options for the service
type Opts struct {
	Logger zerolog.Logger
}

func NewService(dataStore *domain.Store, session SessionManager, opts Opts) *Service {
	return &Service{
		dataStore: dataStore,
		session:   session,
		logger:    opts.Logger.With().Str("pkg", "tracks").Logger(),
	}
}
