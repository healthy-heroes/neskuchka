package session

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestManager() *Manager {
	return NewManager(Opts{
		Issuer:          "test_issuer",
		Secret:          "test_secret",
		SecureCookies:   false,
		SessionDuration: time.Hour,
	})
}

func TestManager_Set(t *testing.T) {
	t.Run("sets valid JWT cookie", func(t *testing.T) {
		m := newTestManager()
		w := httptest.NewRecorder()

		userID := "user-123"

		err := m.Set(w, userID)

		require.NoError(t, err)

		cookies := w.Result().Cookies()
		require.Len(t, cookies, 1, "should set exactly one cookie")

		cookie := cookies[0]
		assert.Equal(t, defaultSessionCookieName, cookie.Name)
		assert.NotEmpty(t, cookie.Value)
		assert.True(t, cookie.HttpOnly)
		assert.Equal(t, "/", cookie.Path)
		assert.Equal(t, int(time.Hour.Seconds()), cookie.MaxAge)
		assert.Equal(t, http.SameSiteLaxMode, cookie.SameSite)
	})

	t.Run("token contains correct user data", func(t *testing.T) {
		m := newTestManager()
		w := httptest.NewRecorder()

		userID := "user-456"

		err := m.Set(w, userID)
		require.NoError(t, err)

		cookie := w.Result().Cookies()[0]

		// Parse the token and verify claims
		claims := Claims{}
		err = m.tokenService.Parse(cookie.Value, &claims)
		require.NoError(t, err)

		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, m.opts.Issuer, claims.Issuer)
		assert.NotEmpty(t, claims.ID, "JTI should be set")
	})

	t.Run("respects SecureCookies option", func(t *testing.T) {
		m := newTestManager()
		m.opts.SecureCookies = true
		w := httptest.NewRecorder()

		err := m.Set(w, "user-1")
		require.NoError(t, err)

		cookie := w.Result().Cookies()[0]
		assert.True(t, cookie.Secure)
	})
}

func TestManager_Clear(t *testing.T) {
	t.Run("clears JWT cookie", func(t *testing.T) {
		m := newTestManager()
		w := httptest.NewRecorder()

		m.Clear(w)

		cookies := w.Result().Cookies()
		require.Len(t, cookies, 1, "should set exactly one cookie")

		cookie := cookies[0]
		assert.Equal(t, defaultSessionCookieName, cookie.Name)
		assert.Empty(t, cookie.Value, "cookie value should be empty")
		assert.Equal(t, -1, cookie.MaxAge, "MaxAge should be -1 to delete cookie")
		assert.True(t, cookie.HttpOnly)
		assert.Equal(t, "/", cookie.Path)
		assert.Equal(t, http.SameSiteLaxMode, cookie.SameSite)
	})

	t.Run("respects SecureCookies option", func(t *testing.T) {
		m := newTestManager()
		m.opts.SecureCookies = true
		w := httptest.NewRecorder()

		m.Clear(w)

		cookie := w.Result().Cookies()[0]
		assert.True(t, cookie.Secure)
	})
}

func TestManager_Get(t *testing.T) {
	m := newTestManager()

	t.Run("extracts user from valid token", func(t *testing.T) {
		expectedUserID := "user-123"

		claims := Claims{
			UserID: expectedUserID,
			RegisteredClaims: jwt.RegisteredClaims{
				ID:        "test_jti",
				Issuer:    "test_issuer",
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			},
		}

		tokenString, err := m.tokenService.Token(claims)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{
			Name:  defaultSessionCookieName,
			Value: tokenString,
		})

		userID, parsedClaims, err := m.Get(req)

		require.NoError(t, err)
		assert.Equal(t, expectedUserID, userID, "user ID should match")
		assert.Equal(t, expectedUserID, parsedClaims.UserID)
	})

	t.Run("returns error when cookie is missing", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)

		_, _, err := m.Get(req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "token cookie was not presented")
	})

	t.Run("returns error for invalid token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{
			Name:  defaultSessionCookieName,
			Value: "invalid.token.here",
		})

		_, _, err := m.Get(req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse token")
	})
}
