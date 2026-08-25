package httpx

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testSchema struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (s testSchema) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Name, validation.Required),
		validation.Field(&s.Email, validation.Required, is.Email),
	)
}

func TestParseBody(t *testing.T) {
	l := zerolog.Nop()

	t.Run("parse valid json", func(t *testing.T) {
		request := httptest.NewRequest("POST", "/", bytes.NewBufferString(`{"name": "John", "email": "john@example.com"}`))
		response := httptest.NewRecorder()

		data, ok := ParseBody[testSchema](response, request, l)
		require.True(t, ok)
		require.Equal(t, "John", data.Name)
		require.Equal(t, "john@example.com", data.Email)
	})

	t.Run("parse valid json", func(t *testing.T) {
		request := httptest.NewRequest("POST", "/", bytes.NewBufferString(`{"name": "John"}`))
		response := httptest.NewRecorder()

		data, ok := ParseBody[testSchema](response, request, l)
		require.True(t, ok)
		require.Equal(t, "John", data.Name)
		require.Equal(t, "", data.Email)
	})

	t.Run("invalid json", func(t *testing.T) {
		request := httptest.NewRequest("POST", "/", bytes.NewBufferString(`{ name: "John"}`))
		response := httptest.NewRecorder()

		_, ok := ParseBody[testSchema](response, request, l)
		require.False(t, ok)
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})
}

func TestParseAndValidateBody(t *testing.T) {
	l := zerolog.Nop()

	t.Run("parse valid json", func(t *testing.T) {
		request := httptest.NewRequest("POST", "/", bytes.NewBufferString(`{"name": "John", "email": "john@example.com"}`))
		response := httptest.NewRecorder()

		data, ok := ParseAndValidateBody[testSchema](response, request, l)
		require.True(t, ok)
		require.Equal(t, "John", data.Name)
		require.Equal(t, "john@example.com", data.Email)
	})

	t.Run("parse valid json with incorrect schema", func(t *testing.T) {
		request := httptest.NewRequest("POST", "/", bytes.NewBufferString(`{"name": "John"}`))
		response := httptest.NewRecorder()

		_, ok := ParseAndValidateBody[testSchema](response, request, l)
		require.False(t, ok)
		assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
	})

	t.Run("invalid json", func(t *testing.T) {
		request := httptest.NewRequest("POST", "/", bytes.NewBufferString(`{ name: "John"}`))
		response := httptest.NewRecorder()

		_, ok := ParseBody[testSchema](response, request, l)
		require.False(t, ok)
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})
}

func TestQueryInt(t *testing.T) {
	tcs := map[string]int{
		"":         7, // отсутствует
		"?n=3":     3,
		"?n=-3":    -3, // границы — дело вызывающего
		"?n=abc":   7,
		"?n=":      7,
		"?n=3.5":   7,
		"?n=99999": 99999,
	}

	for query, expected := range tcs {
		r := httptest.NewRequest(http.MethodGet, "/"+query, nil)
		assert.Equal(t, expected, QueryInt(r, "n", 7), query)
	}
}
