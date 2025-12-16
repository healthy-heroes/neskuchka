package auth

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"

	"github.com/healthy-heroes/neskuchka/backend/app/store"
)

// LoginSchema is the schema for the login request
type LoginSchema struct {
	Email string `json:"email"`
}

func (s LoginSchema) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Email, validation.Required, is.Email),
	)
}

// ConfirmSchema is the schema for the confirm request
type ConfirmationSchema struct {
	Token string `json:"token"`
}

// UserSchema is the schema for data about logged user
type UserSchema struct {
	ID   store.UserID
	Name string
}
