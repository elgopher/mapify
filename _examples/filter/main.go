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
		Filter: func(path string, e mapify.Element) bool {
			return path == ".VisibleField"
		},
	}

	m := mapper.MapAny(s)

	fmt.Printf("%+v", m) // map[VisibleField:visible]
}
