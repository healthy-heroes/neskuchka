package api_user

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// UserSchema is the schema for data about logged user
type UserSchema struct {
	ID     string
	Name   string
	Avatar string `json:",omitempty"`
}

// SettingsSchema is the schema for the user settings page
type SettingsSchema struct {
	Name   string
	Email  string
	Avatar string `json:",omitempty"`
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
