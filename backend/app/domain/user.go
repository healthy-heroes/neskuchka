package domain

import (
	"context"
	"errors"

	"github.com/healthy-heroes/neskuchka/backend/app/internal/uuid"
)

type UserID string
type Email string

// NewUserID generates a new user id
func NewUserID() UserID {
	return UserID(uuid.New())
}

// User is a user aggregate
type User struct {
	ID    UserID
	Name  string
	Email Email
}

// UserStore is a interface for user storage
type UserStore interface {
	Get(context.Context, UserID) (User, error)
	GetByEmail(context.Context, Email) (User, error)
	Create(context.Context, User) (User, error)
	Update(context.Context, User) (User, error)
}

// GetUser gets a user by id
func (s *Service) GetUser(ctx context.Context, id UserID) (User, error) {
	return s.userStore.Get(ctx, id)
}

// FindOrCreateUser finds a user by email or creates a new user
func (s *Service) FindOrCreateUser(ctx context.Context, u User) (User, error) {
	if u.Email == "" {
		return User{}, errors.New("email is required")
	}

	user, err := s.userStore.GetByEmail(ctx, u.Email)
	if err != nil {
		if !errors.Is(err, ErrNotFound) {
			return User{}, err
		}
	} else {
		return user, nil
	}

	u.ID = NewUserID()
	return s.userStore.Create(ctx, u)
}

// UpdateUser updates a user
// updates only safe fields, other should be ignored
func (s *Service) UpdateUser(ctx context.Context, u User) (User, error) {
	user, err := s.userStore.Get(ctx, u.ID)
	if err != nil {
		return User{}, err
	}

	user.Name = u.Name

	return s.userStore.Update(ctx, user)
}
