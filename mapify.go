// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

// Package mapify converts structs (and other maps) into maps.
package mapify

import (
	"fmt"
	"reflect"
	"strconv"
)

// Mapper represents instance of mapper
type Mapper struct {
	Filter   Filter
	Rename   Rename
	MapValue MapValue
}

// Filter returns true when element should be included. If error is returned then the whole conversion is aborted
// and wrapped error is returned from Mapper.MapAny method.
type Filter func(path string, e Element) (bool, error)

// Rename renames element name. If error is returned then the whole conversion is aborted
// and wrapped error is returned from Mapper.MapAny method.
type Rename func(path string, e Element) (string, error)

// MapValue maps (transforms) element value. If error is returned then the whole conversion is aborted
// and wrapped error is returned from Mapper.MapAny method.
type MapValue func(path string, e Element) (interface{}, error)

// Element represents either a map entry, field of a struct or unnamed element of a slice.
type Element struct {
	name string
	reflect.Value
}

func (e Element) Name() string {
	return e.name
}

func (i Mapper) MapAny(v interface{}) (interface{}, error) {
	return i.newInstance().mapAny("", v)
}

func (i Mapper) mapAny(path string, v interface{}) (interface{}, error) {
	reflectValue := reflect.ValueOf(v)

	switch {
	case reflectValue.Kind() == reflect.Ptr && reflectValue.Elem().Kind() == reflect.Struct:
		return i.mapAny(path, reflectValue.Elem().Interface())
	case reflectValue.Kind() == reflect.Struct:
		return i.mapStruct(path, reflectValue)
	case reflectValue.Kind() == reflect.Slice:
		return i.mapSlice(path, reflectValue)
	default:
		return v, nil
	}
}

func (i Mapper) newInstance() Mapper {
	if i.Filter == nil {
		i.Filter = acceptAllFields
	}

	if i.Rename == nil {
		i.Rename = noRename
	}

	if i.MapValue == nil {
		i.MapValue = interfaceValue
	}

	return i
}

func (i Mapper) mapStruct(path string, reflectValue reflect.Value) (map[string]interface{}, error) {
	result := map[string]interface{}{}

	reflectType := reflectValue.Type()

	for j := 0; j < reflectType.NumField(); j++ {
		field := reflectType.Field(j)

		if !field.IsExported() {
			continue
		}

		fieldName := field.Name
		fieldPath := path + "." + fieldName
		value := reflectValue.Field(j)
		element := Element{name: fieldName, Value: value}

		accepted, err := i.Filter(fieldPath, element)
		if err != nil {
			return nil, fmt.Errorf("Filter failed: %w", err)
		}

		if accepted {
			renamed, err := i.Rename(fieldPath, element)
			if err != nil {
				return nil, fmt.Errorf("Rename failed: %w", err)
			}

			mappedValue, err := i.MapValue(fieldPath, element)
			if err != nil {
				return nil, fmt.Errorf("MapValue failed: %w", err)
			}

			result[renamed], err = i.mapAny(fieldPath, mappedValue)
			if err != nil {
				return nil, err
			}
		}
	}

	return result, nil
}

func (i Mapper) mapSlice(path string, reflectValue reflect.Value) (_ interface{}, err error) {
	kind := reflectValue.Type().Elem().Kind()

	switch kind {
	case reflect.Struct:
		slice := make([]map[string]interface{}, reflectValue.Len())

		for j := 0; j < reflectValue.Len(); j++ {
			slice[j], err = i.mapStruct(slicePath(path, j), reflectValue.Index(j))
			if err != nil {
				return nil, err
			}
		}

		return slice, nil
	case reflect.Slice:
		if reflectValue.Type().Elem().Elem().Kind() == reflect.Struct {
			var slice [][]map[string]interface{}

			for j := 0; j < reflectValue.Len(); j++ {
				indexValue, err := i.mapSlice(slicePath(path, j), reflectValue.Index(j))
				if err != nil {
					return nil, err
				}

				slice = append(slice, indexValue.([]map[string]interface{}))
			}

			return slice, nil
		}
	}

	return reflectValue.Interface(), nil
}

func slicePath(path string, index int) string {
	return path + "[" + strconv.Itoa(index) + "]"
}
