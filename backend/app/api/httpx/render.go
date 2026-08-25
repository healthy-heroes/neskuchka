package httpx

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rs/zerolog"

	"github.com/healthy-heroes/neskuchka/backend/app/domain"
)

// Response represents a successful response with data
type Response struct {
	Data any `json:"data"`
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
func Render(w http.ResponseWriter, data any) {
	response := Response{
		Data: data,
	}
	renderJSONWithStatus(w, response, http.StatusOK)
}

// RenderError sends JSON response with error and status code.
//
// Anything below 500 is something the caller did — a closed edit window, a
// stale id — and is logged as a warning. Error level is for the server's own
// failures, so that a log filtered to errors stays worth reading.
func RenderError(w http.ResponseWriter, l zerolog.Logger, code int, err error, msg string) {
	level := zerolog.ErrorLevel
	if code < http.StatusInternalServerError {
		level = zerolog.WarnLevel
	}
	l.WithLevel(level).Err(err).Msg(msg)

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

func renderJSONWithStatus(w http.ResponseWriter, data any, code int) {
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

func RenderDomainError(w http.ResponseWriter, l zerolog.Logger, err error, msg string) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		RenderError(w, l, http.StatusNotFound, err, "Not found")
	case errors.Is(err, domain.ErrForbidden):
		RenderError(w, l, http.StatusForbidden, err, "Forbidden")
	case errors.Is(err, domain.ErrLocked):
		RenderError(w, l, http.StatusConflict, err, "Conflict")
	default:
		RenderError(w, l, http.StatusInternalServerError, err, msg)
	}
}
