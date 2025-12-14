package token

import (
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

type TestClaims struct {
	jwt.RegisteredClaims

	Data struct {
		Name string
	}
}

func TestJWT_New(t *testing.T) {
	j := NewService(Opts{
		Issuer: issuer,
		Secret: secret,
	})

	assert.NotNil(t, j)
	assert.Equal(t, issuer, j.Issuer)
	assert.Equal(t, secret, j.Secret)
}

func TestJWT_RandID(t *testing.T) {
	id, err := RandID()
	assert.NoError(t, err)
	assert.NotEmpty(t, id)

	// Test uniqueness
	id2, err := RandID()
	assert.NoError(t, err)
	assert.NotEqual(t, id, id2)
}

func TestJWT_TokenAndParse(t *testing.T) {
	j := NewService(Opts{
		Issuer: issuer,
		Secret: secret,
	})

	claims := testClaims()
	token, err := j.Token(claims)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	parsedClaims := TestClaims{}
	err = j.Parse(token, &parsedClaims)
	assert.NoError(t, err)
	assert.Equal(t, claims, parsedClaims)
}

func TestJWT_Token(t *testing.T) {
	t.Run("empty secret", func(t *testing.T) {
		j := NewService(Opts{
			Issuer: issuer,
			Secret: "",
		})
		_, err := j.Token(testClaims())
		assert.Error(t, err)
	})
}

func TestJWT_Parse(t *testing.T) {
	t.Run("expired within leeway", func(t *testing.T) {
		j := NewService(Opts{Issuer: issuer, Secret: secret})

		claims := TestClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(-30 * time.Second)), // expired 30 sec ago
				Issuer:    issuer,
				ID:        jti,
			},
			Data: struct{ Name string }{Name: issuer},
		}
		token, err := j.Token(claims)
		assert.NoError(t, err)

		err = j.Parse(token, &TestClaims{})
		assert.NoError(t, err)
	})

	t.Run("expired beyond leeway", func(t *testing.T) {
		j := NewService(Opts{Issuer: issuer, Secret: secret})

		claims := TestClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(-2 * time.Minute)), // expired 2 min ago
				Issuer:    issuer,
				ID:        jti,
			},
			Data: struct{ Name string }{Name: issuer},
		}
		token, err := j.Token(claims)
		assert.NoError(t, err)

		err = j.Parse(token, &TestClaims{})
		assert.Error(t, err)
	})

	// "None" Algorithm Attack (CVE-2015-9235)
	// Attacker sets alg: none and removes signature, bypassing verification.
	// jwt/v5 rejects "none" by default, but we verify this protection is in place.
	t.Run("none algorithm", func(t *testing.T) {
		j := NewService(Opts{Issuer: issuer, Secret: secret})

		token := jwt.NewWithClaims(jwt.SigningMethodNone, testClaims())
		tokenString, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
		assert.NoError(t, err)

		err = j.Parse(tokenString, &TestClaims{})
		assert.Error(t, err)
	})

	// Algorithm Confusion Attack (CVE-2024-54150, CVE-2023-48238)
	// Attacker changes alg from RS256 to HS256 and signs with public key.
	// WithValidMethods prevents accepting unexpected algorithms.
	// This test uses HS384 (same secret) to verify strict algorithm validation.
	t.Run("hs384 algorithm not allowed", func(t *testing.T) {
		j := NewService(Opts{Issuer: issuer, Secret: secret})

		token := jwt.NewWithClaims(jwt.SigningMethodHS384, testClaims())
		tokenString, err := token.SignedString([]byte(secret))
		assert.NoError(t, err)

		err = j.Parse(tokenString, &TestClaims{})
		assert.ErrorContains(t, err, "signing method")
	})

	t.Run("wrong secret", func(t *testing.T) {
		j1 := NewService(Opts{Issuer: issuer, Secret: secret})
		j2 := NewService(Opts{Issuer: issuer, Secret: "wrong_secret"})

		token, err := j1.Token(testClaims())
		assert.NoError(t, err)

		err = j2.Parse(token, &TestClaims{})
		assert.Error(t, err)
	})

	t.Run("invalid token", func(t *testing.T) {
		j := NewService(Opts{Issuer: issuer, Secret: secret})

		err := j.Parse("not.a.valid.jwt.token", &TestClaims{})
		assert.Error(t, err)
	})

	t.Run("empty token", func(t *testing.T) {
		j := NewService(Opts{Issuer: issuer, Secret: secret})

		err := j.Parse("", &TestClaims{})
		assert.Error(t, err)
	})

	t.Run("malformed token", func(t *testing.T) {
		j := NewService(Opts{Issuer: issuer, Secret: secret})

		err := j.Parse("only.two.parts", &TestClaims{})
		assert.Error(t, err)
	})
}

func testClaims() TestClaims {
	return TestClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Date(2100, 2, 14, 12, 30, 59, 0, time.UTC).Local()),
			NotBefore: jwt.NewNumericDate(time.Date(2000, 2, 14, 12, 30, 59, 0, time.UTC).Local()),
			Issuer:    issuer,
			ID:        jti,
		},

		Data: struct{ Name string }{Name: issuer},
	}
}
