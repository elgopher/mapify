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
	ShouldConvert ShouldConvert
	Filter        Filter
	Rename        Rename
	MapValue      MapValue
}

// ShouldConvert returns true when value should be converted to map. The value can be a struct, map[string]any or slice.
type ShouldConvert func(path string, value reflect.Value) (bool, error)

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
	name  string
	field *reflect.StructField
	reflect.Value
}

// Name returns field name of a struct, key of a map or empty string, when it represents element of a slice.
func (e Element) Name() string {
	return e.name
}

// StructField returns the reflect.StructField if e represents a field of a struct. If not, ok is false.
func (e Element) StructField() (_ reflect.StructField, ok bool) {
	if e.field == nil {
		return reflect.StructField{}, false
	}

	return *e.field, true
}

// MapAny maps any object (struct, map, slice etc.) by converting each struct found to a map.
//
//  * for struct the returned type will be map[string]interface{}
//  * for slice of structs the returned type will be []map[string]interface{}
func (i Mapper) MapAny(v interface{}) (interface{}, error) {
	return i.newInstance().mapAny("", v)
}

func (i Mapper) mapAny(path string, v interface{}) (interface{}, error) {
	reflectValue := reflect.ValueOf(v)

	switch {
	case reflectValue.Kind() == reflect.Struct ||
		(reflectValue.Kind() == reflect.Ptr && reflectValue.Elem().Kind() == reflect.Struct):
		shouldConvert, err := i.ShouldConvert(path, reflectValue)
		if err != nil {
			return nil, fmt.Errorf("ShouldConvert failed: %w", err)
		}

		if !shouldConvert {
			return reflectValue.Interface(), nil
		}

		return i.mapStruct(path, reflectValue)
	case reflectValue.Kind() == reflect.Map && reflectValue.Type().Key().Kind() == reflect.String:
		shouldConvert, err := i.ShouldConvert(path, reflectValue)
		if err != nil {
			return nil, fmt.Errorf("ShouldConvert failed: %w", err)
		}

		if !shouldConvert {
			return reflectValue.Interface(), nil
		}

		return i.mapStringMap(path, reflectValue)
	case reflectValue.Kind() == reflect.Slice:
		return i.mapSlice(path, reflectValue)
	default:
		return v, nil
	}
}

func (i Mapper) newInstance() Mapper {
	if i.ShouldConvert == nil {
		i.ShouldConvert = convertAll
	}

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

	reflectValue = dereference(reflectValue)

	reflectType := reflectValue.Type()

	for j := 0; j < reflectType.NumField(); j++ {
		field := reflectType.Field(j)

		if !field.IsExported() {
			continue
		}

		fieldName := field.Name
		fieldPath := path + "." + fieldName
		value := reflectValue.Field(j)
		element := Element{name: fieldName, Value: value, field: &field}

		if err := i.mapElement(fieldPath, element, result); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func dereference(value reflect.Value) reflect.Value {
	for value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	return value
}

func (i Mapper) mapStringMap(path string, reflectValue reflect.Value) (map[string]interface{}, error) {
	result := map[string]interface{}{}

	keys := reflectValue.MapKeys()
	for _, key := range keys {
		fieldName := key.String()
		fieldPath := path + "." + fieldName
		value := reflectValue.MapIndex(key)
		element := Element{name: fieldName, Value: value}

		if err := i.mapElement(fieldPath, element, result); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (i Mapper) mapElement(fieldPath string, element Element, result map[string]interface{}) error {
	accepted, filterErr := i.Filter(fieldPath, element)
	if filterErr != nil {
		return fmt.Errorf("Filter failed: %w", filterErr)
	}

	if accepted {
		renamed, renameErr := i.Rename(fieldPath, element)
		if renameErr != nil {
			return fmt.Errorf("Rename failed: %w", renameErr)
		}

		mappedValue, mapErr := i.MapValue(fieldPath, element)
		if mapErr != nil {
			return fmt.Errorf("MapValue failed: %w", mapErr)
		}

		finalValue, err := i.mapAny(fieldPath, mappedValue)
		if err != nil {
			return err
		}

		result[renamed] = finalValue
	}

	return nil
}

func (i Mapper) mapSlice(path string, reflectValue reflect.Value) (_ interface{}, err error) {
	kind := reflectValue.Type().Elem().Kind()

	switch kind {
	case reflect.Struct:
		shouldConvert, err := i.ShouldConvert(path, reflectValue)
		if err != nil {
			return nil, fmt.Errorf("ShouldConvert failed: %w", err)
		}

		if !shouldConvert {
			return reflectValue.Interface(), nil
		}

		slice := make([]map[string]interface{}, reflectValue.Len())

		for j := 0; j < reflectValue.Len(); j++ {
			slice[j], err = i.mapStruct(slicePath(path, j), reflectValue.Index(j))
			if err != nil {
				return nil, err
			}
		}

		return slice, nil
	case reflect.Map:
		if reflectValue.Type().Elem().Key().Kind() != reflect.String {
			return reflectValue.Interface(), nil
		}

		shouldConvert, err := i.ShouldConvert(path, reflectValue)
		if err != nil {
			return nil, fmt.Errorf("ShouldConvert failed: %w", err)
		}

		if !shouldConvert {
			return reflectValue.Interface(), nil
		}

		slice := make([]map[string]interface{}, reflectValue.Len())

		for j := 0; j < reflectValue.Len(); j++ {
			slice[j], err = i.mapStringMap(slicePath(path, j), reflectValue.Index(j))
			if err != nil {
				return nil, err
			}
		}

		return slice, nil
	case reflect.Slice:
		sliceElem := reflectValue.Type().Elem().Elem()

		if sliceElem.Kind() == reflect.Struct ||
			(sliceElem.Kind() == reflect.Map && sliceElem.Key().Kind() == reflect.String) {

			shouldConvert, err := i.ShouldConvert(path, reflectValue)
			if err != nil {
				return nil, fmt.Errorf("ShouldConvert failed: %w", err)
			}

			if !shouldConvert {
				return reflectValue.Interface(), nil
			}

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
