package httpx

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	R "github.com/go-pkgz/rest"
	"github.com/rs/zerolog"
)

// ParseBody parses the body to the given type
// Returns the data and a boolean indicating if the parsing was successful
func ParseBody[T any](w http.ResponseWriter, r *http.Request, l zerolog.Logger) (T, bool) {
	var data T
	err := R.DecodeJSON(r, &data)
	if err != nil {
		RenderError(w, l, http.StatusBadRequest, err, "Failed to decode body")
		return data, false
	}

	return data, true
}

// ParseAndValidateBody parses the body to the given type and runs validation
// Returns the data and a boolean indicating if the parsing and validation were successful
func ParseAndValidateBody[T validation.Validatable](w http.ResponseWriter, r *http.Request, l zerolog.Logger) (T, bool) {
	data, ok := ParseBody[T](w, r, l)
	if !ok {
		return data, false
	}

	err := data.Validate()
	if err != nil {
		RenderValidationError(w, l, err)
		return data, false
	}

	return data, true
}
