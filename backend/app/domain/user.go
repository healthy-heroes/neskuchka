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

// UserRepo is a interface for user storage
type UserRepo interface {
	Get(context.Context, UserID) (User, error)
	GetByEmail(context.Context, Email) (User, error)
	Create(context.Context, User) (User, error)
	Update(context.Context, User) (User, error)
}

// GetUser gets a user by id
func (s *Store) GetUser(ctx context.Context, id UserID) (User, error) {
	return s.userRepo.Get(ctx, id)
}

// FindOrCreateUser finds a user by email or creates a new user
func (s *Store) FindOrCreateUser(ctx context.Context, u User) (User, error) {
	if u.Email == "" {
		return User{}, errors.New("email is required")
	}

	user, err := s.userRepo.GetByEmail(ctx, u.Email)
	if err != nil {
		if !errors.Is(err, ErrNotFound) {
			return User{}, err
		}
	} else {
		return user, nil
	}

	u.ID = NewUserID()
	return s.userRepo.Create(ctx, u)
}

// UpdateUser updates a user
// updates only safe fields, other should be ignored
func (s *Store) UpdateUser(ctx context.Context, u User) (User, error) {
	user, err := s.userRepo.Get(ctx, u.ID)
	if err != nil {
		return User{}, err
	}

	user.Name = u.Name

	return s.userRepo.Update(ctx, user)
}
