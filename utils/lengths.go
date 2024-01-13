package utils

import "reflect"

// getLengthFallback checks if the given value has a
// length or count method that can be used to determine the length.
func getLengthFallback(value any) (int, bool) {
	switch t := value.(type) {
	case interface{ Count() int }:
		return t.Count(), true
	case interface{ Len() int }:
		return t.Len(), true
	case interface{ Length() int }:
		return t.Length(), true
	default:
		return 0, false
	}
}

// Length determines the length of the given
// value or returns false if the value does not have a length.
func Length[T any](value T) (length int, ok bool) {
	defer func() {
		if r := recover(); r != nil {
			length, ok = getLengthFallback(value)
		}
	}()
	return reflect.ValueOf(value).Len(), true
}
