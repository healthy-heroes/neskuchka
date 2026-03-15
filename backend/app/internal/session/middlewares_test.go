package session

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVerifier(t *testing.T) {
	m := newTestManager()

	execMw := func(t *testing.T, req *http.Request) (string, bool) {
		t.Helper()

		var capturedUserID string
		var hasUserID bool
		handler := m.Verifier()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			capturedUserID, hasUserID = GetUserID(r)
			w.WriteHeader(http.StatusOK)
		}))

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		return capturedUserID, hasUserID
	}

	t.Run("sets userID in context when valid bearer token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", createToken(t, m, "user-123")))

		userID, ok := execMw(t, req)
		assert.True(t, ok)
		assert.Equal(t, "user-123", userID)
	})

	t.Run("gets from bearer first, then from cookie", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", createToken(t, m, "user-123")))
		req.AddCookie(&http.Cookie{
			Name:  defaultSessionCookieName,
			Value: createToken(t, m, "user-456"),
		})

		userID, ok := execMw(t, req)
		assert.True(t, ok)
		assert.Equal(t, "user-123", userID)
	})

	t.Run("no using cookie if bearer token is invalid", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", "invalid.token.here"))
		req.AddCookie(&http.Cookie{
			Name:  defaultSessionCookieName,
			Value: createToken(t, m, "user-456"),
		})

		_, ok := execMw(t, req)
		assert.False(t, ok)
	})

	t.Run("sets userID in context when valid cookie token", func(t *testing.T) {
		w := httptest.NewRecorder()
		err := m.Set(w, "user-123")
		require.NoError(t, err)
		cookie := w.Result().Cookies()[0]

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(cookie)

		userID, ok := execMw(t, req)
		assert.True(t, ok)
		assert.Equal(t, "user-123", userID)
	})

	t.Run("continues without userID when no cookie", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)

		_, ok := execMw(t, req)
		assert.False(t, ok)
	})

	t.Run("continues without userID when invalid token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{
			Name:  defaultSessionCookieName,
			Value: "invalid.token.here",
		})

		_, ok := execMw(t, req)
		assert.False(t, ok)
	})
}

func TestAuthenticator(t *testing.T) {
	m := newTestManager()

	t.Run("calls next handler when userID in context", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		ctx := context.WithValue(req.Context(), ctxKey, "user-123")
		req = req.WithContext(ctx)

		nextCalled := false
		unauthorizedCalled := false

		handler := m.Authenticator(func(w http.ResponseWriter) {
			unauthorizedCalled = true
		})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			nextCalled = true
			w.WriteHeader(http.StatusOK)
		}))

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		assert.True(t, nextCalled)
		assert.False(t, unauthorizedCalled)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("calls unauthorizedHandler when no userID in context", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)

		nextCalled := false
		unauthorizedCalled := false

		handler := m.Authenticator(func(w http.ResponseWriter) {
			unauthorizedCalled = true
			w.WriteHeader(http.StatusUnauthorized)
		})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			nextCalled = true
		}))

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		assert.False(t, nextCalled)
		assert.True(t, unauthorizedCalled)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("calls unauthorizedHandler when userID is empty", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		ctx := context.WithValue(req.Context(), ctxKey, "")
		req = req.WithContext(ctx)

		unauthorizedCalled := false

		handler := m.Authenticator(func(w http.ResponseWriter) {
			unauthorizedCalled = true
			w.WriteHeader(http.StatusUnauthorized)
		})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		assert.True(t, unauthorizedCalled)
	})
}

func TestGetUserID(t *testing.T) {
	t.Run("returns userID when present", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		ctx := context.WithValue(req.Context(), ctxKey, "user-456")
		req = req.WithContext(ctx)

		userID, ok := GetUserID(req)

		assert.True(t, ok)
		assert.Equal(t, "user-456", userID)
	})

	t.Run("returns false when not present", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)

		userID, ok := GetUserID(req)

		assert.False(t, ok)
		assert.Empty(t, userID)
	})
}

func TestMustGetUserID(t *testing.T) {
	t.Run("returns userID when present", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		ctx := context.WithValue(req.Context(), ctxKey, "user-789")
		req = req.WithContext(ctx)

		userID := MustGetUserID(req)

		assert.Equal(t, "user-789", userID)
	})

	t.Run("panics when not present", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)

		assert.Panics(t, func() {
			MustGetUserID(req)
		})
	})
}

func createToken(t *testing.T, m *Manager, userID string) string {
	t.Helper()

	tokenString, err := m.Token(userID)
	require.NoError(t, err)
	return tokenString
}
