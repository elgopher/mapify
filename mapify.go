// Package mapify converts structs (and other maps) into maps.
package mapify

// Instance represents instance of mapper
type Instance struct{}

func (i Instance) MapAny(v interface{}) interface{} {
	return v
}
