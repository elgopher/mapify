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

		t.Run("should map pointer to primitive", func(t *testing.T) {
			str := "str"
			number := 3
			expected := []interface{}{&str, &number}

			for _, val := range expected {
				result := instance.MapAny(val)
				assert.Equal(t, val, result)
			}
		})

		t.Run("should map nil", func(t *testing.T) {
			actual := instance.MapAny(nil)
			assert.Nil(t, actual)
		})

		t.Run("should map an empty struct", func(t *testing.T) {
			actual := instance.MapAny(struct{}{})
			assert.Equal(t, map[string]interface{}{}, actual)
		})

		t.Run("should map a pointer to empty struct", func(t *testing.T) {
			s := struct{}{}
			actual := instance.MapAny(&s)
			assert.Equal(t, map[string]interface{}{}, actual)
		})

		t.Run("should map a pointer to nil struct", func(t *testing.T) {
			var s *struct{}
			actual := instance.MapAny(s)
			assert.Nil(t, actual)
		})

		t.Run("should map a pointer to pointer to nil struct", func(t *testing.T) {
			var s **struct{}
			actual := instance.MapAny(s)
			assert.Nil(t, actual)
		})

		t.Run("should map a zero-value struct with one field", func(t *testing.T) {
			s := struct {
				Field string
			}{}
			actual := instance.MapAny(s)
			assert.Equal(t,
				map[string]interface{}{
					"Field": "",
				},
				actual)
		})

		t.Run("should map an empty slice of strings", func(t *testing.T) {
			actual := instance.MapAny([]string{})
			assert.Equal(t, []string{}, actual)
		})
	})
}
