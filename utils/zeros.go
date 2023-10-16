package utils

import (
	"slices"

	"goToolbox/internal/liteUtils"
)

// Zero gets the zero value for the given type.
func Zero[T any]() T {
	return liteUtils.Zero[T]()
}

// IsZero determines if the given value is equivalent
// to the zero value of the given type.
//
// This will also return true for any nil type values.
func IsZero[T any](value T) bool {
	return liteUtils.IsZero[T](value)
}

// RemoveZeros creates a new slice without modifying the
// given values which has all the zero values remove from it.
func RemoveZeros[T any, S ~[]T](s S) S {
	copy := slices.Clone(s)
	copy = slices.DeleteFunc(copy, IsZero)
	return slices.Clip(copy)
}

// SetToZero sets the given range of the slice to zero.
// The start is inclusive, the stop is exclusive.
//
// This will panic if the indices are not valid.
// The start must be `[0..len)` and stop must be `(start..len]`.
func SetToZero[T any, S ~[]T](s S, start, stop int) {
	for i, z := start, Zero[T](); i < stop; i++ {
		s[i] = z
	}
}
