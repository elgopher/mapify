package mapify_test

import (
	"testing"

	"github.com/elgopher/mapify"
	"github.com/stretchr/testify/assert"
)

func TestInstance_MapAny(t *testing.T) {
	t.Run("for default Instance", func(t *testing.T) {
		instance := mapify.Instance{}

		t.Run("should map primitive", func(t *testing.T) {
			expected := []interface{}{1, 1.0, "str"}

			for _, val := range expected {
				result := instance.MapAny(val)
				assert.Equal(t, val, result)
			}
		})
	})
}
