package auth

import (
	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type UserRegistrationSchema struct {
	Name  string
	Email string
}

func (u UserRegistrationSchema) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Name, validation.Required, validation.Length(1, 100)),
		validation.Field(&u.Email, validation.Required, is.Email),
	)
}
