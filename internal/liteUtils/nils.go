package liteUtils

import "reflect"

// TryIsNil determines if the value is nil
// or returns false if the value isn't able to check for nil.
func TryIsNil(value any) (isNil, ok bool) {
	defer func() {
		if r := recover(); r != nil {
			isNil = false
			ok = false
		}
	}()

	switch value.(type) {
	case nil:
		return true, true
	default:
		return reflect.ValueOf(value).IsNil(), true
	}
}

// IsNil determines if the value is nil
// or returns false if the value isn't able to check for nil.
func IsNil(value any) bool {
	isNil, _ := TryIsNil(value)
	return isNil
}
