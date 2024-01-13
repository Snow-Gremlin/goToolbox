package utils

import "github.com/Snow-Gremlin/goToolbox/internal/liteUtils"

// TryIsNil determines if the value is nil
// or returns false if the value isn't able to check for nil.
func TryIsNil(value any) (isNil, ok bool) {
	return liteUtils.TryIsNil(value)
}

// IsNil determines if the value is nil
// or returns false if the value isn't able to check for nil.
func IsNil(value any) bool {
	return liteUtils.IsNil(value)
}
