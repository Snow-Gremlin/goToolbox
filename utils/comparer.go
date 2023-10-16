package utils

import (
	"cmp"
	"time"
)

// Comparer returns the comparison result between the two given values.
//
// The comparison results should be:
// `< 0` if x is less than y,
// `== 0` if x equals y,
// `> 0` if x is greater than y.
type Comparer[T any] func(x, y T) int

// DefaultComparer tries to create a comparer for the given type.
// If a comparer isn't able to be created, this will return nil.
func DefaultComparer[T any]() Comparer[T] {
	zero := Zero[T]()
	var c any
	switch any(zero).(type) {
	case int:
		c = OrderedComparer[int]()
	case int8:
		c = OrderedComparer[int8]()
	case int16:
		c = OrderedComparer[int16]()
	case int32:
		c = OrderedComparer[int32]()
	case int64:
		c = OrderedComparer[int64]()
	case uint:
		c = OrderedComparer[uint]()
	case uint8:
		c = OrderedComparer[uint8]()
	case uint16:
		c = OrderedComparer[uint16]()
	case uint32:
		c = OrderedComparer[uint32]()
	case uint64:
		c = OrderedComparer[uint64]()
	case uintptr:
		c = OrderedComparer[uintptr]()
	case float32:
		c = OrderedComparer[float32]()
	case float64:
		c = OrderedComparer[float64]()
	case string:
		c = OrderedComparer[string]()
	case Comparable[T]:
		return flexForComparable[T]()
	case time.Duration:
		c = DurationComparer()
	case time.Time:
		c = TimeComparer()
	default:
		return nil
	}
	return c.(Comparer[T])
}

// flexForComparable returns a comparer which compares the given comparable type.
// This is designed to handle how Go performs type checking in generics for `For`.
func flexForComparable[T any]() Comparer[T] {
	return func(x, y T) int {
		if !IsNil(x) {
			if c, ok := any(x).(Comparable[T]); ok {
				return c.CompareTo(y)
			}
		}
		if !IsNil(y) {
			if c, ok := any(y).(Comparable[T]); ok {
				return -c.CompareTo(x)
			}
		}
		return 0
	}
}

// OrderedComparer returns a comparer which compares the given ordered type.
func OrderedComparer[T cmp.Ordered]() Comparer[T] {
	return cmp.Compare[T]
}

// ComparableComparer returns a comparer which compares the given comparable type.
func ComparableComparer[T Comparable[T]]() Comparer[T] {
	return func(x, y T) int {
		if IsNil(x) {
			if IsNil(y) {
				return 0
			}
			return -y.CompareTo(x)
		}
		return x.CompareTo(y)
	}
}

// DurationComparer returns a comparer which compares the given time duration.
func DurationComparer() Comparer[time.Duration] {
	return func(x, y time.Duration) int {
		return cmp.Compare[int64](int64(x), int64(y))
	}
}

// TimeComparer returns a comparer which compares the given time duration.
func TimeComparer() Comparer[time.Time] {
	return func(x, y time.Time) int {
		return x.Compare(y)
	}
}

// EpsilonComparer returns an epsilon comparer which compares the given floating point types.
//
// An epsilon comparator should be used
// when comparing calculated floating point numbers since calculations may accrue small
// errors and make the actual value very close to the expected value but not exactly equal.
// The two floating point values are equal if very close to each other.
// The values must be within the given epsilon to be considered equal.
//
// The given epsilon must be greater than zero. If the epsilon is
// less than or equal to zero, this will fallback to an ordered comparer.
func EpsilonComparer[T NumConstraint](epsilon T) Comparer[T] {
	if epsilon <= 0 {
		return OrderedComparer[T]()
	}
	return func(a, b T) int {
		if a < b {
			if b-a > epsilon {
				return -1
			}
			return 0
		}
		if a-b > epsilon {
			return 1
		}
		return 0
	}
}

// Descender returns a comparer which negates the given comparer.
//
// Typically negating a comparer will change a sort from ascending to descending.
func Descender[T any](cmp Comparer[T]) Comparer[T] {
	return func(x, y T) int {
		return -cmp(x, y)
	}
}

// LessComparer returns a comparer which compares two values using
// a IsLessThan function to perform the comparison.
func ComparerForLess[T any](less IsLessThan[T]) Comparer[T] {
	return func(x, y T) int {
		switch {
		case less(x, y):
			return -1
		case less(y, x):
			return 1
		default:
			return 0
		}
	}
}
