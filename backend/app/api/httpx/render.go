package httpx

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rs/zerolog"
)

// Response represents a successful response with data
type Response struct {
	Data interface{} `json:"data"`
}

// ErrorResponse represents an error response with details
type ErrorResponse struct {
	Error string `json:"error"`
}

// ValidationErrorResponse represents a validation error response with details
type ValidationErrorResponse struct {
	Error   string            `json:"error"`
	Details validation.Errors `json:"details"`
}

// Render sends JSON response with data and status code 200
func Render(w http.ResponseWriter, data interface{}) {
	response := Response{
		Data: data,
	}
	renderJSONWithStatus(w, response, http.StatusOK)
}

// RenderError sends JSON response with error and status code
func RenderError(w http.ResponseWriter, l zerolog.Logger, code int, err error, msg string) {
	l.Error().Err(err).Msg(msg)

	response := ErrorResponse{
		Error: msg,
	}
	renderJSONWithStatus(w, response, code)
}

// RenderValidationError sends JSON response with validation error and status code 422
func RenderValidationError(w http.ResponseWriter, l zerolog.Logger, err error) {
	l.Error().Msgf("Failed to validate data: %s", err)

	var validationResults validation.Errors
	var validationError validation.InternalError
	if errors.As(err, &validationError) || !errors.As(err, &validationResults) {
		RenderError(w, l, http.StatusBadRequest, validationError, "Failed to run validation")
		return
	}

	response := ValidationErrorResponse{
		Error:   "Failed to validate data",
		Details: validationResults,
	}
	renderJSONWithStatus(w, response, http.StatusUnprocessableEntity)
}

func RenderUnauthorized(w http.ResponseWriter) {
	renderJSONWithStatus(w, ErrorResponse{
		Error: "Unauthorized",
	}, http.StatusUnauthorized)
}

func renderJSONWithStatus(w http.ResponseWriter, data interface{}, code int) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)
	if err := enc.Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_, _ = w.Write(buf.Bytes())
}
