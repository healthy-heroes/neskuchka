package store

import "github.com/google/uuid"

type UserID string

type User struct {
	ID      UserID
	Name    string
	Email   string
	Picture string
}

func CreateUserId() UserID {
	return UserID(uuid.New().String())
}

type UserStore interface {
	Store

	Create(user *User) (*User, error)
	Get(id UserID) (*User, error)
	FindByEmail(email string) (*User, error)
}
