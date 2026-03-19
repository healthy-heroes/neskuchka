package api_user

import (
	"context"

	"github.com/rs/zerolog"

	"github.com/healthy-heroes/neskuchka/backend/app/domain"
)

type AvatarStore interface {
	Get(context.Context, domain.UserID) (domain.Avatar, error)
	Save(context.Context, domain.UserID, domain.Avatar) error
	Delete(context.Context, domain.UserID) error
	Exists(context.Context, domain.UserID) (bool, error)
}

type AvatarURLFunc func(domain.UserID) string

// Service is a service for user API
type Service struct {
	logger        zerolog.Logger
	dataStore     *domain.Store
	avatarStore   AvatarStore
	avatarURLFunc AvatarURLFunc
}

type Opts struct {
	Logger zerolog.Logger

	AvatarStore   AvatarStore
	AvatarURLFunc AvatarURLFunc
}

func NewService(dataStore *domain.Store, opts Opts) *Service {
	return &Service{
		logger:        opts.Logger,
		dataStore:     dataStore,
		avatarStore:   opts.AvatarStore,
		avatarURLFunc: opts.AvatarURLFunc,
	}
}
