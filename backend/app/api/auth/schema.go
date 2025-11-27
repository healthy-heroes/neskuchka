package auth

import (
	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/healthy-heroes/neskuchka/backend/app/store"
)

type UserRegistrationSchema struct {
	Name  string
	Email string
}

type UserSchema struct {
	ID   store.UserID
	Name string
}

func (u UserRegistrationSchema) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Name, validation.Required, validation.Length(1, 100)),
		validation.Field(&u.Email, validation.Required, is.Email),
	)
}
