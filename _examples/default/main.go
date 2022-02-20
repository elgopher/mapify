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

	fmt.Printf("%+v", m) // map[IntField:3 Nested:map[Slice:[map[Field:1] map[Field:2]]] PointerToStringField:0x521e00 StringField:str]
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
