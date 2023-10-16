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
