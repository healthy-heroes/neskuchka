package middlewares

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	testLimit = 2
	testQueue = 1

	// testWait is a whole second because Retry-After has no finer unit: any
	// shorter wait would round to a "0" the assertion below couldn't tell from
	// a header nobody set.
	testWait = time.Second
)

// blockingHandler reports every entry on entered and stays inside until
// release is closed, so a test can hold every slot for as long as it needs.
func blockingHandler(entered chan<- struct{}, release <-chan struct{}) http.Handler {
	return http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		entered <- struct{}{}
		<-release
	})
}

// post sends a request in the background. Nothing here waits for an answer
// inline: a request the throttle fails to hold back sits in the handler until
// the test ends, and waiting for it would hang the suite instead of failing it.
func post(handler http.Handler) <-chan *httptest.ResponseRecorder {
	answered := make(chan *httptest.ResponseRecorder, 1)

	go func() {
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/", nil))
		answered <- recorder
	}()

	return answered
}

// fillSlots occupies every slot of the throttle and returns once each blocked
// request is inside the handler.
func fillSlots(t *testing.T, handler http.Handler, entered <-chan struct{}) {
	t.Helper()

	for range testLimit {
		post(handler)
	}

	for range testLimit {
		select {
		case <-entered:
		case <-time.After(time.Second):
			t.Fatal("a request never reached the handler")
		}
	}
}

func TestThrottle(t *testing.T) {
	t.Run("turns a request away once the wait runs out", func(t *testing.T) {
		entered := make(chan struct{}, testLimit)
		release := make(chan struct{})
		defer close(release)

		handler := Throttle(testLimit, testQueue, testWait)(blockingHandler(entered, release))
		fillSlots(t, handler, entered)

		select {
		case recorder := <-post(handler):
			assert.Equal(t, http.StatusServiceUnavailable, recorder.Code)
			assert.Equal(t, "1", recorder.Header().Get("Retry-After"))
		case <-time.After(2 * time.Second):
			t.Fatal("a request past the limit was neither served nor refused")
		}
	})

	t.Run("holds a request until a slot frees instead of failing it", func(t *testing.T) {
		entered := make(chan struct{}, testLimit+1)
		release := make(chan struct{})
		// Freed by hand once the queued request is confirmed waiting, and by
		// the defer if an assertion gives up before that.
		freeSlots := sync.OnceFunc(func() { close(release) })
		defer freeSlots()

		handler := Throttle(testLimit, testQueue, time.Minute)(blockingHandler(entered, release))
		fillSlots(t, handler, entered)

		answered := post(handler)

		// With every slot held, the extra request has nowhere to run: it must
		// be queued rather than served or refused.
		select {
		case <-entered:
			t.Fatal("a request ran while every slot was busy")
		case recorder := <-answered:
			t.Fatalf("a queued request was answered with %d", recorder.Code)
		case <-time.After(100 * time.Millisecond):
		}

		freeSlots()

		select {
		case recorder := <-answered:
			assert.Equal(t, http.StatusOK, recorder.Code)
		case <-time.After(time.Second):
			t.Fatal("a queued request never got its slot")
		}
	})
}
