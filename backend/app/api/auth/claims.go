package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/healthy-heroes/neskuchka/backend/app/internal/token"
)

const (
	ClaimsTypeRegistration token.ClaimsType = "registration"
)

type RegistrationClaims struct {
	jwt.RegisteredClaims

	Data       *UserRegistrationSchema `json:"data"`
	ClaimsType token.ClaimsType        `json:"type"`
}

func NewRegistrationClaims(data *UserRegistrationSchema, regClaims jwt.RegisteredClaims) RegistrationClaims {
	return RegistrationClaims{
		ClaimsType:       ClaimsTypeRegistration,
		Data:             data,
		RegisteredClaims: regClaims,
	}
}

func (r RegistrationClaims) GetType() token.ClaimsType {
	return r.ClaimsType
}
