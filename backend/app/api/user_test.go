package api

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"math/rand/v2"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/healthy-heroes/neskuchka/backend/app/domain"
	"github.com/healthy-heroes/neskuchka/backend/app/internal/avatarimg"
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

func Test_ApiUserService_GetSettings(t *testing.T) {
	app := NewTestApp(t)

	t.Run("should return settings with email", func(t *testing.T) {
		user, err := app.DataStorage.CreateUser(t.Context(), testutil.CreateUser())
		require.NoError(t, err)

		resp := app.GET(t, "/api/v1/user/me/settings", WithCookie(app.LoginAs(t, user.ID)))
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		type settingsResp struct {
			Name  string
			Email string
		}
		data := ReadJSON[settingsResp](t, resp)

		assert.Equal(t, settingsResp{user.Name, string(user.Email)}, data)
	})

	t.Run("should return 401 if user is not logged in", func(t *testing.T) {
		resp := app.GET(t, "/api/v1/user/me/settings")
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func Test_ApiUserService_UpdateSettings(t *testing.T) {
	app := NewTestApp(t)

	t.Run("should update user name", func(t *testing.T) {
		user, err := app.DataStorage.CreateUser(t.Context(), testutil.CreateUser())
		require.NoError(t, err)

		resp := app.PUT(t, "/api/v1/user/me/settings",
			WithCookie(app.LoginAs(t, user.ID)),
			WithJSON(map[string]string{"Name": "New Name"}),
		)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		type settingsResp struct {
			Name  string
			Email string
		}
		data := ReadJSON[settingsResp](t, resp)

		assert.Equal(t, "New Name", data.Name)
		assert.Equal(t, string(user.Email), data.Email)
	})

	t.Run("should return 422 when name is empty", func(t *testing.T) {
		user, err := app.DataStorage.CreateUser(t.Context(), testutil.CreateUser())
		require.NoError(t, err)

		resp := app.PUT(t, "/api/v1/user/me/settings",
			WithCookie(app.LoginAs(t, user.ID)),
			WithJSON(map[string]string{"Name": ""}),
		)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("should return 401 if user is not logged in", func(t *testing.T) {
		resp := app.PUT(t, "/api/v1/user/me/settings",
			WithJSON(map[string]string{"Name": "Test"}),
		)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func Test_ApiUserService_DeleteAvatar(t *testing.T) {
	app := NewTestApp(t)

	t.Run("should delete avatar", func(t *testing.T) {
		user, err := app.DataStorage.CreateUser(t.Context(), testutil.CreateUser())
		require.NoError(t, err)

		err = app.AvatarStorage.Save(t.Context(), user.ID, domain.Avatar{
			MimeType: "image/png",
			Data:     []byte("test"),
		})
		require.NoError(t, err)

		resp := app.DELETE(t, "/api/v1/user/me/avatar", WithCookie(app.LoginAs(t, user.ID)))
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		exists, err := app.AvatarStorage.Exists(t.Context(), user.ID)
		require.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("should succeed even if no avatar exists", func(t *testing.T) {
		userID := domain.NewUserID()

		resp := app.DELETE(t, "/api/v1/user/me/avatar", WithCookie(app.LoginAs(t, userID)))
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("should return 401 if user is not logged in", func(t *testing.T) {
		resp := app.DELETE(t, "/api/v1/user/me/avatar")
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

	// makeNoisePNG builds a PNG that does not compress, so that its encoded
	// size is roughly w*h*3 bytes and the upload limit can be exercised.
	makeNoisePNG := func(t *testing.T, w, h int) []byte {
		t.Helper()
		rnd := rand.New(rand.NewPCG(1, 2))
		img := image.NewRGBA(image.Rect(0, 0, w, h))
		for y := range h {
			for x := range w {
				img.Set(x, y, color.RGBA{
					R: uint8(rnd.UintN(256)),
					G: uint8(rnd.UintN(256)),
					B: uint8(rnd.UintN(256)),
					A: 255,
				})
			}
		}
		buf := new(bytes.Buffer)
		require.NoError(t, png.Encode(buf, img))
		return buf.Bytes()
	}

	decodeSaved := func(t *testing.T, data []byte) image.Image {
		t.Helper()
		img, format, err := image.Decode(bytes.NewReader(data))
		require.NoError(t, err)
		require.Equal(t, "png", format)
		return img
	}

	t.Run("should downscale an uploaded png avatar", func(t *testing.T) {
		userID := domain.NewUserID()

		resp := app.POST(t, "/api/v1/user/me/avatar",
			WithCookie(app.LoginAs(t, userID)),
			WithMultipartFile("avatar", "photo.png", "image/png", makePNG(t, 1000, 1000)),
		)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		saved, err := app.AvatarStorage.Get(t.Context(), userID)
		require.NoError(t, err)
		assert.Equal(t, "image/png", saved.MimeType)

		img := decodeSaved(t, saved.Data)
		assert.Equal(t, avatarimg.MaxSide, img.Bounds().Dx())
		assert.Equal(t, avatarimg.MaxSide, img.Bounds().Dy())
	})

	t.Run("should store an uploaded jpeg avatar as png", func(t *testing.T) {
		userID := domain.NewUserID()

		resp := app.POST(t, "/api/v1/user/me/avatar",
			WithCookie(app.LoginAs(t, userID)),
			WithMultipartFile("avatar", "photo.jpg", "image/jpeg", makeJPEG(t, 80, 80)),
		)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		saved, err := app.AvatarStorage.Get(t.Context(), userID)
		require.NoError(t, err)
		assert.Equal(t, "image/png", saved.MimeType)

		// 80x80 is below MaxSide, so the picture keeps its size and only the
		// encoding changes.
		img := decodeSaved(t, saved.Data)
		assert.Equal(t, 80, img.Bounds().Dx())
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

	t.Run("should return 400 for an image that fails to decode", func(t *testing.T) {
		userID := domain.NewUserID()

		// Passes the content type sniffing, falls apart on decoding.
		broken := append([]byte("\x89PNG\r\n\x1a\n"), []byte("truncated junk")...)

		resp := app.POST(t, "/api/v1/user/me/avatar",
			WithCookie(app.LoginAs(t, userID)),
			WithMultipartFile("avatar", "photo.png", "image/png", broken),
		)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should accept a photo weighing more than a megabyte", func(t *testing.T) {
		userID := domain.NewUserID()

		photo := makeNoisePNG(t, 900, 900)
		require.Greater(t, len(photo), 1024*1024, "fixture is not big enough to test the limit")

		resp := app.POST(t, "/api/v1/user/me/avatar",
			WithCookie(app.LoginAs(t, userID)),
			WithMultipartFile("avatar", "photo.png", "image/png", photo),
		)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
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
