package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func setupTestHandler(expiration time.Duration, version string) http.Handler {
	// Setup server
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	})

	middleware := CacheControl(expiration, version)
	handler := middleware(testHandler)

	return handler
}

func TestCacheControl(t *testing.T) {
	tests := []struct {
		name          string
		expiration    time.Duration
		version       string
		url           string
		expectedCache string
	}{
		{
			name:          "Basic cache control with 1 hour expiration",
			expiration:    time.Hour,
			version:       "v1.0.0",
			url:           "/api/test",
			expectedCache: "max-age=3600, no-cache",
		},
		{
			name:          "Cache control with 30 minutes expiration",
			expiration:    time.Minute * 30,
			version:       "v2.0.0",
			url:           "/api/data",
			expectedCache: "max-age=1800, no-cache",
		},
		{
			name:          "Cache control with zero expiration",
			expiration:    0,
			version:       "v3.0.0",
			url:           "/api/zero",
			expectedCache: "max-age=0, no-cache",
		},
		{
			name:          "Cache control with negative expiration",
			expiration:    -time.Hour,
			version:       "v4.0.0",
			url:           "/api/negative",
			expectedCache: "max-age=-3600, no-cache",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := setupTestHandler(tt.expiration, tt.version)

			req := httptest.NewRequest("GET", tt.url, nil)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, http.StatusOK, rr.Code)
			assert.NotEmpty(t, rr.Header().Get("Etag"))
			assert.Equal(t, tt.expectedCache, rr.Header().Get("Cache-Control"))
		})
	}
}

func TestCacheControl_IfNoneMatch_Ok(t *testing.T) {
	handler := setupTestHandler(time.Hour, "v1.0.0")

	req := httptest.NewRequest("GET", "/api/if-none-match-ok", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	req.Header.Set("If-None-Match", rr.Header().Get("Etag"))
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotModified, rr.Code)
}

func TestCacheControl_IfNoneMatch_Wrong(t *testing.T) {
	handler := setupTestHandler(time.Hour, "v1.0.0")

	req := httptest.NewRequest("GET", "/api/if-none-match-wrong", nil)
	req.Header.Set("If-None-Match", "wrong-etag")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestCacheControl_ETagConsistency(t *testing.T) {
	handler := setupTestHandler(time.Hour, "v1.0.0")

	etags := make(map[string]bool)
	for range 4 {
		req := httptest.NewRequest("GET", "/api/etag-consistency", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		etag := rr.Header().Get("Etag")
		etags[etag] = true
	}

	// All ETags should be the same
	assert.Equal(t, 1, len(etags), "All ETags should be identical for the same URL and version")
}

func TestCacheControl_DifferentURLs(t *testing.T) {
	handler := setupTestHandler(time.Hour, "v1.0.0")

	urls := []string{
		"/api/test1",
		"/api/test2",
		"/api/test3",
		"/api/test1?param=value",
		"/api/test1?param=value&other=123",
	}

	etags := make(map[string]bool)
	for _, url := range urls {
		req := httptest.NewRequest("GET", url, nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		etag := rr.Header().Get("Etag")
		etags[etag] = true
	}

	assert.Equal(t, len(urls), len(etags), "Different URLs should produce different ETags")
}

func TestCacheControl_DifferentETagsForDifferentVersions(t *testing.T) {
	expiration := time.Hour
	url := "/api/test"

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	versions := []string{"v1.0.0", "v1.0.1", "v2.0.0", "v2.1.0"}
	etags := make(map[string]bool)

	for _, version := range versions {
		middleware := CacheControl(expiration, version)
		handler := middleware(testHandler)

		req := httptest.NewRequest("GET", url, nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		etag := rr.Header().Get("Etag")
		etags[etag] = true
	}

	assert.Equal(t, len(versions), len(etags), "Different versions should produce different ETags")
}
