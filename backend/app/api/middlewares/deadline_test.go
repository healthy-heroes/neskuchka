package middlewares

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// readingServer serves a body-reading handler over a real connection, which is
// the point: a deadline needs one, and httptest's recorder has none.
func readingServer(t *testing.T, d time.Duration) (addr string, readErr <-chan error) {
	t.Helper()

	errs := make(chan error, 1)

	handler := ReadDeadline(zerolog.Nop(), d)(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			_, err := io.ReadAll(r.Body)
			errs <- err
			w.WriteHeader(http.StatusOK)
		}))

	server := httptest.NewUnstartedServer(handler)
	server.Config.ErrorLog = nil
	server.Start()
	t.Cleanup(server.Close)

	return server.Listener.Addr().String(), errs
}

func TestReadDeadline(t *testing.T) {
	t.Run("gives up on a body that stops coming", func(t *testing.T) {
		addr, readErr := readingServer(t, 100*time.Millisecond)

		// A body promised in full and then abandoned — indistinguishable from a
		// client that means to send it one byte at a time.
		conn, err := net.Dial("tcp", addr)
		require.NoError(t, err)
		t.Cleanup(func() { conn.Close() })

		_, err = fmt.Fprint(conn, "POST / HTTP/1.1\r\nHost: test\r\nContent-Length: 1000\r\n\r\na")
		require.NoError(t, err)

		select {
		case err := <-readErr:
			require.Error(t, err, "the handler read a body that was never sent")
			assert.ErrorIs(t, err, os.ErrDeadlineExceeded)
		case <-time.After(time.Second):
			t.Fatal("the handler is still waiting for a body nobody is sending")
		}
	})

	t.Run("leaves a request that arrives in time alone", func(t *testing.T) {
		addr, readErr := readingServer(t, time.Minute)

		resp, err := http.Post("http://"+addr+"/", "text/plain", strings.NewReader("hello"))
		require.NoError(t, err)
		t.Cleanup(func() { resp.Body.Close() })

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		select {
		case err := <-readErr:
			assert.NoError(t, err)
		case <-time.After(time.Second):
			t.Fatal("the handler never finished reading")
		}
	})
}
