// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

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
				assert.Same(t, val, result)
			}
		})

		t.Run("should map nil", func(t *testing.T) {
			actual := instance.MapAny(nil)
			assert.Nil(t, actual)
		})

		t.Run("should map pointer to nil primitive", func(t *testing.T) {
			var str *string
			actual := instance.MapAny(str)
			assert.Same(t, str, actual)
		})

		t.Run("should map an empty struct", func(t *testing.T) {
			actual := instance.MapAny(struct{}{})
			assert.IsType(t, map[string]interface{}{}, actual)
			assert.Empty(t, actual)
		})

		t.Run("should map a pointer to empty struct", func(t *testing.T) {
			s := struct{}{}
			actual := instance.MapAny(&s)
			assert.IsType(t, map[string]interface{}{}, actual)
			assert.Empty(t, actual)
		})

		t.Run("should map a pointer to nil struct", func(t *testing.T) {
			var s *struct{}
			actual := instance.MapAny(s)
			assert.Same(t, s, actual)
		})

		t.Run("should map a zero-value struct with two fields", func(t *testing.T) {
			s := struct {
				Field1 string
				Field2 string
			}{}
			actual := instance.MapAny(s)
			assert.Equal(t,
				map[string]interface{}{
					"Field1": "",
					"Field2": "",
				},
				actual)
		})

		t.Run("should map a struct with only private fields", func(t *testing.T) {
			s := struct {
				field1 string
				field2 string
			}{}
			actual := instance.MapAny(s)
			assert.IsType(t, map[string]interface{}{}, actual)
			assert.Empty(t, actual)
		})

		t.Run("should map a struct with field specified", func(t *testing.T) {
			s := struct {
				Field string
			}{
				Field: "value",
			}
			actual := instance.MapAny(s)
			assert.Equal(t,
				map[string]interface{}{
					"Field": s.Field,
				},
				actual)
		})

		t.Run("should map a struct with field pointer specified", func(t *testing.T) {
			str := "value"
			s := struct {
				Field *string
			}{
				Field: &str,
			}
			// when
			actual := instance.MapAny(s)
			// then
			assert.Equal(t,
				map[string]interface{}{
					"Field": s.Field,
				},
				actual)
		})

		t.Run("should map a struct with nil field", func(t *testing.T) {
			s := struct {
				Field *string
			}{}
			actual := instance.MapAny(s)
			assert.Equal(t,
				map[string]interface{}{
					"Field": s.Field,
				},
				actual)
		})

		t.Run("should map a struct with nested struct", func(t *testing.T) {
			type nestedStruct struct {
				Field string
			}

			s := struct {
				Nested nestedStruct
			}{
				Nested: nestedStruct{
					Field: "value",
				},
			}
			actual := instance.MapAny(s)
			assert.Equal(t,
				map[string]interface{}{
					"Nested": map[string]interface{}{
						"Field": s.Nested.Field,
					},
				},
				actual)
		})

		t.Run("should map a struct with nested nil struct", func(t *testing.T) {
			s := struct {
				Nested *struct{}
			}{}
			actual := instance.MapAny(s)
			assert.Equal(t,
				map[string]interface{}{
					"Nested": s.Nested,
				},
				actual)
		})

		t.Run("should map an empty slice of strings", func(t *testing.T) {
			actual := instance.MapAny([]string{})
			assert.Equal(t, []string{}, actual)
		})

		t.Run("should map an nil slice of strings", func(t *testing.T) {
			var given []string
			actual := instance.MapAny(given)
			assert.Equal(t, given, actual)
		})

		t.Run("should map an slice of two strings", func(t *testing.T) {
			given := []string{"1", "2"}
			actual := instance.MapAny(given)
			assert.Equal(t, given, actual)
		})

		t.Run("should map an slice of pointer to string", func(t *testing.T) {
			str1 := "1"
			given := []*string{&str1}
			actual := instance.MapAny(given)
			assert.Equal(t, given, actual)
		})

		t.Run("should map a slice of empty structs", func(t *testing.T) {
			s := []struct{}{
				{},
				{},
			}
			actual := instance.MapAny(s)
			assert.Equal(t,
				[]map[string]interface{}{
					{},
					{},
				},
				actual)
		})

		t.Run("should map a slice of structs", func(t *testing.T) {
			type structWithField struct {
				Field string
			}
			s := []structWithField{
				{Field: "value1"},
				{Field: "value2"},
			}
			actual := instance.MapAny(s)
			assert.Equal(t,
				[]map[string]interface{}{
					{
						"Field": s[0].Field,
					},
					{
						"Field": s[1].Field,
					},
				},
				actual)
		})

		t.Run("should map slice of slices of structs", func(t *testing.T) {
			type structWithField struct {
				Field string
			}
			s := [][]structWithField{
				{{Field: "A1"}, {Field: "A2"}},
				{{Field: "B1"}, {Field: "B2"}},
			}
			actual := instance.MapAny(s)
			assert.Equal(t,
				[][]map[string]interface{}{
					{
						map[string]interface{}{"Field": s[0][0].Field},
						map[string]interface{}{"Field": s[0][1].Field},
					},
					{
						map[string]interface{}{"Field": s[1][0].Field},
						map[string]interface{}{"Field": s[1][1].Field},
					},
				},
				actual)
		})

		t.Run("should map a struct with nested slice of structs", func(t *testing.T) {
			type nestedStruct struct {
				Field string
			}

			s := struct {
				Nested []nestedStruct
			}{
				Nested: []nestedStruct{
					{Field: "1"},
					{Field: "2"},
				},
			}
			actual := instance.MapAny(s)
			assert.Equal(t,
				map[string]interface{}{
					"Nested": []map[string]interface{}{
						{"Field": s.Nested[0].Field},
						{"Field": s.Nested[1].Field},
					},
				},
				actual)
		})

	})
}

