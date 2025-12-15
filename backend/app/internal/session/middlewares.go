package session

import (
	"context"
	"errors"
	"net/http"
)

// key is an unexported type for keys defined in this package.
// This prevents collisions with keys defined in other packages.
type key int

// ctxKey is the key for UserID values in Contexts. It is
// unexported; clients use Verifier middleware and GetUserID function
// instead of using this key directly.
var ctxKey key

// UnauthorizedHandler is a function that will be called if the user is not authenticated.
type UnauthorizedHandler func(w http.ResponseWriter)

// Verifier http middleware handler will verify jwt token from a http request.
// Sets userID to the request context.
func (m *Manager) Verifier() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			userID, _, err := m.Get(r)
			if err != nil {
				if !errors.Is(err, ErrSessionNotFound) {
					m.logger.Error().Err(err).Msg("failed to verify session")
				}

				next.ServeHTTP(w, r)
				return

			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, ctxKey, userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

// Authenticator checks that userID exists in context.
// Returns 401 if not found. Must be used AFTER Verifier.
func (m *Manager) Authenticator(unauthorizedHandler UnauthorizedHandler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := GetUserID(r)
			if !ok || userID == "" {
				unauthorizedHandler(w)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// GetUserID returns the userID value stored in ctx, if any.
func GetUserID(r *http.Request) (string, bool) {
	userID, ok := r.Context().Value(ctxKey).(string)

	return userID, ok
}

// MustGetUserID gets user ID and panics if can't extract it from the request.
// should be called from authenticated controllers only
func MustGetUserID(r *http.Request) string {
	userID, ok := GetUserID(r)
	if !ok {
		panic(errors.New("failed to get user ID"))
	}
	return userID
}
