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
		if reflectValue.IsNil() {
			return nil
		}

		return i.MapAny(reflectValue.Elem().Interface())
	case reflectValue.Kind() == reflect.Struct:
		return i.mapStruct(reflectValue)
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
