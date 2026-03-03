package api_user

import (
	"github.com/rs/zerolog"

	"github.com/healthy-heroes/neskuchka/backend/app/domain"
)

// Service is a service for user API
type Service struct {
	logger    zerolog.Logger
	dataStore *domain.Store
}

type Opts struct {
	Logger zerolog.Logger
}

func NewService(dataStore *domain.Store, opts Opts) *Service {
	return &Service{
		logger:    opts.Logger,
		dataStore: dataStore,
	}
}
