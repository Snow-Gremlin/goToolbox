package comp

import (
	"cmp"
	"time"

	"github.com/Snow-Gremlin/goToolbox/internal/liteUtils"
)

// Default tries to create a comparer for the given type.
// If a comparer isn't able to be created, this will return nil.
func Default[T any]() Comparer[T] {
	zero := liteUtils.Zero[T]()
	var c any
	switch any(zero).(type) {
	case bool:
		c = Bool()
	case int:
		c = Ordered[int]()
	case int8:
		c = Ordered[int8]()
	case int16:
		c = Ordered[int16]()
	case int32:
		c = Ordered[int32]()
	case int64:
		c = Ordered[int64]()
	case uint:
		c = Ordered[uint]()
	case uint8:
		c = Ordered[uint8]()
	case uint16:
		c = Ordered[uint16]()
	case uint32:
		c = Ordered[uint32]()
	case uint64:
		c = Ordered[uint64]()
	case uintptr:
		c = Ordered[uintptr]()
	case float32:
		c = Ordered[float32]()
	case float64:
		c = Ordered[float64]()
	case string:
		c = Ordered[string]()
	case Comparable[T]:
		return flexForComparable[T]()
	case time.Duration:
		c = Duration()
	case time.Time:
		c = Time()
	default:
		return nil
	}
	return c.(Comparer[T])
}

// Bool is a comparer that compares the given boolean values.
//
// |  x  |  y  | result |
// |:---:|:---:|:------:|
// |  F  |  F  |    0   |
// |  F  |  T  |   -1   |
// |  T  |  F  |    1   |
// |  T  |  T  |    0   |
func Bool() Comparer[bool] {
	return func(x, y bool) int {
		result := 0
		if x {
			result++
		}
		if y {
			result--
		}
		return result
	}
}

// flexForComparable returns a comparer which compares the given comparable type.
// This is designed to handle how Go performs type checking when switching on type.
func flexForComparable[T any]() Comparer[T] {
	return func(x, y T) int {
		if !liteUtils.IsNil(x) {
			if c, ok := any(x).(Comparable[T]); ok {
				return c.CompareTo(y)
			}
		}
		if !liteUtils.IsNil(y) {
			if c, ok := any(y).(Comparable[T]); ok {
				return -c.CompareTo(x)
			}
		}
		return 0
	}
}

// Ordered returns a comparer which compares the given ordered type.
func Ordered[T cmp.Ordered]() Comparer[T] {
	return cmp.Compare[T]
}

// ComparableComparer returns a comparer which compares the given comparable type.
func ComparableComparer[T Comparable[T]]() Comparer[T] {
	return func(x, y T) int {
		if liteUtils.IsNil(x) {
			if liteUtils.IsNil(y) {
				return 0
			}
			return -y.CompareTo(x)
		}
		return x.CompareTo(y)
	}
}

// Duration returns a comparer which compares the given time duration.
func Duration() Comparer[time.Duration] {
	return func(x, y time.Duration) int {
		return cmp.Compare[int64](int64(x), int64(y))
	}
}

// Time returns a comparer which compares the given time duration.
func Time() Comparer[time.Time] {
	return func(x, y time.Time) int {
		return x.Compare(y)
	}
}

// Num is any number value type.
type Num interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~uintptr | ~float32 | ~float64
}

// Epsilon returns an epsilon comparer which compares the given floating point types.
//
// An epsilon comparator should be used
// when comparing calculated floating point numbers since calculations may accrue small
// errors and make the actual value very close to the expected value but not exactly equal.
// The two floating point values are equal if very close to each other.
// The values must be within the given epsilon to be considered equal.
//
// The given epsilon must be greater than zero. If the epsilon is
// less than or equal to zero, this will fallback to an ordered comparer.
func Epsilon[T Num](epsilon T) Comparer[T] {
	if epsilon <= 0 {
		return Ordered[T]()
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

// FromLess returns a comparer which compares two values using
// a IsLessThan function to perform the comparison.
func FromLess[T any](less IsLessThan[T]) Comparer[T] {
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

// Slice compares two slices with the given element comparer.
func Slice[S ~[]T, T any](elemCmp Comparer[T]) Comparer[S] {
	return func(a, b S) int {
		ca, cb := len(a), len(b)
		cMin := min(ca, cb)
		for i := 0; i < cMin; i++ {
			if cmp := elemCmp(a[i], b[i]); cmp != 0 {
				return cmp
			}
		}
		return cmp.Compare(ca, cb)
	}
}

// Or will return the first non-zero value returned
// by a comparison or it will return zero.
//
// The given functions will only be evaluated if all
// prior tests have returned zero.
func Or(comps ...func() int) int {
	for _, cmp := range comps {
		if c := cmp(); c != 0 {
			return c
		}
	}
	return 0
}
