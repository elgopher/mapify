package main

import (
	"fmt"
	"strings"

	"github.com/elgopher/mapify"
)

// This example shows how to filter maps and rename keys
func main() {
	s := map[string]interface{}{
		"key":     "value",
		"another": "another value",
	}

	mapper := mapify.Mapper{
		Filter: func(path string, e mapify.Element) (bool, error) {
			return path == ".key", nil
		},
		Rename: func(path string, e mapify.Element) (string, error) {
			return strings.ToUpper(e.Name()), nil
		},
	}

	result, err := mapper.MapAny(s)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", result) // map[KEY:value]
}
