package liteUtils

import "reflect"

type equatable interface{ Equals(other any) bool }

// Equal determines if the two values are equal.
//
// This will check if the objects are Equatable, otherwise it will fallback
// to a DeepEqual. This will not check for Equatable within a slice, array,
// map, etc only in the top level object.
//
// This will not check for Comparable types. Any struct that implements
// Comparable should also implement Equatable where both agree.
func Equal[T any](a, b T) bool {
	if IsNil(a) {
		return reflect.DeepEqual(a, b)
	}
	if IsNil(b) {
		return reflect.DeepEqual(a, b)
	}
	if ae, ok := any(a).(equatable); ok {
		return ae.Equals(b)
	}
	if be, ok := any(b).(equatable); ok {
		return be.Equals(a)
	}
	return reflect.DeepEqual(a, b)
}
