package tracks

import (
	"github.com/healthy-heroes/neskuchka/backend/app/domain"
	"github.com/rs/zerolog"
)

// Service represents tracks endpoints
type Service struct {
	dataStore *domain.Store
	logger    zerolog.Logger
}

// Opts contains options for the service
type Opts struct {
	Logger zerolog.Logger
}

func NewService(dataStore *domain.Store, opts Opts) *Service {
	return &Service{
		dataStore: dataStore,
		logger:    opts.Logger.With().Str("pkg", "tracks").Logger(),
	}
}
