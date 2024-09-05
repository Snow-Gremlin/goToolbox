package utils

import "reflect"

// IntConstraint is a type constraint for integer types.
type IntConstraint interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// FloatingConstraint is a type constraint for floating point types.
type FloatingConstraint interface {
	~float32 | ~float64
}

// NumConstraint is a type constraint for numerical types.
type NumConstraint interface {
	IntConstraint | FloatingConstraint
}

// ParsableConstraint is the set of types that can be parsed.
type ParsableConstraint interface {
	~string | ~bool | NumConstraint | ~complex64 | ~complex128
}

// TypeOf gets the reflect type of the given generic value.
func TypeOf[T any]() reflect.Type {
	var zero [0]T
	return reflect.TypeOf(zero).Elem()
}

// Is tries to cast the given value `v` into the given type.
// If `out` is not nil, then it will set with the cast value.
// Returns true if successfully cast, false otherwise.
//
// This is designed to allow for inlining a cast into an if-statement
// without having to start another if-statement for each cast.
func Is[T any](v any, out *T) bool {
	if u, ok := v.(T); ok {
		if !IsNil(out) {
			*out = u
		}
		return true
	}
	return false
}
