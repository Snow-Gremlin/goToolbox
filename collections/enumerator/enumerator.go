package enumerator

import (
	"strings"

	"goToolbox/collections"
	"goToolbox/collections/iterator"
	"goToolbox/collections/tuple2"
	"goToolbox/terrors/terror"
	"goToolbox/utils"
)

// New creates a new enumerator around the given iterator factory.
// This will pull values from the given iterable.
func New[T any](iterable collections.Iterable[T]) collections.Enumerator[T] {
	return enumeratorImp[T]{
		iterable: iterable,
	}
}

// Enumerate creates an enumerator to return the given values in the given order.
func Enumerate[T any](values ...T) collections.Enumerator[T] {
	return New(func() collections.Iterator[T] {
		return iterator.Iterate(values...)
	})
}

// Range creates an enumerator that counts from he given start the given number of values.
// The range monotonically increments by one from the given start value.
func Range[T utils.NumConstraint](start T, count int) collections.Enumerator[T] {
	return New(func() collections.Iterator[T] {
		return iterator.Range(start, count)
	})
}

// Stride creates an enumerator from the given count of values.
// This will add the given step between each number,
// starting from the given start value.
func Stride[T utils.NumConstraint](start, step T, count int) collections.Enumerator[T] {
	return New(func() collections.Iterator[T] {
		return iterator.Stride(start, step, count)
	})
}

// Repeat creates an enumerator that repeats the given value the given number of times.
func Repeat[T any](value T, count int) collections.Enumerator[T] {
	return New(func() collections.Iterator[T] {
		return iterator.Repeat(value, count)
	})
}

// SplitFunc creates an enumerator that enumerates all the strings from
// splitting the given string with the given separator function.
// The matched separators will not be returned.
func SplitFunc(value string, separator utils.StringMatcher) collections.Enumerator[string] {
	if utils.IsNil(separator) {
		panic(terror.NilArg(`separator`))
	}
	return New(func() collections.Iterator[string] {
		start, count := 0, len(value)
		return iterator.New(func() (string, bool) {
			if start >= count {
				return utils.Zero[string](), false
			}
			index, stride := separator(value[start:])
			if index >= 0 {
				end := min(start+index, count)
				result := value[start:end]
				start = end + stride
				return result, true
			}
			result := value[start:]
			start = count
			return result, true
		})
	})
}

// Split creates an enumerator that enumerates all the string from
// splitting the given string with the given separator.
func Split(value, sep string) collections.Enumerator[string] {
	stride := len(sep)
	return SplitFunc(value, func(part string) (int, int) {
		return strings.Index(part, sep), stride
	})
}

// Lines creates an enumerator that enumerates all the lines in
// the given string. This splits on `\n\r`, `\n`, `\r`, `\r\n`,
// or `\u2029` (paragraph separator, UTF-8 is \xE2\x80\xA9).
// This is useful for reading through a loaded file line by line.
func Lines(value string) collections.Enumerator[string] {
	return SplitFunc(value, func(part string) (int, int) {
		for i, count := 0, len(part); i < count; i++ {
			c := part[i]
			switch c {
			case '\r':
				if i+1 < count && part[i+1] == '\n' {
					return i, 2
				}
				return i, 1
			case '\n':
				if i+1 < count && part[i+1] == '\r' {
					return i, 2
				}
				return i, 1
			case '\xE2':
				if i+2 < count && part[i+1] == '\x80' && part[i+2] == '\xA9' {
					return i, 3
				}
			}
		}
		return -1, 0
	})
}

// Errors creates an enumerator that walks all of the wrapped errors
// inside of the given error.
func Errors(err error) collections.Enumerator[error] {
	return New(func() collections.Iterator[error] {
		return terror.Walk(err)
	})
}

// DoUntilNotZero runs the given function for each value from the given enumerator.
// When a non-zero value is returned by the selector, the enumeration ends and
// returns that non-zero value. If no non-zero value is hit, a zero value is returned.
func DoUntilNotZero[TIn, TOut any](it collections.Enumerator[TIn], s collections.Selector[TIn, TOut]) TOut {
	return iterator.DoUntilNotZero[TIn, TOut](it.Iterate(), s)
}

// Select changes one enumerator type into another by converting each value.
// Typically this is used to select one value out of an enumerated value.
func Select[TIn, TOut any](e collections.Enumerator[TIn], selector collections.Selector[TIn, TOut]) collections.Enumerator[TOut] {
	return New(func() collections.Iterator[TOut] {
		return iterator.Select(e.Iterate(), selector)
	})
}

// Expand creates an enumerator which enumerates through all the values from
// all of the enumerators selected from the values in the given enumerator.
//
// If the expander is expanding objects which are enumerate or enumerable
// then use the `Iterate` method in the enumerator itself as the iterable function.
func Expand[TIn, TOut any, TEnum collections.Iterable[TOut]](e collections.Enumerator[TIn], expander collections.Selector[TIn, TEnum]) collections.Enumerator[TOut] {
	return New(func() collections.Iterator[TOut] {
		return iterator.Expand(e.Iterate(), expander)
	})
}

// Reduce performs a reduction of the values in the given enumerator.
// The reduce method is called with the prior returned value from the previous call.
// The first call is given the initial value.
// The last returned value from reduce is returned. or init if no values.
func Reduce[TIn, TOut any](e collections.Enumerator[TIn], init TOut, reducer collections.Reducer[TIn, TOut]) TOut {
	return iterator.Reduce(e.Iterate(), init, reducer)
}

