package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/healthy-heroes/neskuchka/backend/app/store"
)

type RegistrationClaims struct {
	jwt.RegisteredClaims

	Data *UserRegistrationSchema `json:"data"`
}

type AccessClaims struct {
	jwt.RegisteredClaims

	Data *store.User
}
