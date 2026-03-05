package api_user

import (
	"context"

	"github.com/rs/zerolog"

	"github.com/healthy-heroes/neskuchka/backend/app/domain"
)

type AvatarStorage interface {
	Get(context.Context, domain.UserID) (domain.Avatar, error)
	Save(context.Context, domain.UserID, domain.Avatar) error
	Exists(context.Context, domain.UserID) (bool, error)
}

type AvatarURLFunc func(domain.UserID) string

// Service is a service for user API
type Service struct {
	logger        zerolog.Logger
	dataStore     *domain.Store
	avatarStorage AvatarStorage
	avatarURLFunc AvatarURLFunc
}

type Opts struct {
	Logger zerolog.Logger

	AvatarStorage AvatarStorage
	AvatarURLFunc AvatarURLFunc
}

func NewService(dataStore *domain.Store, opts Opts) *Service {
	return &Service{
		logger:        opts.Logger,
		dataStore:     dataStore,
		avatarStorage: opts.AvatarStorage,
		avatarURLFunc: opts.AvatarURLFunc,
	}
}
