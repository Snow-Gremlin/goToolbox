package utils

import "goToolbox/internal/liteUtils"

// Stringer is an interface for an object which
// returns a string from a `String()` method.
type Stringer interface{ String() string }

// StringMatcher is a function for finding the first match in the given string.
//
// This returns the starting index of the first character (UTF-8) in the match
// and the length of the matched substring.
// If no match found then start should return -1.
// If a match is found the length must be greater than zero.
// For a match, the start index plus the length should not go past the end
// of the given value.
type StringMatcher func(value string) (start, length int)

// GetMaxStringLen gets the maximum length of the given strings.
func GetMaxStringLen(values []string) int {
	maxWidth := 0
	for _, s := range values {
		maxWidth = max(maxWidth, len(s))
	}
	return maxWidth
}

// String gets the string for the given value.
func String[T any](value T) string {
	return liteUtils.String(value)
}

// Strings gets the strings for all the given values in the slice.
func Strings[T any, S ~[]T](s S) []string {
	return liteUtils.Strings(s)
}
