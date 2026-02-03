package domain

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type UserStoreStub struct {
	GetFunc        func(context.Context, UserID) (User, error)
	GetByEmailFunc func(context.Context, Email) (User, error)
	CreateFunc     func(context.Context, User) (User, error)
	UpdateFunc     func(context.Context, User) (User, error)
}

func (s UserStoreStub) Get(ctx context.Context, id UserID) (User, error) {
	return s.GetFunc(ctx, id)
}

func (s UserStoreStub) GetByEmail(ctx context.Context, email Email) (User, error) {
	return s.GetByEmailFunc(ctx, email)
}

func (s UserStoreStub) Create(ctx context.Context, user User) (User, error) {
	return s.CreateFunc(ctx, user)
}

func (s UserStoreStub) Update(ctx context.Context, user User) (User, error) {
	return s.UpdateFunc(ctx, user)
}

func TestNewUserID(t *testing.T) {
	t.Run("should generate a new user id", func(t *testing.T) {
		userID := NewUserID()
		assert.NotEmpty(t, userID)
	})
}

func TestGetUser(t *testing.T) {
	t.Run("should return user", func(t *testing.T) {
		existingUser := User{
			ID:   UserID("1"),
			Name: "Test User",
		}
		service := NewService(Opts{
			UserStore: UserStoreStub{
				GetFunc: func(ctx context.Context, id UserID) (User, error) {
					return existingUser, nil
				},
			},
		})
		user, err := service.GetUser(context.Background(), UserID("1"))

		assert.Nil(t, err)
		assert.Equal(t, existingUser, user)
	})
}

func TestFindOrCreateUser(t *testing.T) {
	t.Run("should return existing user", func(t *testing.T) {
		existingUser := User{
			ID:    UserID("1"),
			Email: Email("test@example.com"),
			Name:  "Test User",
		}
		service := NewService(Opts{
			UserStore: UserStoreStub{
				GetByEmailFunc: func(ctx context.Context, email Email) (User, error) {
					return existingUser, nil
				},
			},
		})
		user, err := service.FindOrCreateUser(context.Background(), User{
			Email: Email("test@example.com"),
			Name:  "Test User New",
		})

		assert.Nil(t, err)
		assert.Equal(t, existingUser, user)
	})

	t.Run("should create new user", func(t *testing.T) {
		service := NewService(Opts{
			UserStore: UserStoreStub{
				GetByEmailFunc: func(ctx context.Context, email Email) (User, error) {
					return User{}, ErrNotFound
				},
				CreateFunc: func(ctx context.Context, user User) (User, error) {
					return user, nil
				},
			},
		})

		newUser := User{
			ID:    UserID("1"),
			Email: Email("test@example.com"),
			Name:  "Test User",
		}
		user, err := service.FindOrCreateUser(context.Background(), newUser)

		assert.Nil(t, err)
		assert.NotEmpty(t, user.ID)
		assert.NotEqual(t, newUser.ID, user.ID)
		assert.Equal(t, newUser.Email, user.Email)
		assert.Equal(t, newUser.Name, user.Name)
	})

	t.Run("should check required email", func(t *testing.T) {
		service := NewService(Opts{})
		_, err := service.FindOrCreateUser(context.Background(), User{
			Name: "Test User",
		})

		assert.Error(t, err)
		assert.Equal(t, "email is required", err.Error())
	})

	t.Run("should return error", func(t *testing.T) {
		expectedErr := errors.New("some error")
		service := NewService(Opts{
			UserStore: UserStoreStub{
				GetByEmailFunc: func(ctx context.Context, email Email) (User, error) {
					return User{}, expectedErr
				},
			},
		})
		_, err := service.FindOrCreateUser(context.Background(), User{
			Email: Email("test@example.com"),
		})

		assert.Error(t, err)
		assert.ErrorIs(t, err, expectedErr)
	})
}

func TestUpdateUser(t *testing.T) {
	t.Run("should update user", func(t *testing.T) {
		existingUser := User{
			ID:    UserID("1"),
			Email: Email("test@example.com"),
			Name:  "Test User",
		}
		service := NewService(Opts{
			UserStore: UserStoreStub{
				GetFunc: func(ctx context.Context, id UserID) (User, error) {
					return existingUser, nil
				},
				UpdateFunc: func(ctx context.Context, user User) (User, error) {
					return user, nil
				},
			},
		})
		user, err := service.UpdateUser(context.Background(), User{
			ID:    UserID("2"),
			Email: Email("wrong@example.com"),
			Name:  "Test User New",
		})

		assert.Nil(t, err)
		assert.Equal(t, existingUser.ID, user.ID)
		assert.Equal(t, existingUser.Email, user.Email)
		assert.Equal(t, "Test User New", user.Name)
	})

	t.Run("should return error if user not found", func(t *testing.T) {
		service := NewService(Opts{
			UserStore: UserStoreStub{
				GetFunc: func(ctx context.Context, id UserID) (User, error) {
					return User{}, ErrNotFound
				},
			},
		})
		_, err := service.UpdateUser(context.Background(), User{
			ID:   UserID("1"),
			Name: "Test User New",
		})

		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
	})
}
