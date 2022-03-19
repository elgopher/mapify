// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package mapify_test

import (
	"testing"

	"github.com/elgopher/mapify"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMapper_MapAny(t *testing.T) {
	t.Run("for default Mapper", func(t *testing.T) {
		mapper := mapify.Mapper{}

		t.Run("should map primitive", func(t *testing.T) {
			expected := []interface{}{1, 1.0, "str"}

			for _, val := range expected {
				result, err := mapper.MapAny(val)
				require.NoError(t, err)
				assert.Equal(t, val, result)
			}
		})

		t.Run("should map pointer to primitive", func(t *testing.T) {
			str := "str"
			number := 3
			expected := []interface{}{&str, &number}

			for _, val := range expected {
				result, err := mapper.MapAny(val)
				require.NoError(t, err)
				assert.Same(t, val, result)
			}
		})

		t.Run("should map nil", func(t *testing.T) {
			actual, err := mapper.MapAny(nil)
			require.NoError(t, err)
			assert.Nil(t, actual)
		})

		t.Run("should map pointer to nil primitive", func(t *testing.T) {
			var str *string
			actual, err := mapper.MapAny(str)
			require.NoError(t, err)
			assert.Same(t, str, actual)
		})

		t.Run("should map an empty struct", func(t *testing.T) {
			actual, err := mapper.MapAny(struct{}{})
			require.NoError(t, err)
			assert.IsType(t, map[string]interface{}{}, actual)
			assert.Empty(t, actual)
		})

		t.Run("should map a pointer to empty struct", func(t *testing.T) {
			s := struct{}{}
			actual, err := mapper.MapAny(&s)
			require.NoError(t, err)
			assert.IsType(t, map[string]interface{}{}, actual)
			assert.Empty(t, actual)
		})

		t.Run("should map a pointer to nil struct", func(t *testing.T) {
			var s *struct{}
			actual, err := mapper.MapAny(s)
			require.NoError(t, err)
			assert.Same(t, s, actual)
		})

		t.Run("should map a zero-value struct with two fields", func(t *testing.T) {
			s := struct {
				Field1 string
				Field2 string
			}{}
			actual, err := mapper.MapAny(s)
			require.NoError(t, err)
			expected := map[string]interface{}{
				"Field1": "",
				"Field2": "",
			}
			assert.Equal(t, expected, actual)
		})

		t.Run("should map a struct with only private fields", func(t *testing.T) {
			s := struct {
				field1 string
				field2 string
			}{}
			actual, err := mapper.MapAny(s)
			require.NoError(t, err)
			assert.IsType(t, map[string]interface{}{}, actual)
			assert.Empty(t, actual)
		})

		t.Run("should map a struct with field specified", func(t *testing.T) {
			s := struct{ Field string }{Field: "value"}
			actual, err := mapper.MapAny(s)
			require.NoError(t, err)
			expected := map[string]interface{}{
				"Field": s.Field,
			}
			assert.Equal(t, expected, actual)
		})

		t.Run("should map a struct with field pointer specified", func(t *testing.T) {
			str := "value"
			s := struct{ Field *string }{Field: &str}
			// when
			actual, err := mapper.MapAny(s)
			// then
			require.NoError(t, err)
			expected := map[string]interface{}{
				"Field": s.Field,
			}
			assert.Equal(t, expected, actual)
		})

		t.Run("should map a struct with nil field", func(t *testing.T) {
			s := struct{ Field *string }{}
			actual, err := mapper.MapAny(s)
			require.NoError(t, err)
			expected := map[string]interface{}{
				"Field": s.Field,
			}
			assert.Equal(t, expected, actual)
		})

		t.Run("should map a struct with nested struct", func(t *testing.T) {
			type nestedStruct struct{ Field string }
			s := struct{ Nested nestedStruct }{
				Nested: nestedStruct{Field: "value"},
			}
			actual, err := mapper.MapAny(s)
			require.NoError(t, err)
			expected := map[string]interface{}{
				"Nested": map[string]interface{}{
					"Field": s.Nested.Field,
				},
			}
			assert.Equal(t, expected, actual)
		})

		t.Run("should map a struct with nested nil struct", func(t *testing.T) {
			s := struct{ Nested *struct{} }{}
			actual, err := mapper.MapAny(s)
			require.NoError(t, err)
			expected := map[string]interface{}{
				"Nested": s.Nested,
			}
			assert.Equal(t, expected, actual)
		})

		t.Run("should map an empty slice of strings", func(t *testing.T) {
			actual, err := mapper.MapAny([]string{})
			require.NoError(t, err)
			assert.Equal(t, []string{}, actual)
		})

		t.Run("should map an nil slice of strings", func(t *testing.T) {
			var given []string
			actual, err := mapper.MapAny(given)
			require.NoError(t, err)
			assert.Equal(t, given, actual)
		})

		t.Run("should map an slice of two strings", func(t *testing.T) {
			given := []string{"1", "2"}
			actual, err := mapper.MapAny(given)
			require.NoError(t, err)
			assert.Equal(t, given, actual)
		})

		t.Run("should map an slice of pointer to string", func(t *testing.T) {
			str1 := "1"
			given := []*string{&str1}
			actual, err := mapper.MapAny(given)
			require.NoError(t, err)
			assert.Equal(t, given, actual)
		})

		t.Run("should map a slice of empty structs", func(t *testing.T) {
			s := []struct{}{
				{},
				{},
			}
			actual, err := mapper.MapAny(s)
			require.NoError(t, err)
			expected := []map[string]interface{}{
				{},
				{},
			}
			assert.Equal(t, expected, actual)
		})

		t.Run("should map a slice of structs", func(t *testing.T) {
			type structWithField struct{ Field string }
			s := []structWithField{
				{Field: "value1"},
				{Field: "value2"},
			}
			actual, err := mapper.MapAny(s)
			require.NoError(t, err)
			expected := []map[string]interface{}{
				{
					"Field": s[0].Field,
				},
				{
					"Field": s[1].Field,
				},
			}
			assert.Equal(t, expected, actual)
		})

		t.Run("should map slice of slices of structs", func(t *testing.T) {
			type structWithField struct{ Field string }
			s := [][]structWithField{
				{{Field: "A1"}, {Field: "A2"}},
				{{Field: "B1"}, {Field: "B2"}},
			}
			actual, err := mapper.MapAny(s)
			require.NoError(t, err)
			expected := [][]map[string]interface{}{
				{
					map[string]interface{}{"Field": s[0][0].Field},
					map[string]interface{}{"Field": s[0][1].Field},
				},
				{
					map[string]interface{}{"Field": s[1][0].Field},
					map[string]interface{}{"Field": s[1][1].Field},
				},
			}
			assert.Equal(t, expected, actual)
		})

		t.Run("should map a struct with nested slice of structs", func(t *testing.T) {
			type nestedStruct struct{ Field string }
			s := struct {
				Nested []nestedStruct
			}{
				Nested: []nestedStruct{
					{Field: "1"},
					{Field: "2"},
				},
			}
			actual, err := mapper.MapAny(s)
			require.NoError(t, err)
			expected := map[string]interface{}{
				"Nested": []map[string]interface{}{
					{"Field": s.Nested[0].Field},
					{"Field": s.Nested[1].Field},
				},
			}
			assert.Equal(t, expected, actual)
		})
	})

	t.Run("should convert map[string]any to map[string]interface{}", func(t *testing.T) {
		m := map[string]string{"key1": "value1", "key2": "value2"}
		mapper := mapify.Mapper{}
		result, err := mapper.MapAny(m)
		require.NoError(t, err)
		assert.Equal(t, map[string]interface{}{"key1": "value1", "key2": "value2"}, result)
	})

	t.Run("should not convert map which does not have a string key", func(t *testing.T) {
		tests := map[string]interface{}{
			"map[struct]":     map[struct{ A string }]string{},
			"[]map[struct]":   []map[struct{ A string }]string{{}},
			"[][]map[struct]": [][]map[struct{ A string }]string{{}},
		}

		for name, theMap := range tests {
			t.Run(name, func(t *testing.T) {
				mapper := mapify.Mapper{}
				result, err := mapper.MapAny(theMap)
				require.NoError(t, err)
				assert.Equal(t, theMap, result)
			})
		}
	})
}

func TestFilter(t *testing.T) {
	t.Run("should filter out all elements", func(t *testing.T) {
		tests := map[string]interface{}{
			"struct": struct{ A, B string }{},
			"map":    map[string]string{"A": "", "B": ""},
		}

		for name, object := range tests {
			t.Run(name, func(t *testing.T) {
				mapper := mapify.Mapper{
					Filter: func(path string, e mapify.Element) (bool, error) {
						return false, nil
					},
				}
				// when
				v, err := mapper.MapAny(object)
				// then
				require.NoError(t, err)
				assert.Empty(t, v)
			})
		}
	})

	t.Run("should filter by element path", func(t *testing.T) {
		tests := map[string]interface{}{
			"struct": struct{ A, B string }{},
			"map":    map[string]string{"A": "", "B": ""},
		}

		for name, object := range tests {
			t.Run(name, func(t *testing.T) {
				mapper := mapify.Mapper{
					Filter: func(path string, e mapify.Element) (bool, error) {
						return path == ".A", nil
					},
				}
				// when
				v, err := mapper.MapAny(object)
				// then
				require.NoError(t, err)
				expected := map[string]interface{}{
					"A": "",
				}
				assert.Equal(t, expected, v)
			})
		}
	})

	t.Run("should filter by nested element path", func(t *testing.T) {
		tests := map[string]interface{}{
			"struct": struct {
				Nested struct{ A string }
			}{},
			"map": map[string]map[string]string{
				"Nested": {"A": ""},
			},
		}

		for name, object := range tests {
			t.Run(name, func(t *testing.T) {
				mapper := mapify.Mapper{
					Filter: func(path string, e mapify.Element) (bool, error) {
						return path == ".Nested" || path == ".Nested.A", nil
					},
				}
				// when
				v, err := mapper.MapAny(object)
				// then
				require.NoError(t, err)
				assert.Equal(t, map[string]interface{}{
					"Nested": map[string]interface{}{"A": ""},
				}, v)
			})
		}
	})

	t.Run("should filter by slice element path", func(t *testing.T) {
		tests := map[string]interface{}{
			"struct": []struct{ Field string }{
				{Field: "0"},
				{Field: "1"},
			},
			"map": []map[string]string{
				{"Field": "0"},
				{"Field": "1"},
			},
		}

		for name, object := range tests {
			t.Run(name, func(t *testing.T) {
				mapper := mapify.Mapper{
					Filter: func(path string, e mapify.Element) (bool, error) {
						return path == "[1].Field", nil
					},
				}
				// when
				v, err := mapper.MapAny(object)
				// then
				require.NoError(t, err)
				expected := []map[string]interface{}{
					{},
					{"Field": "1"},
				}
				assert.Equal(t, expected, v)
			})
		}
	})

	t.Run("should filter by 2d slice element path", func(t *testing.T) {
		tests := map[string]interface{}{
			"struct": [][]struct{ Field string }{
				{
					{Field: "A0"},
				},
				{
					{Field: "B0"},
					{Field: "B1"},
				},
			},
			"map": [][]map[string]string{
				{
					{"Field": "A0"},
				},
				{
					{"Field": "B0"},
					{"Field": "B1"},
				},
			},
		}

		for name, object := range tests {
			t.Run(name, func(t *testing.T) {
				mapper := mapify.Mapper{
					Filter: func(path string, e mapify.Element) (bool, error) {
						return path == "[1][1].Field", nil
					},
				}
				// when
				v, err := mapper.MapAny(object)
				// then
				require.NoError(t, err)
				expected := [][]map[string]interface{}{
					{
						{},
					},
					{
						{},
						{"Field": "B1"},
					},
				}
				assert.Equal(t, expected, v)
			})
		}
	})

	t.Run("should filter by element name", func(t *testing.T) {
		tests := map[string]interface{}{
			"struct": struct{ Field string }{
				Field: "v",
			},
			"map": map[string]string{
				"Field": "v",
			},
		}

		for name, object := range tests {
			t.Run(name, func(t *testing.T) {
				mapper := mapify.Mapper{
					Filter: func(path string, e mapify.Element) (bool, error) {
						return e.Name() == "Field", nil
					},
				}
				// when
				v, err := mapper.MapAny(object)
				// then
				require.NoError(t, err)
				expected := map[string]interface{}{
					"Field": "v",
				}
				assert.Equal(t, expected, v)
			})
		}
	})

	t.Run("should filter by value", func(t *testing.T) {
		tests := map[string]interface{}{
			"struct": struct{ Field1, Field2 string }{
				Field1: "keep it",
				Field2: "omit this",
			},
			"map": map[string]string{
				"Field1": "keep it",
				"Field2": "omit this",
			},
		}

		for name, object := range tests {
			t.Run(name, func(t *testing.T) {
				mapper := mapify.Mapper{
					Filter: func(path string, e mapify.Element) (bool, error) {
						return e.String() == "keep it", nil
					},
				}
				// when
				v, err := mapper.MapAny(object)
				// then
				require.NoError(t, err)
				expected := map[string]interface{}{
					"Field1": "keep it",
				}
				assert.Equal(t, expected, v)
			})
		}
	})

	t.Run("should return error when Filter returned error", func(t *testing.T) {
		tests := map[string]interface{}{
			"struct": struct{ Field string }{},
			"map":    map[string]string{"Field": ""},
		}

		for name, object := range tests {
			t.Run(name, func(t *testing.T) {
				givenError := stringError("err")
				mapper := mapify.Mapper{
					Filter: func(path string, e mapify.Element) (bool, error) {
						return false, givenError
					},
				}
				// when
				result, actualErr := mapper.MapAny(object)
				// then
				assert.Nil(t, result)
				assert.ErrorIs(t, actualErr, givenError)
			})
		}
	})
}

func TestRename(t *testing.T) {
	t.Run("should rename element", func(t *testing.T) {
		tests := map[string]interface{}{
			"struct": struct{ OldName string }{
				OldName: "v",
			},
			"map": map[string]string{
				"OldName": "v",
			},
		}
		for name, object := range tests {
			t.Run(name, func(t *testing.T) {
				mapper := mapify.Mapper{
					Rename: func(path string, e mapify.Element) (string, error) {
						return "newName", nil
					},
				}
				// when
				v, err := mapper.MapAny(object)
				// then
				require.NoError(t, err)
				expected := map[string]interface{}{
					"newName": "v",
				}
				assert.Equal(t, expected, v)
			})
		}
	})

	t.Run("should return error when Rename returned error", func(t *testing.T) {
		tests := map[string]interface{}{
			"struct": struct{ Field string }{},
			"map":    map[string]string{"Field": ""},
		}

		for name, object := range tests {
			t.Run(name, func(t *testing.T) {
				givenError := stringError("err")
				mapper := mapify.Mapper{
					Rename: func(path string, e mapify.Element) (string, error) {
						return e.Name(), givenError
					},
				}
				// when
				result, actualErr := mapper.MapAny(object)
				// then
				assert.Nil(t, result)
				assert.ErrorIs(t, actualErr, givenError)
			})
		}
	})
}

func TestMapValue(t *testing.T) {
	mappedValue := "str"
	field2Value := 2

	tests := map[string]interface{}{
		"struct": struct{ Field1, Field2 int }{
			Field1: 1, Field2: field2Value,
		},
		"map": map[string]int{
			"Field1": 1, "Field2": field2Value,
		},
	}

	t.Run("should map struct field", func(t *testing.T) {
		for name, object := range tests {
			t.Run(name, func(t *testing.T) {
				mapper := mapify.Mapper{
					MapValue: func(path string, e mapify.Element) (interface{}, error) {
						if e.Name() == "Field1" {
							return mappedValue, nil
						}

						return e.Interface(), nil
					},
				}
				// when
				v, err := mapper.MapAny(object)
				// then
				require.NoError(t, err)
				expected := map[string]interface{}{
					"Field1": mappedValue,
					"Field2": field2Value,
				}
				assert.Equal(t, expected, v)
			})
		}
	})

	t.Run("should map struct field by path", func(t *testing.T) {
		for name, object := range tests {
			t.Run(name, func(t *testing.T) {
				mapper := mapify.Mapper{
					MapValue: func(path string, e mapify.Element) (interface{}, error) {
						if path == ".Field1" {
							return mappedValue, nil
						}

						return e.Interface(), nil
					},
				}
				// when
				v, err := mapper.MapAny(object)
				// then
				require.NoError(t, err)
				expected := map[string]interface{}{
					"Field1": mappedValue,
					"Field2": field2Value,
				}
				assert.Equal(t, expected, v)
			})
		}
	})

	t.Run("should return error when MapValue returned error", func(t *testing.T) {
		tests := map[string]interface{}{
			"struct": struct{ Field string }{},
			"map":    map[string]string{"Field": ""},
		}

		for name, object := range tests {
			t.Run(name, func(t *testing.T) {
				givenError := stringError("err")
				mapper := mapify.Mapper{
					MapValue: func(path string, e mapify.Element) (interface{}, error) {
						return nil, givenError
					},
				}
				// when
				result, actualErr := mapper.MapAny(object)
				// then
				assert.Nil(t, result)
				assert.ErrorIs(t, actualErr, givenError)
			})
		}
	})
}

type stringError string

func (d stringError) Error() string {
	return string(d)
}

func TestElement_StructField(t *testing.T) {
	t.Run("should return struct field in callbacks", func(t *testing.T) {
		mapper := mapify.Mapper{
			Filter: func(path string, e mapify.Element) (bool, error) {
				assertStructField(t, "Field", e)

				return true, nil
			},
			Rename: func(path string, e mapify.Element) (string, error) {
				assertStructField(t, "Field", e)

				return e.Name(), nil
			},
			MapValue: func(path string, e mapify.Element) (interface{}, error) {
				assertStructField(t, "Field", e)

				return e.Interface(), nil
			},
		}
		_, _ = mapper.MapAny(struct{ Field string }{})
	})
}

func assertStructField(t *testing.T, fieldName string, e mapify.Element) {
	t.Helper()

	field, ok := e.StructField()
	require.True(t, ok)
	assert.Equal(t, fieldName, field.Name)
}
