package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/healthy-heroes/neskuchka/backend/app/internal/token"
	"github.com/healthy-heroes/neskuchka/backend/app/store"
)

func newTestService() *Service {
	tokenService := token.NewService(token.Opts{
		Issuer: "test_issuer",
		Secret: "test_secret",
	})

	return &Service{
		opts: Opts{
			Issuer:          "test_issuer",
			Secret:          "test_secret",
			Logger:          zerolog.Nop(),
			SecureCookies:   false,
			SessionDuration: time.Hour,
		},
		tokenService: tokenService,
	}
}

func TestService_SetToken(t *testing.T) {
	t.Run("sets valid JWT cookie", func(t *testing.T) {
		s := newTestService()
		w := httptest.NewRecorder()

		user := UserSchema{
			ID:   store.UserID("user-123"),
			Name: "Test User",
		}

		err := s.setToken(w, user)

		require.NoError(t, err)

		cookies := w.Result().Cookies()
		require.Len(t, cookies, 1, "should set exactly one cookie")

		cookie := cookies[0]
		assert.Equal(t, JWTCookieName, cookie.Name)
		assert.NotEmpty(t, cookie.Value)
		assert.True(t, cookie.HttpOnly)
		assert.Equal(t, "/", cookie.Path)
		assert.Equal(t, int(time.Hour.Seconds()), cookie.MaxAge)
		assert.Equal(t, http.SameSiteLaxMode, cookie.SameSite)
	})

	t.Run("token contains correct user data", func(t *testing.T) {
		s := newTestService()
		w := httptest.NewRecorder()

		user := UserSchema{
			ID:   store.UserID("user-456"),
			Name: "Another User",
		}

		err := s.setToken(w, user)
		require.NoError(t, err)

		cookie := w.Result().Cookies()[0]

		// Parse the token and verify claims
		claims := UserClaims{}
		err = s.tokenService.Parse(cookie.Value, &claims)
		require.NoError(t, err)

		assert.Equal(t, user.ID, claims.Data.ID)
		assert.Equal(t, user.Name, claims.Data.Name)
		assert.Equal(t, s.opts.Issuer, claims.Issuer)
		assert.NotEmpty(t, claims.ID, "JTI should be set")
	})

	t.Run("respects SecureCookies option", func(t *testing.T) {
		s := newTestService()
		s.opts.SecureCookies = true
		w := httptest.NewRecorder()

		err := s.setToken(w, UserSchema{ID: "user-1", Name: "User"})
		require.NoError(t, err)

		cookie := w.Result().Cookies()[0]
		assert.True(t, cookie.Secure)
	})
}

func TestService_ClearToken(t *testing.T) {
	t.Run("clears JWT cookie", func(t *testing.T) {
		s := newTestService()
		w := httptest.NewRecorder()

		s.clearToken(w)

		cookies := w.Result().Cookies()
		require.Len(t, cookies, 1, "should set exactly one cookie")

		cookie := cookies[0]
		assert.Equal(t, JWTCookieName, cookie.Name)
		assert.Empty(t, cookie.Value, "cookie value should be empty")
		assert.Equal(t, -1, cookie.MaxAge, "MaxAge should be -1 to delete cookie")
		assert.True(t, cookie.HttpOnly)
		assert.Equal(t, "/", cookie.Path)
		assert.Equal(t, http.SameSiteLaxMode, cookie.SameSite)
	})

	t.Run("respects SecureCookies option", func(t *testing.T) {
		s := newTestService()
		s.opts.SecureCookies = true
		w := httptest.NewRecorder()

		s.clearToken(w)

		cookie := w.Result().Cookies()[0]
		assert.True(t, cookie.Secure)
	})
}

func TestService_GetUser(t *testing.T) {
	s := newTestService()

	t.Run("extracts user from valid token", func(t *testing.T) {
		expectedUser := UserSchema{
			ID:   store.UserID("user-123"),
			Name: "Test User",
		}

		claims := UserClaims{
			Data: expectedUser,
			RegisteredClaims: jwt.RegisteredClaims{
				ID:        "test_jti",
				Issuer:    "test_issuer",
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			},
		}

		tokenString, err := s.tokenService.Token(claims)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{
			Name:  JWTCookieName,
			Value: tokenString,
		})

		user, err := s.getUser(req)

		require.NoError(t, err)
		assert.Equal(t, expectedUser.ID, user.ID, "user ID should match")
		assert.Equal(t, expectedUser.Name, user.Name, "user Name should match")
	})

	t.Run("returns error when cookie is missing", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)

		_, err := s.getUser(req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "token cookie was not presented")
	})

	t.Run("returns error for invalid token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{
			Name:  JWTCookieName,
			Value: "invalid.token.here",
		})

		_, err := s.getUser(req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse token")
	})
}
