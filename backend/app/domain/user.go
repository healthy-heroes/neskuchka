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

// GetUser gets a user by id
func (s *Store) GetUser(ctx context.Context, id UserID) (User, error) {
	return s.dataStorage.GetUser(ctx, id)
}

// FindOrCreateUser finds a user by email or creates a new user
func (s *Store) FindOrCreateUser(ctx context.Context, u User) (User, error) {
	if u.Email == "" {
		return User{}, errors.New("email is required")
	}

	user, err := s.dataStorage.GetUserByEmail(ctx, u.Email)
	if err != nil {
		if err != ErrNotFound {
			return User{}, err
		}
	} else {
		return user, nil
	}

	u.ID = NewUserID()
	return s.dataStorage.CreateUser(ctx, u)
}

// UpdateUser updates a user
// updates only safe fields, other should be ignored
func (s *Store) UpdateUser(ctx context.Context, u User) (User, error) {
	user, err := s.dataStorage.GetUser(ctx, u.ID)
	if err != nil {
		return User{}, err
	}

	user.Name = u.Name

	return s.dataStorage.UpdateUser(ctx, user)
}
