package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/healthy-heroes/neskuchka/backend/app/domain"
	"github.com/healthy-heroes/neskuchka/backend/app/storage/database"
)

type TestApp struct {
	Secret  string
	Version string

	Server      *httptest.Server
	DB          *database.Engine
	DataStorage *database.DataStorage
	Store       *domain.Store
}

func NewTestApp(t *testing.T) *TestApp {
	t.Helper()

	app := &TestApp{
		Secret:  "test_secret",
		Version: "test_version",
	}

	engine, err := database.NewSqliteEngine(":memory:", zerolog.Nop())
	require.NoError(t, err)

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

func (app *TestApp) GET(t *testing.T, path string) *http.Response {
	t.Helper()

	url := app.Server.URL + "/" + path
	req, err := http.NewRequest("GET", url, nil)
	require.NoError(t, err)

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

func Test_Api(t *testing.T) {
	app := NewTestApp(t)

	t.Run("/ping", func(t *testing.T) {
		resp := app.GET(t, "ping")

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "pong", ReadBody(t, resp))
	})

	t.Run("/index.html", func(t *testing.T) {
		resp := app.GET(t, "/index.html")

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "<html>test</html>", ReadBody(t, resp))
	})
}
