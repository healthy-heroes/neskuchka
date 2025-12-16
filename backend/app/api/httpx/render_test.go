package httpx

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func body(output string) string {
	return output + "\n"
}

func TestRender(t *testing.T) {
	t.Run("render data", func(t *testing.T) {
		var data = map[string]string{
			"name":  "John",
			"email": "john@example.com",
		}

		response := httptest.NewRecorder()
		Render(response, data)

		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, body(`{"data":{"email":"john@example.com","name":"John"}}`), response.Body.String())
	})

	t.Run("render invalid data", func(t *testing.T) {
		response := httptest.NewRecorder()
		Render(response, func() {})

		assert.Equal(t, "text/plain; charset=utf-8", response.Header().Get("Content-Type"))
		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})
}

func TestRenderError(t *testing.T) {
	t.Run("render error", func(t *testing.T) {
		response := httptest.NewRecorder()
		RenderError(response, zerolog.Nop(), http.StatusInternalServerError, errors.New("error"), "caused by error")

		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.Equal(t, body(`{"error":"caused by error"}`), response.Body.String())
	})
}

func TestRenderValidationError(t *testing.T) {
	t.Run("render validation error", func(t *testing.T) {
		validationErrors := validation.Errors{
			"Name":  errors.New("is bad"),
			"Email": errors.New("is required"),
		}

		response := httptest.NewRecorder()
		RenderValidationError(response, zerolog.Nop(), validationErrors)

		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
		assert.Equal(t, http.StatusUnprocessableEntity, response.Code)

		expected := body(`{"error":"Failed to validate data","details":{"Email":"is required","Name":"is bad"}}`)
		assert.Equal(t, expected, response.Body.String())
	})
}

func TestRenderUnauthorized(t *testing.T) {
	t.Run("returns 401 with JSON error", func(t *testing.T) {
		response := httptest.NewRecorder()

		RenderUnauthorized(response)

		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
		assert.Equal(t, http.StatusUnauthorized, response.Code)
		assert.Equal(t, body(`{"error":"Unauthorized"}`), response.Body.String())
	})
}
