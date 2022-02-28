# Mapify

[![Project Status: Active – The project has reached a stable, usable state and is being actively developed.](https://www.repostatus.org/badges/latest/active.svg)](https://www.repostatus.org/#active)
[![Build](https://github.com/elgopher/mapify/actions/workflows/build.yml/badge.svg)](https://github.com/elgopher/mapify/actions/workflows/build.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/elgopher/mapify.svg)](https://pkg.go.dev/github.com/elgopher/mapify)
[![Go Report Card](https://goreportcard.com/badge/github.com/elgopher/mapify)](https://goreportcard.com/report/github.com/elgopher/mapify)
[![codecov](https://codecov.io/gh/elgopher/mapify/branch/master/graph/badge.svg)](https://codecov.io/gh/elgopher/mapify)

**Highly configurable** struct to map converter. _Will convert maps into other maps as well (work in progress)._

## Features

* **configuration outside the struct**
  * could be in a different package - separation of concerns
  * no need to modify original structs (by adding tags, implementing methods etc.)
  * **behaviour as a code** - you provide code which will be run during conversion
* ability to **rename keys** during conversion
* ability to **omit keys** based on field name, value or tag etc.
* ability to **map elements** during conversion
