package utils

import (
	"regexp"
	"strconv"
	"sync"

	"github.com/Snow-Gremlin/goToolbox/internal/liteUtils"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
)

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

// LazyRegex is a regular expression which
// lazy compiles the pattern until the first time it is needed.
//
// If the pattern is invalid, this will panic the first time
// it is called to get the regex. It will panic the same
// error for any following attempts.
func LazyRegex(pattern string) func() *regexp.Regexp {
	return sync.OnceValue(func() *regexp.Regexp {
		return regexp.MustCompile(pattern)
	})
}

// LazyMatcher is a regular expression matcher which
// lazy compiles the pattern until the first time it is needed.
//
// If the pattern is invalid, this will panic the first time
// it is called to perform a match. It will panic the same
// error for any following match attempts.
//
// Example:
//
//	var hex = LazyMatcher(`^[0-9A-Fa-f]+$`)
//	func Foo() {
//		...
//		isHex := hex(`572A6F`) // true
//		...
//	}
func LazyMatcher(pattern string) func(value string) bool {
	r := LazyRegex(pattern)
	return func(value string) bool {
		return r().MatchString(value)
	}
}

// Parse interprets a string and returns the
// corresponding value of the given type.
func Parse[T ParsableConstraint](s string) (T, error) {
	var value any
	var err error
	z := Zero[T]()
	switch any(z).(type) {
	case string:
		value = s
	case bool:
		value, err = strconv.ParseBool(s)
		value = value.(bool)
	case int:
		value, err = strconv.ParseInt(s, 0, 0)
		value = int(value.(int64))
	case int8:
		value, err = strconv.ParseInt(s, 0, 8)
		value = int8(value.(int64))
	case int16:
		value, err = strconv.ParseInt(s, 0, 16)
		value = int16(value.(int64))
	case int32:
		value, err = strconv.ParseInt(s, 0, 32)
		value = int32(value.(int64))
	case int64:
		value, err = strconv.ParseInt(s, 0, 64)
	case uint:
		value, err = strconv.ParseUint(s, 0, 0)
		value = uint(value.(uint64))
	case uint8:
		value, err = strconv.ParseUint(s, 0, 8)
		value = uint8(value.(uint64))
	case uint16:
		value, err = strconv.ParseUint(s, 0, 16)
		value = uint16(value.(uint64))
	case uint32:
		value, err = strconv.ParseUint(s, 0, 32)
		value = uint32(value.(uint64))
	case uint64:
		value, err = strconv.ParseUint(s, 0, 64)
	case float32:
		value, err = strconv.ParseFloat(s, 32)
		value = float32(value.(float64))
	case float64:
		value, err = strconv.ParseFloat(s, 64)
	case complex64:
		value, err = strconv.ParseComplex(s, 64)
		value = complex64(value.(complex128))
	case complex128:
		value, err = strconv.ParseComplex(s, 128)
	}
	if err != nil {
		return z, terror.New(`unable to parse value`, err).
			With(`input`, s).
			With(`type`, TypeOf[T]())
	}
	return value.(T), err
}
