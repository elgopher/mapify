package main

import (
	"fmt"

	"github.com/elgopher/mapify"
)

var str = "str"

func main() {
	s := Struct{
		StringField:          str,
		IntField:             3,
		PointerToStringField: &str,
		Nested: Nested{
			Slice: []AnotherNested{
				{Field: "1"},
				{Field: "2"},
			},
		},
	}

	m := mapify.Instance{}.MapAny(s)

	fmt.Printf("%+v", m)
}

type Struct struct {
	StringField          string
	IntField             int
	PointerToStringField *string
	Nested               Nested
}

type Nested struct {
	Slice []AnotherNested
}

type AnotherNested struct {
	Field string
}
