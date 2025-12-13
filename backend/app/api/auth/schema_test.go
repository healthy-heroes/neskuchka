package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoginSchemaValidation(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		schema := LoginSchema{Email: "test@example.com"}

		err := schema.Validate()
		assert.NoError(t, err)
	})

	t.Run("empty schema", func(t *testing.T) {
		schema := LoginSchema{}

		err := schema.Validate()
		assert.Error(t, err)
	})

	t.Run("invalid email", func(t *testing.T) {
		schema := LoginSchema{Email: "invalid-email"}

		err := schema.Validate()
		assert.Error(t, err)
	})
}
