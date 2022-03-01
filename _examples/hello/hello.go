package main

import (
	"fmt"

	"github.com/elgopher/mapify"
)

// This example shows how to convert struct into map
func main() {
	s := SomeStruct{Field: "value"}

	// create Mapper instance. Here default parameters are used.
	mapper := mapify.Mapper{}
	// MapAny maps any object - this can be a struct, slice or map. The whole object is traversed in order to find
	// all nested structs. Each struct will be converted to map[string]interface{}
	result, err := mapper.MapAny(s)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", result) // map[Field:value]
}

type SomeStruct struct {
	Field string
}
