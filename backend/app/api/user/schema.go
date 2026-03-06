package api_user

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/healthy-heroes/neskuchka/backend/app/domain"
)

// UserSchema is the schema for data about logged user
type UserSchema struct {
	ID     string
	Name   string
	Avatar string `json:",omitempty"`
}

// SettingsSchema is the schema for the user settings page
type SettingsSchema struct {
	Name  string
	Email string
}

func MakeSettingsSchema(user domain.User) SettingsSchema {
	return SettingsSchema{
		Name:  user.Name,
		Email: string(user.Email),
	}
}

// UpdateSettingsSchema is the schema for updating user settings
type UpdateSettingsSchema struct {
	Name string `json:"Name"`
}

func (s UpdateSettingsSchema) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Name, validation.Required, validation.Length(1, 100)),
	)
}
