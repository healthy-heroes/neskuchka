package api

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/healthy-heroes/neskuchka/backend/app/domain"
	"github.com/healthy-heroes/neskuchka/backend/app/internal/testutil"
)

func Test_ApiUserService_User(t *testing.T) {
	app := NewTestApp(t)

	t.Run("should returns current user", func(t *testing.T) {
		user, err := app.DataStorage.CreateUser(t.Context(), testutil.CreateUser())
		require.NoError(t, err)

		resp := app.GET(t, "/api/v1/user/me", WithCookie(app.LoginAs(t, user.ID)))
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		type userResp struct {
			ID   string
			Name string
		}

		data := ReadJSON[userResp](t, resp)

		assert.Equal(t, userResp{string(user.ID), user.Name}, data)
	})

	t.Run("should return avatar url if avatar exists", func(t *testing.T) {
		user, err := app.DataStorage.CreateUser(t.Context(), testutil.CreateUser())
		require.NoError(t, err)

		err = app.AvatarStorage.Save(t.Context(), user.ID, domain.Avatar{
			MimeType: "image/png",
			Data:     []byte("test"),
		})
		require.NoError(t, err)

		resp := app.GET(t, "/api/v1/user/me", WithCookie(app.LoginAs(t, user.ID)))
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		type userResp struct {
			ID     string
			Name   string
			Avatar string
		}
		data := ReadJSON[userResp](t, resp)

		assert.Equal(t,
			userResp{
				string(user.ID),
				user.Name,
				fmt.Sprintf("%s/user/%s/avatar", prefixApi, string(user.ID)),
			},
			data,
		)
	})

	t.Run("should return 401 if user is not logged in", func(t *testing.T) {
		resp := app.GET(t, "/api/v1/user/me")
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 401 if user is not found", func(t *testing.T) {
		resp := app.GET(t, "/api/v1/user/me", WithCookie(app.LoginAs(t, domain.NewUserID())))
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func Test_ApiUserService_MyAvatar(t *testing.T) {
	app := NewTestApp(t)

	t.Run("should return avatar bytes", func(t *testing.T) {
		userID := domain.NewUserID()

		err := app.AvatarStorage.Save(t.Context(), userID, domain.Avatar{
			MimeType: "image/jpeg",
			Data:     []byte("test"),
		})
		require.NoError(t, err)

		resp := app.GET(t, "/api/v1/user/me/avatar", WithCookie(app.LoginAs(t, userID)))
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "image/jpeg", resp.Header.Get("Content-Type"))
		assert.Equal(t, "test", ReadBody(t, resp))
	})

	t.Run("should return 401 if user is not logged in", func(t *testing.T) {
		resp := app.GET(t, "/api/v1/user/me/avatar")
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func Test_ApiUserService_UserAvatar(t *testing.T) {
	app := NewTestApp(t)

	t.Run("should return avatar bytes", func(t *testing.T) {
		userID := domain.NewUserID()
		err := app.AvatarStorage.Save(t.Context(), userID, domain.Avatar{
			MimeType: "image/jpeg",
			Data:     []byte("test"),
		})
		require.NoError(t, err)

		resp := app.GET(t, fmt.Sprintf("/api/v1/user/%s/avatar", string(userID)))
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "image/jpeg", resp.Header.Get("Content-Type"))
		assert.Equal(t, "test", ReadBody(t, resp))
	})

	t.Run("should return 404 if avatar does not exist", func(t *testing.T) {
		resp := app.GET(t, "/api/v1/user/1/avatar")
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func Test_ApiUserService_UploadAvatar(t *testing.T) {
	app := NewTestApp(t)

	makePNG := func(t *testing.T, w, h int) []byte {
		t.Helper()
		img := image.NewRGBA(image.Rect(0, 0, w, h))
		buf := new(bytes.Buffer)
		require.NoError(t, png.Encode(buf, img))
		return buf.Bytes()
	}

	makeJPEG := func(t *testing.T, w, h int) []byte {
		t.Helper()
		img := image.NewRGBA(image.Rect(0, 0, w, h))
		buf := new(bytes.Buffer)
		require.NoError(t, jpeg.Encode(buf, img, nil))
		return buf.Bytes()
	}

	t.Run("should upload png avatar", func(t *testing.T) {
		userID := domain.NewUserID()
		pngData := makePNG(t, 100, 100)

		resp := app.POST(t, "/api/v1/user/me/avatar",
			WithCookie(app.LoginAs(t, userID)),
			WithMultipartFile("avatar", "photo.png", "image/png", pngData),
		)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		saved, err := app.AvatarStorage.Get(t.Context(), userID)
		require.NoError(t, err)
		assert.Equal(t, "image/png", saved.MimeType)
		assert.Equal(t, pngData, saved.Data)
	})

	t.Run("should upload jpeg avatar", func(t *testing.T) {
		userID := domain.NewUserID()
		jpegData := makeJPEG(t, 80, 80)

		resp := app.POST(t, "/api/v1/user/me/avatar",
			WithCookie(app.LoginAs(t, userID)),
			WithMultipartFile("avatar", "photo.jpg", "image/jpeg", jpegData),
		)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		saved, err := app.AvatarStorage.Get(t.Context(), userID)
		require.NoError(t, err)
		assert.Equal(t, "image/jpeg", saved.MimeType)
		assert.Equal(t, jpegData, saved.Data)
	})

	t.Run("should return 401 if user is not logged in", func(t *testing.T) {
		resp := app.POST(t, "/api/v1/user/me/avatar",
			WithMultipartFile("avatar", "photo.png", "image/png", makePNG(t, 10, 10)),
		)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 400 for unsupported mime type in header", func(t *testing.T) {
		userID := domain.NewUserID()

		resp := app.POST(t, "/api/v1/user/me/avatar",
			WithCookie(app.LoginAs(t, userID)),
			WithMultipartFile("avatar", "doc.pdf", "application/pdf", []byte("%PDF-1.4")),
		)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 when content doesn't match allowed image types", func(t *testing.T) {
		userID := domain.NewUserID()

		resp := app.POST(t, "/api/v1/user/me/avatar",
			WithCookie(app.LoginAs(t, userID)),
			WithMultipartFile("avatar", "fake.png", "image/png", []byte("not an image at all")),
		)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 when file field is missing", func(t *testing.T) {
		userID := domain.NewUserID()

		resp := app.POST(t, "/api/v1/user/me/avatar",
			WithCookie(app.LoginAs(t, userID)),
			WithMultipartFile("wrong_field", "photo.png", "image/png", makePNG(t, 10, 10)),
		)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