func TestFilter(t *testing.T) {
	t.Run("should filter out all struct fields", func(t *testing.T) {
		s := struct {
			A, B string
		}{}
		instance := mapify.Instance{
			Filter: func(path string, e mapify.Element) bool {
				return false
			},
		}
		// when
		v := instance.MapAny(s)
		// then
		assert.Empty(t, v)
	})

	t.Run("should filter by struct field path", func(t *testing.T) {
		s := struct {
			A, B string
		}{}
		instance := mapify.Instance{
			Filter: func(path string, e mapify.Element) bool {
				return path == ".A"
			},
		}
		// when
		v := instance.MapAny(s)
		// then
		assert.Equal(t, map[string]interface{}{
			"A": "",
		}, v)
	})

	t.Run("should filter by nested struct field path", func(t *testing.T) {
		s := struct {
			Nested struct {
				A string
			}
		}{}
		instance := mapify.Instance{
			Filter: func(path string, e mapify.Element) bool {
				return path == ".Nested" || path == ".Nested.A"
			},
		}
		// when
		v := instance.MapAny(s)
		// then
		assert.Equal(t, map[string]interface{}{
			"Nested": map[string]interface{}{"A": ""},
		}, v)
	})

	t.Run("should filter by slice element path", func(t *testing.T) {
		s := []struct{ Field string }{
			{Field: "0"},
			{Field: "1"},
		}
		instance := mapify.Instance{
			Filter: func(path string, e mapify.Element) bool {
				return path == "[1].Field"
			},
		}
		// when
		v := instance.MapAny(s)
		// then
		assert.Equal(t,
			[]map[string]interface{}{
				{},
				{"Field": s[1].Field},
			},
			v)
	})

	t.Run("should filter by 2d slice element path", func(t *testing.T) {
		s := [][]struct{ Field string }{
			{
				{Field: "A0"},
			},
			{
				{Field: "B0"},
				{Field: "B1"},
			},
		}
		instance := mapify.Instance{
			Filter: func(path string, e mapify.Element) bool {
				return path == "[1][1].Field"
			},
		}
		// when
		v := instance.MapAny(s)
		// then
		assert.Equal(t,
			[][]map[string]interface{}{
				{
					{},
				},
				{
					{},
					{"Field": s[1][1].Field},
				},
			},
			v)
	})

	t.Run("should filter by field name", func(t *testing.T) {
		instance := mapify.Instance{
			Filter: func(path string, e mapify.Element) bool {
				return e.Name() == "Field"
			},
		}
		// when
		v := instance.MapAny(
			struct{ Field string }{
				Field: "v",
			},
		)
		// then
		assert.Equal(t,
			map[string]interface{}{
				"Field": "v",
			},
			v)
	})
}

func TestRename(t *testing.T) {
	t.Run("should rename struct field", func(t *testing.T) {
		instance := mapify.Instance{
			Rename: func(path string, e mapify.Element) string {
				return "newName"
			},
		}
		// when
		v := instance.MapAny(
			struct{ OldName string }{
				OldName: "v",
			},
		)
		// then
		assert.Equal(t,
			map[string]interface{}{
				"newName": "v",
			},
			v)
	})
}
