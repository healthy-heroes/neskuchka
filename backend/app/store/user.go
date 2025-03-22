package store

import "github.com/google/uuid"

type UserID string

type User struct {
	ID    UserID
	Name  string
	Login string
	Email string
}

func CreateUserId() UserID {
	return UserID(uuid.New().String())
}

type UserStore interface {
	Create(user *User) (*User, error)
	Get(id UserID) (*User, error)
}
