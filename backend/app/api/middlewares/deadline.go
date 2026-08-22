package middlewares

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

// ReadDeadline gives a request d to arrive in full, body included. Routes that
// read a body need it: a client trickling its bytes otherwise sits in the
// handler for as long as it pleases, holding everything the handler holds — a
// Throttle slot, most of all. chi's Timeout is no substitute, as it cancels the
// context and a blocked Read never consults one.
//
// A writer with no connection behind it can't take a deadline. That is worth a
// line in the log rather than a failed upload: the request is still perfectly
// serviceable, it just isn't bounded.
func ReadDeadline(logger zerolog.Logger, d time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := http.NewResponseController(w).SetReadDeadline(time.Now().Add(d)); err != nil {
				logger.Warn().Err(err).Msg("Failed to bound the request read")
			}

			next.ServeHTTP(w, r)
		})
	}
}
