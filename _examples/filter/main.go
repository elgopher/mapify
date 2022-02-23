package main

import (
	"fmt"

	"github.com/elgopher/mapify"
)

func main() {
	s := struct {
		OmittedField1 string
		OmittedField2 string
		VisibleField  string
	}{
		OmittedField1: "hidden",
		OmittedField2: "hidden",
		VisibleField:  "visible",
	}

	mapper := mapify.Mapper{
		Filter: func(path string, e mapify.Element) (bool, error) {
			return path == ".VisibleField", nil
		},
	}

	result, err := mapper.MapAny(s)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", result) // map[VisibleField:visible]
}
