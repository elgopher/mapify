// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

// Package mapify converts structs (and other maps) into maps.
package mapify

import (
	"reflect"
	"strconv"
)

// Instance represents instance of mapper
type Instance struct {
	Filter   Filter
	Rename   Rename
	MapValue MapValue
}

// Filter returns true when element should be included.
type Filter func(path string, e Element) bool

// Rename renames element name.
type Rename func(path string, e Element) string

// MapValue map (transform) element value.
type MapValue func(path string, e Element) interface{}

// Element represents either a map entry, field of a struct or unnamed element of a slice.
type Element struct {
	name string
}

func (e Element) Name() string {
	return e.name
}

func (i Instance) MapAny(v interface{}) interface{} {
	return i.newInstance().mapAny("", v)
}

func (i Instance) mapAny(path string, v interface{}) interface{} {
	reflectValue := reflect.ValueOf(v)

	switch {
	case reflectValue.Kind() == reflect.Ptr && reflectValue.Elem().Kind() == reflect.Struct:
		return i.mapAny(path, reflectValue.Elem().Interface())
	case reflectValue.Kind() == reflect.Struct:
		return i.mapStruct(path, reflectValue)
	case reflectValue.Kind() == reflect.Slice:
		return i.mapSlice(path, reflectValue)
	default:
		return v
	}
}

func (i Instance) newInstance() Instance {
	if i.Filter == nil {
		i.Filter = acceptAllFields
	}

	if i.Rename == nil {
		i.Rename = noRename
	}

	return i
}

func (i Instance) mapStruct(path string, reflectValue reflect.Value) map[string]interface{} {
	result := map[string]interface{}{}

	reflectType := reflectValue.Type()

	for j := 0; j < reflectType.NumField(); j++ {
		field := reflectType.Field(j)

		if !field.IsExported() {
			continue
		}

		fieldName := field.Name
		fieldPath := path + "." + fieldName
		element := Element{name: fieldName}

		if i.Filter(fieldPath, element) {
			value := reflectValue.Field(j)
			renamed := i.Rename(fieldPath, element)
			result[renamed] = i.mapAny(fieldPath, value.Interface())
		}
	}

	return result
}

func (i Instance) mapSlice(path string, reflectValue reflect.Value) interface{} {
	kind := reflectValue.Type().Elem().Kind()

	switch kind {
	case reflect.Struct:
		slice := make([]map[string]interface{}, reflectValue.Len())

		for j := 0; j < reflectValue.Len(); j++ {
			slice[j] = i.mapStruct(slicePath(path, j), reflectValue.Index(j))
		}

		return slice
	case reflect.Slice:
		if reflectValue.Type().Elem().Elem().Kind() == reflect.Struct {
			var slice [][]map[string]interface{}

			for j := 0; j < reflectValue.Len(); j++ {
				indexValue := i.mapSlice(slicePath(path, j), reflectValue.Index(j)).([]map[string]interface{})
				slice = append(slice, indexValue)
			}
			return slice
		}
	}

	return reflectValue.Interface()
}

func slicePath(path string, index int) string {
	return path + "[" + strconv.Itoa(index) + "]"
}
