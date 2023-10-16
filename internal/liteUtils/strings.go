package liteUtils

import "fmt"

type stringer interface{ String() string }

// String gets the string for the given value.
func String[T any](value T) string {
	switch t := any(value).(type) {
	case nil:
		return `<nil>`
	case string:
		return t
	case error:
		return t.Error()
	case stringer:
		return t.String()
	default:
		return fmt.Sprint(value)
	}
}

// Strings gets the strings for all the given values in the slice.
func Strings[T any, S ~[]T](s S) []string {
	result := make([]string, len(s))
	for i, v := range s {
		result[i] = String(v)
	}
	return result
}
