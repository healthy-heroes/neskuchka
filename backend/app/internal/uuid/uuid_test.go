package uuid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("should generate a new UUID", func(t *testing.T) {
		uuid := New()
		assert.NotEmpty(t, uuid)
	})
}
