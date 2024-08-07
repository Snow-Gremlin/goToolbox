package comp

import "github.com/Snow-Gremlin/goToolbox/internal/liteUtils"

// Equatable is an object which can be checked for equality against another object.
type Equatable interface {
	// Equals returns true if this object and the given object are equal.
	Equals(other any) bool
}

// Equal determines if the two values are equal.
//
// This will check if the objects are Equatable, otherwise it will fallback
// to a DeepEqual. This will not check for Equatable within a slice, array,
// map, etc only in the top level object.
//
// This will not check for Comparable types. Any struct that implements
// Comparable should also implement Equatable where both agree.
func Equal[T any](a, b T) bool {
	return liteUtils.Equal(a, b)
}
