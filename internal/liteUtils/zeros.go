package liteUtils

import "reflect"

// Zero gets the zero value for the given type.
func Zero[T any]() T {
	var zero T
	return zero
}

// IsZero determines if the given value is equivalent
// to the zero value of the given type.
//
// This will also return true for any nil type values.
func IsZero[T any](value T) bool {
	switch any(value).(type) {
	case nil:
		return true
	default:
		return reflect.ValueOf(value).IsZero()
	}
}
