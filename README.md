# Mapify

[![Project Status: Active – The project has reached a stable, usable state and is being actively developed.](https://www.repostatus.org/badges/latest/active.svg)](https://www.repostatus.org/#active)
[![Build](https://github.com/elgopher/mapify/actions/workflows/build.yml/badge.svg)](https://github.com/elgopher/mapify/actions/workflows/build.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/elgopher/mapify.svg)](https://pkg.go.dev/github.com/elgopher/mapify)
[![Go Report Card](https://goreportcard.com/badge/github.com/elgopher/mapify)](https://goreportcard.com/report/github.com/elgopher/mapify)
[![codecov](https://codecov.io/gh/elgopher/mapify/branch/master/graph/badge.svg)](https://codecov.io/gh/elgopher/mapify)

**Highly configurable** struct to map converter. Converts `map[string]any` into other maps as well.

## Features

* **configuration outside the struct**
  * could be in a different package
  * no need to modify original structs (by adding tags, implementing methods etc.)
  * **behaviour as a code** - you provide code which will be run during conversion
* ability to:
  * **rename keys** during conversion
  * **omit keys** based on field name, value or tag etc.
  * **map elements** during conversion

## Installation

```shell
# Add mapify to your Go module:
go get github.com/elgopher/mapify        
```

## Hello, world!

```go
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
```

## MapAny algorithm

1. Take an object which is a struct, map or slice.
2. Traverse entire object looking for nested structs or maps.
3. **Filter** elements (struct fields and map keys).
4. **Rename** field names or map keys.
5. **Map** (struct field or map values).
