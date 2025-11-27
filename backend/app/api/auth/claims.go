package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

type RegistrationClaims struct {
	jwt.RegisteredClaims

	Data *UserRegistrationSchema `json:"data"`
}

type AccessClaims struct {
	jwt.RegisteredClaims

	Data *UserSchema
}