// SlidingWindow creates an enumerator which handles a sliding window over
// the values from the given enumerator.
//
// The sliding widow is always the specified size. The size must be greater than zero.
// The given stride is how far to advance the window, it must be [1..size].
// If there isn't enough values to fill a window frame then it will not be returned.
// The same slice is reused for the window so modifying it could cause problems and
// it should be copied if the window is kept.
func SlidingWindow[TIn, TOut any](e collections.Enumerator[TIn], size, stride int, window collections.Window[TIn, TOut]) collections.Enumerator[TOut] {
	return New(func() collections.Iterator[TOut] {
		return iterator.SlidingWindow(e.Iterate(), size, stride, window)
	})
}

// Chunk creates an enumerator which has the values grouped into chunks of the given size.
// Each chunk of values are copied into a slice.
// The last chunk may not be the given size, it may be the remaining values.
func Chunk[T any](e collections.Enumerator[T], size int) collections.Enumerator[[]T] {
	return New(func() collections.Iterator[[]T] {
		return iterator.Chunk(e.Iterate(), size)
	})
}

// Sum gets the sum of all value in the given enumerator
// and the number of values that were summed.
func Sum[T utils.NumConstraint](e collections.Enumerator[T]) (T, int) {
	return iterator.Sum(e.Iterate())
}

// IsUnique determines if all the items in the enumerator are unique.
// Returns false if there are any duplicates.
func IsUnique[T comparable](e collections.Enumerator[T]) bool {
	return iterator.IsUnique(e.Iterate())
}

// Unique creates an enumerator that returns only the unique items in the enumerator.
func Unique[T comparable](e collections.Enumerator[T]) collections.Enumerator[T] {
	return New(func() collections.Iterator[T] {
		return iterator.Unique(e.Iterate())
	})
}

// DuplicateCounts creates a map with the given enumerator values as the keys
// and number of duplicates of those values as the map's value.
func DuplicateCounts[T comparable](e collections.Enumerator[T]) map[T]int {
	counters := map[T]int{}
	e.Foreach(func(value T) { counters[value]++ })
	return counters
}

// Zip merges two enumerators together while both enumerators have values
// and returns an enumerator with a combined value of two values from both enumerators.
func Zip[TFirst, TSecond, TOut any](firsts collections.Enumerator[TFirst], seconds collections.Enumerator[TSecond], combiner collections.Combiner[TFirst, TSecond, TOut]) collections.Enumerator[TOut] {
	return New(func() collections.Iterator[TOut] {
		return iterator.Zip(firsts.Iterate(), seconds.Iterate(), combiner)
	})
}

// ZipToTuples merges two enumerators together while both enumerators have values
// and returns an enumerator with a tuple containing values from both enumerators.
func ZipToTuples[TFirst, TSecond any](firsts collections.Enumerator[TFirst], seconds collections.Enumerator[TSecond]) collections.Enumerator[collections.Tuple2[TFirst, TSecond]] {
	return Zip(firsts, seconds, tuple2.New)
}

// Interweave will pull values from each enumerator, one at a time,
// and return them in the cycling order as an enumerator.
// When an enumerator runs out the remaining will interweave until all are empty.
func Interweave[T any](es ...collections.Enumerator[T]) collections.Enumerator[T] {
	return New(func() collections.Iterator[T] {
		its := make([]collections.Iterator[T], len(es))
		for i, e := range es {
			its[i] = e.Iterate()
		}
		return iterator.Interweave(its)
	})
}

// Indexed returns an enumerator that returns a tuple containing the value
// from the given enumerator and an index of the value starting with zero.
func Indexed[T any](e collections.Enumerator[T]) collections.Enumerator[collections.Tuple2[int, T]] {
	return New(func() collections.Iterator[collections.Tuple2[int, T]] {
		return iterator.Indexed(e.Iterate())
	})
}

// OfType creates an enumerator that only enumerates the values of the given target type.
func OfType[Target, T any](e collections.Enumerator[T]) collections.Enumerator[Target] {
	return New(func() collections.Iterator[Target] {
		return iterator.OfType[T, Target](e.Iterate())
	})
}

// Cast creates an enumerator that casts each value from the given enumerator into the
// given target type. If the cast isn't possible then zero is returned for that value.
func Cast[Target, T any](e collections.Enumerator[T]) collections.Enumerator[Target] {
	return New(func() collections.Iterator[Target] {
		return iterator.Cast[T, Target](e.Iterate())
	})
}

// Union creates an enumerator that is the union of the two values.
// No unique values are returned.
func Union[T comparable](left, right collections.Enumerator[T]) collections.Enumerator[T] {
	return Unique(left.Concat(right))
}

// Intersection creates an enumerator that contains the intersection of the two enumerators.
// Only returns values which exists in both enumerators.
func Intersection[T comparable](left, right collections.Enumerator[T]) collections.Enumerator[T] {
	return New(func() collections.Iterator[T] {
		return iterator.Intersection(left.Iterate(), right.Iterate())
	})
}

// Subtract creates an enumerator that contains the all the values from
// the right enumerators but not in the left enumerator.
// This will subtract the left set from the right set.
// Only returns values which exists in only the right enumerator.
func Subtract[T comparable](left, right collections.Enumerator[T]) collections.Enumerator[T] {
	return New(func() collections.Iterator[T] {
		return iterator.Subtract(left.Iterate(), right.Iterate())
	})
}
