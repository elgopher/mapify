package main

import (
	"fmt"
	"reflect"
	"time"

	"github.com/elgopher/mapify"
)

// This example shows how to skip conversion of some objects
func main() {
	mapper := mapify.Mapper{
		ShouldConvert: shouldConvert,
	}

	v := struct {
		Time          time.Time // time.Time is a struct too
		DontConvertMe struct{ Field string }
	}{
		Time:          time.Now(),
		DontConvertMe: struct{ Field string }{Field: "v"},
	}

	value, err := mapper.MapAny(v)
	if err != nil {
		panic(err)
	}

	fmt.Println(value)
}

func shouldConvert(path string, v reflect.Value) (bool, error) {
	t := v.Type()

	if t.PkgPath() == "time" && t.Name() == "Time" {
		// time.Time struct will not be converted to map
		return false, nil
	}

	if path == ".DontConvertMe" {
		return false, nil
	}

	return true, nil
}
