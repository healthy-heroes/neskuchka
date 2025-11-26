package token

import (
	"net/http"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

const (
	issuer = "custom_iss"
	secret = "test_secret"
	jti    = "test_jti"
)

type Claims struct {
	jwt.RegisteredClaims

	Data struct {
		name string
	}
}

func TestJWT_New(t *testing.T) {
	j := NewService(Opts{
		Issuer:         issuer,
		Secret:         secret,
		TokenDuration:  time.Minute,
		CookieDuration: time.Hour,
		SameSite:       http.SameSiteDefaultMode,
	})

	assert.NotNil(t, j)
	assert.Equal(t, issuer, j.Issuer)
	assert.Equal(t, secret, j.Secret)
	assert.Equal(t, time.Minute, j.TokenDuration)
	assert.Equal(t, time.Hour, j.CookieDuration)
	assert.Equal(t, http.SameSiteDefaultMode, j.SameSite)
}

func TestJWT_RandID(t *testing.T) {
	id, err := RandID()
	assert.NoError(t, err)
	assert.NotEmpty(t, id)
}

func TestJWT_Token(t *testing.T) {
	j := NewService(Opts{
		Issuer:         issuer,
		Secret:         secret,
		TokenDuration:  time.Minute,
		CookieDuration: time.Hour,
		SameSite:       http.SameSiteDefaultMode,
	})

	claims := testClaims()
	token, err := j.Token(claims)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// todo Parse
	// check idempotent (save result)
}

func testClaims() Claims {
	return Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Date(2100, 2, 14, 12, 30, 59, 0, time.UTC).Local()),
			NotBefore: jwt.NewNumericDate(time.Date(2000, 2, 14, 12, 30, 59, 0, time.UTC).Local()),
			Issuer:    issuer,
			ID:        "jti",
		},

		Data: struct{ name string }{name: issuer},
	}
}
