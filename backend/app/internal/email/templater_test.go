package email

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTemplater_AuthLink(t *testing.T) {
	t.Run("should contain confirmation link with token", func(t *testing.T) {
		templater := NewTemplate("https://neskuchka.ru")

		text, err := templater.AuthLink("test-token-123")

		require.NoError(t, err)
		assert.Contains(t, text, "https://neskuchka.ru/login/confirm?token=test-token-123")
	})

	t.Run("should work with localhost base URL", func(t *testing.T) {
		templater := NewTemplate("http://localhost:5173")

		text, err := templater.AuthLink("abc")

		require.NoError(t, err)
		assert.Contains(t, text, "http://localhost:5173/login/confirm?token=abc")
	})

	t.Run("should not have trailing slash duplication", func(t *testing.T) {
		templater := NewTemplate("https://neskuchka.ru/")

		text, err := templater.AuthLink("token")

		require.NoError(t, err)
		assert.NotContains(t, text, "ru//login")
	})
}
