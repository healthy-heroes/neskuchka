package api_user

import (
	"context"

	"github.com/rs/zerolog"

	"github.com/healthy-heroes/neskuchka/backend/app/domain"
)

type AvatarStorage interface {
	Get(context.Context, domain.UserID) (domain.Avatar, error)
	Save(context.Context, domain.UserID, domain.Avatar) error
}

// Service is a service for user API
type Service struct {
	logger        zerolog.Logger
	dataStore     *domain.Store
	avatarStorage AvatarStorage
}

type Opts struct {
	Logger zerolog.Logger

	AvatarStorage AvatarStorage
}

func NewService(dataStore *domain.Store, opts Opts) *Service {
	return &Service{
		logger:        opts.Logger,
		dataStore:     dataStore,
		avatarStorage: opts.AvatarStorage,
	}
}
