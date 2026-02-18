package api

import (
	"net/http"
	"testing"

	"github.com/healthy-heroes/neskuchka/backend/app/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ApiAuthService_User(t *testing.T) {
	app := NewTestApp(t)

	t.Run("should returns current user", func(t *testing.T) {
		user, err := app.DataStorage.CreateUser(t.Context(), domain.User{
			ID:    domain.NewUserID(),
			Name:  "Test name",
			Email: "test@example.com",
		})
		require.NoError(t, err)

		resp := app.GET(t, "/api/v1/auth/user", WithCookie(app.LoginAs(t, user.ID)))
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		type userResp struct {
			ID   string
			Name string
		}

		data := ReadJSON[userResp](t, resp)

		assert.Equal(t, userResp{string(user.ID), user.Name}, data)
	})

	t.Run("should return 401 if user is not logged in", func(t *testing.T) {
		resp := app.GET(t, "/api/v1/auth/user")
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 401 if user is not found", func(t *testing.T) {
		resp := app.GET(t, "/api/v1/auth/user", WithCookie(app.LoginAs(t, domain.NewUserID())))
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
