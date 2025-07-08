package middlewares

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"time"
)

// CacheControl is a middleware setting cache expiration. Using url+version as etag
func CacheControl(expiration time.Duration, version string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			etag := makeEtag(r, version)

			w.Header().Set("Etag", etag)
			w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d, no-cache", int(expiration.Seconds())))

			if match := r.Header.Get("If-None-Match"); match != "" {
				if match == etag {
					w.WriteHeader(http.StatusNotModified)
					return
				}
			}
			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func makeEtag(r *http.Request, version string) string {
	data := version + ":" + r.URL.String()

	hash := sha256.Sum256([]byte(data))

	return fmt.Sprintf("\"%x\"", hash)
}
