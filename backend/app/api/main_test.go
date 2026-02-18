package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/healthy-heroes/neskuchka/backend/app/domain"
	"github.com/healthy-heroes/neskuchka/backend/app/internal/session"
	"github.com/healthy-heroes/neskuchka/backend/app/storage/database"
)

type TestApp struct {
	Secret  string
	Version string

	Server      *httptest.Server
	DB          *database.Engine
	DataStorage *database.DataStorage
	Store       *domain.Store

	SessionManager *session.Manager
}

func NewTestApp(t *testing.T) *TestApp {
	t.Helper()

	app := &TestApp{
		Secret:  "test_secret",
		Version: "test_version",
	}

	engine, err := database.NewSqliteEngine(":memory:", zerolog.Nop())
	require.NoError(t, err)

	app.SessionManager = session.NewManager(session.Opts{
		Logger: zerolog.Nop(),
		Issuer: Issuer,
		Secret: app.Secret,
	})

	app.DB = engine
	app.DataStorage = database.NewDataStorage(engine, zerolog.Nop())
	app.Store = domain.NewStore(domain.Opts{
		DataStorage: app.DataStorage,
	})

	api := &Api{
		Version:   app.Version,
		Secret:    app.Secret,
		DataStore: app.Store,

		WebFS: fstest.MapFS{
			"web/index.html": &fstest.MapFile{Data: []byte("<html>test</html>")},
		},
	}

	app.Server = httptest.NewServer(api.Handler())
	t.Cleanup(func() {
		app.DB.Close()
		app.Server.Close()
	})

	return app
}

func (app *TestApp) LoginAs(t *testing.T, uid domain.UserID) *http.Cookie {
	t.Helper()

	w := httptest.NewRecorder()
	err := app.SessionManager.Set(w, string(uid))
	require.NoError(t, err)

	cookies := w.Result().Cookies()
	require.NotEmpty(t, cookies)

	return cookies[0]
}

type RequestOption func(*http.Request)

func WithCookie(c *http.Cookie) RequestOption {
	return func(r *http.Request) {
		r.AddCookie(c)
	}
}

func (app *TestApp) GET(t *testing.T, path string, opts ...RequestOption) *http.Response {
	t.Helper()

	url := app.Server.URL + path
	req, err := http.NewRequest("GET", url, nil)
	require.NoError(t, err)

	for _, patch := range opts {
		patch(req)
	}

	resp, err := app.Server.Client().Do(req)
	require.NoError(t, err)

	t.Cleanup(func() {
		resp.Body.Close()
	})

	return resp
}

// ReadBody reads the body of a response and returns it as a string
func ReadBody(t *testing.T, resp *http.Response) string {
	t.Helper()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return string(body)
}

type JSONResponse[T any] struct {
	Data T `json:"data"`
}

func ReadJSON[T any](t *testing.T, resp *http.Response) T {
	t.Helper()

	var result JSONResponse[T]
	err := json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	return result.Data
}

func Test_Api(t *testing.T) {
	app := NewTestApp(t)

	t.Run("/ping", func(t *testing.T) {
		resp := app.GET(t, "/ping")

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "pong", ReadBody(t, resp))
	})

	t.Run("/index.html", func(t *testing.T) {
		resp := app.GET(t, "/index.html")

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "<html>test</html>", ReadBody(t, resp))
	})
}
