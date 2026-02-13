package namegen

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func generateAllVariants() map[string]bool {
	variants := make(map[string]bool)

	for v := range 2 {
		for _, a := range adjectives[v] {
			for _, n := range nouns[v] {
				variants[fmt.Sprintf("%s %s", a, n)] = true
			}
		}
	}
	return variants
}

func TestGenerateName(t *testing.T) {
	allVariants := generateAllVariants()

	t.Run("should generate accurate names", func(t *testing.T) {
		for range 100 {
			name := GenerateName()
			require.True(t, allVariants[name], "name %s is not in the list of all variants", name)
		}
	})
}
