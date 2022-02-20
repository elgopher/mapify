// Package mapify converts structs (and other maps) into maps.
package mapify

import (
	"reflect"
)

// Instance represents instance of mapper
type Instance struct{}

func (i Instance) MapAny(v interface{}) interface{} {
	reflectValue := reflect.ValueOf(v)

	switch {
	case reflectValue.Kind() == reflect.Ptr && reflectValue.Elem().Kind() == reflect.Struct:
		return i.MapAny(reflectValue.Elem().Interface())
	case reflectValue.Kind() == reflect.Struct:
		return i.mapStruct(reflectValue)
	case reflectValue.Kind() == reflect.Slice:
		return i.mapSlice(reflectValue)
	default:
		return v
	}
}

func (i Instance) mapStruct(reflectValue reflect.Value) map[string]interface{} {
	result := map[string]interface{}{}

	reflectType := reflectValue.Type()

	for j := 0; j < reflectType.NumField(); j++ {
		field := reflectType.Field(j)

		if !field.IsExported() {
			continue
		}

		value := reflectValue.Field(j)
		result[field.Name] = i.MapAny(value.Interface())
	}

	return result
}

func (i Instance) mapSlice(reflectValue reflect.Value) interface{} {
	kind := reflectValue.Type().Elem().Kind()

	switch kind {
	case reflect.Struct:
		slice := make([]map[string]interface{}, reflectValue.Len())

		for j := 0; j < reflectValue.Len(); j++ {
			slice[j] = i.mapStruct(reflectValue.Index(j))
		}

		return slice
	case reflect.Slice:
		if reflectValue.Type().Elem().Elem().Kind() == reflect.Struct {
			slice := make([][]map[string]interface{}, reflectValue.Len())

			for j := 0; j < reflectValue.Len(); j++ {
				slice[j] = i.mapSlice(reflectValue.Index(j)).([]map[string]interface{})
			}
			return slice
		}
	}

	return reflectValue.Interface()
}
