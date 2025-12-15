package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

// ConfirmationClaims is the claims for confirmation token
type ConfirmationClaims struct {
	jwt.RegisteredClaims

	Data LoginSchema `json:"data"`
}
