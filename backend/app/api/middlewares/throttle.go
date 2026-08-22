package middlewares

import (
	"net/http"
	"time"

	chiMW "github.com/go-chi/chi/v5/middleware"
)

// Throttle bounds how many requests run at once. Past limit they wait, past
// limit+queue they are turned away at once, as is anyone still waiting after
// wait. A refusal is 503 rather than chi's default 429: the ceiling belongs to
// the server and says nothing about the client that happened to hit it.
func Throttle(limit, queue int, wait time.Duration) func(http.Handler) http.Handler {
	return chiMW.ThrottleWithOpts(chiMW.ThrottleOpts{
		Limit:          limit,
		BacklogLimit:   queue,
		BacklogTimeout: wait,
		StatusCode:     http.StatusServiceUnavailable,
		RetryAfterFn:   func(_ bool) time.Duration { return wait },
	})
}
