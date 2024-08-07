package iterator

import (
	"reflect"
	"slices"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/predicate"
	"github.com/Snow-Gremlin/goToolbox/collections/tuple2"
	"github.com/Snow-Gremlin/goToolbox/comp"
	"github.com/Snow-Gremlin/goToolbox/internal/optional"
	"github.com/Snow-Gremlin/goToolbox/internal/simpleSet"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

// New creates a new iterator for stepping through values using a fetcher.
// As soon as the fetcher returns false this iterator will stop.
func New[T any](fetcher Fetcher[T]) collections.Iterator[T] {
	return &iteratorImp[T]{
		fetcher: fetcher,
		current: utils.Zero[T](),
	}
}

// Iterate will iterate the given values.
func Iterate[T any](values ...T) collections.Iterator[T] {
	index := -1
	count := len(values)
	return New(func() (T, bool) {
		if count > 0 {
			count--
			index++
			return values[index], true
		}
		return utils.Zero[T](), false
	})
}

// Range will iterate from the given count of values.
// The range monotonically increments by one from the given start value.
func Range[T utils.NumConstraint](start T, count int) collections.Iterator[T] {
	return Stride(start, 1, count)
}

// Stride will iterate from the given count of values.
// This will add the given step between each number,
// starting from the given start value.
func Stride[T utils.NumConstraint](start, step T, count int) collections.Iterator[T] {
	return New(func() (T, bool) {
		if count > 0 {
			count--
			value := start
			start += step
			return value, true
		}
		return utils.Zero[T](), false
	})
}

// Repeat will iterate the same value the given count.
func Repeat[T any](value T, count int) collections.Iterator[T] {
	return New(func() (T, bool) {
		if count > 0 {
			count--
			return value, true
		}
		return utils.Zero[T](), false
	})
}

// Where creates an iterator which reads from the given iterator
// but only returns values which the given predicate returns true for.
func Where[T any](it collections.Iterator[T], p collections.Predicate[T]) collections.Iterator[T] {
	return New(func() (T, bool) {
		for it.Next() {
			if value := it.Current(); p(value) {
				return value, true
			}
		}
		return utils.Zero[T](), false
	})
}

// WhereNot creates an iterator which reads from the given iterator
// but only returns values which the given predicate returns false for.
func WhereNot[T any](it collections.Iterator[T], p collections.Predicate[T]) collections.Iterator[T] {
	return Where(it, predicate.Not(p))
}

// ToSlice reads all the values out of the given iterator and returns them as a slice.
func ToSlice[T any](it collections.Iterator[T]) []T {
	s := []T{}
	for it.Next() {
		s = append(s, it.Current())
	}
	return s
}

// CopyToSlice reads values out of the given iterator adds them to the given slice.
// This will stop when either values run out or there is no more room in the slice.
func CopyToSlice[T any](it collections.Iterator[T], s []T) {
	for i, room := 0, len(s); i < room && it.Next(); i++ {
		s[i] = it.Current()
	}
}

// Foreach runs the given function for each values from the given iterator.
func Foreach[T any](it collections.Iterator[T], m func(value T)) {
	for it.Next() {
		m(it.Current())
	}
}

// DoUntilError runs the given function for each value from the given iterator.
// When an error is returned by the selector, the iteration ends and
// returns that error. If no error is hit, nil is returned.
func DoUntilError[T any](it collections.Iterator[T], s collections.Selector[T, error]) error {
	for it.Next() {
		if err := s(it.Current()); err != nil {
			return err
		}
	}
	return nil
}

// DoUntilNotZero runs the given function for each value from the given iterator.
// When a non-zero value is returned by the selector, the iteration ends and
// returns that non-zero value. If no non-zero value is hit, a zero value is returned.
func DoUntilNotZero[TIn, TOut any](it collections.Iterator[TIn], s collections.Selector[TIn, TOut]) TOut {
	for it.Next() {
		if v := s(it.Current()); !utils.IsZero(v) {
			return v
		}
	}
	return utils.Zero[TOut]()
}

// Any reads values from the given iterator until one of the values causes
// the given predicate to return true, then true is returned.
// If the predicated returns false for all values, then false is returned.
func Any[T any](it collections.Iterator[T], p collections.Predicate[T]) bool {
	for it.Next() {
		if p(it.Current()) {
			return true
		}
	}
	return false
}

// All reads value from the given iterator until one of the values causes
// the given predicate to return false, then false is returned.
// If the predicated returns true for all values, then true is returned.
func All[T any](it collections.Iterator[T], p collections.Predicate[T]) bool {
	for it.Next() {
		if !p(it.Current()) {
			return false
		}
	}
	return true
}

// StepsUntil determines the number of values in the iterator are read until
// a value satisfies the given predicate.
//
// The count will not include the value which satisfied the predicate such that
// if the first value satisfies the predicate then this will return zero.
// If no value satisfies the predicate then -1 is returned.
func StepsUntil[T any](it collections.Iterator[T], p collections.Predicate[T]) int {
	count := 0
	for it.Next() {
		if p(it.Current()) {
			return count
		}
		count++
	}
	return -1
}

// StartsWith determines if the first iterator starts with the given prefix.
func StartsWith[T any](it, prefix collections.Iterator[T]) bool {
	for {
		n1, n2 := it.Next(), prefix.Next()
		if !n1 {
			return !n2
		}
		if !n2 {
			return true
		}
		if !comp.Equal(it.Current(), prefix.Current()) {
			return false
		}
	}
}

// Equal determines if the two iterators contain the same values.
func Equal[T any](it1, it2 collections.Iterator[T]) bool {
	for {
		next1, next2 := it1.Next(), it2.Next()
		if !next1 {
			return !next2
		}
		if !next2 || !comp.Equal(it1.Current(), it2.Current()) {
			return false
		}
	}
}

// Empty attempts to read the next value off the given iterator.
// If a value exists, the iterator isn't empty, otherwise false.
func Empty[T any](it collections.Iterator[T]) bool {
	return !it.Next()
}

// Count reads all the values from the given iterator and
// returns how many values were in the iterator.
func Count[T any](it collections.Iterator[T]) int {
	count := 0
	for it.Next() {
		count++
	}
	return count
}

// AtLeast reads only enough values from the given iterator to determine
// if there is at least the given number of values exists.
func AtLeast[T any](it collections.Iterator[T], min int) bool {
	count := 0
	for it.Next() {
		count++
		if count >= min {
			return true
		}
	}
	return false
}

// AtMost reads all the values from the given iterator to determine
// if there is at most the given number of values exists.
func AtMost[T any](it collections.Iterator[T], max int) bool {
	count := 0
	for it.Next() {
		count++
		if count > max {
			return false
		}
	}
	return true
}

// First reads one value off the given iterator, if one exists,
// otherwise the zero value is returned.
func First[T any](it collections.Iterator[T]) (T, bool) {
	if it.Next() {
		return it.Current(), true
	}
	return utils.Zero[T](), false
}

// Last reads all the values off the given iterator and returns the last value.
// Returns zero and false if there were no values in the iterator.
func Last[T any](it collections.Iterator[T]) (T, bool) {
	found := false
	last := utils.Zero[T]()
	for it.Next() {
		found = true
		last = it.Current()
	}
	return last, found
}

// Single reads two values off the given iterator.
// If there is only one, then it is returned.
func Single[T any](it collections.Iterator[T]) (T, bool) {
	if it.Next() {
		value := it.Current()
		if !it.Next() {
			return value, true
		}
	}
	return utils.Zero[T](), false
}

// Skip creates a new iterator which skips over the first given count number of values
// from the given iterator before iterating the remaining values.
// The skipped values aren't read until the first value read from the returned iterator.
func Skip[T any](it collections.Iterator[T], count int) collections.Iterator[T] {
	return New(func() (T, bool) {
		for ; count > 0; count-- {
			if !it.Next() {
				return utils.Zero[T](), false
			}
		}
		if it.Next() {
			return it.Current(), true
		}
		return utils.Zero[T](), false
	})
}

// SkipWhile creates a new iterator which skips over values until the given predicate returns false.
// The values from the given iterator are returned after and including the first false from the predicate.
// The skipped values aren't read until the first value read from the returned iterator.
func SkipWhile[T any](it collections.Iterator[T], p collections.Predicate[T]) collections.Iterator[T] {
	return New(func() (T, bool) {
		if p != nil {
			for it.Next() {
				if value := it.Current(); !p(value) {
					p = nil
					return value, true
				}
			}
			return utils.Zero[T](), false
		}

		if it.Next() {
			return it.Current(), true
		}
		return utils.Zero[T](), false
	})
}

// Take creates a new iterator which only takes the given count of values
// form the given iterator before stopping iteration.
func Take[T any](it collections.Iterator[T], count int) collections.Iterator[T] {
	return New(func() (T, bool) {
		if count > 0 && it.Next() {
			count--
			return it.Current(), true
		}
		count = 0
		return utils.Zero[T](), false
	})
}

// TakeWhile creates a new iterator which only takes values until the given predicate returns false.
// The values from the given iterator are returned until and excluding the first false from the predicate.
func TakeWhile[T any](it collections.Iterator[T], p collections.Predicate[T]) collections.Iterator[T] {
	return New(func() (T, bool) {
		if p != nil && it.Next() {
			if value := it.Current(); p(value) {
				return value, true
			}
		}
		p = nil
		return utils.Zero[T](), false
	})
}

// Replace creates a new iterator which checks each value using the given replacer.
// This gives an opportunity for values to be replaced with a different value or returned unchanged
// whilst iterating over the values from the given iterator.
func Replace[T any](it collections.Iterator[T], replacer collections.Selector[T, T]) collections.Iterator[T] {
	return New(func() (T, bool) {
		if it.Next() {
			return replacer(it.Current()), true
		}
		return utils.Zero[T](), false
	})
}

// Reverse reads all the values from the given iterator and returns them
// in reverse order. No values are read from the given iterator,
// until a value is read from the new iterator.
//
// Try to maximize the number of values used after a reverse, since
// it is a waste of time and memory to reverse a large iteration if only
// the first three values (the last three in the iteration) are used.
// Reducing down to a smaller iterator where all the values are needed
// in reverse is better than reversing unneeded values.
func Reverse[T any](it collections.Iterator[T]) collections.Iterator[T] {
	first := true
	var index int
	var values []T
	return New(func() (T, bool) {
		if first {
			first = false
			values = ToSlice[T](it)
			index = len(values)
		}
		index--
		if index >= 0 {
			return values[index], true
		}
		values = nil
		return utils.Zero[T](), false
	})
}

// Append creates an iterator with the given value appended to the end of the values.
func Append[T any](it collections.Iterator[T], tails []T) collections.Iterator[T] {
	index, count := 0, len(tails)
	return New(func() (T, bool) {
		if it.Next() {
			return it.Current(), true
		}
		if index < count {
			tail := tails[index]
			index++
			return tail, true
		}
		return utils.Zero[T](), false
	})
}

// Concat concatenates the given iterators into one iterator.
func Concat[T any](its []collections.Iterator[T]) collections.Iterator[T] {
	index, count := 0, len(its)
	return New(func() (T, bool) {
		for ; index < count; index++ {
			if it := its[index]; it.Next() {
				return it.Current(), true
			}
		}
		return utils.Zero[T](), false
	})
}

// Select changes one iterator type into another by converting each value.
// Typically this is used to select one value out of a value.
func Select[TIn, TOut any](
	it collections.Iterator[TIn],
	selector collections.Selector[TIn, TOut],
) collections.Iterator[TOut] {
	return New(func() (TOut, bool) {
		if it.Next() {
			return selector(it.Current()), true
		}
		return utils.Zero[TOut](), false
	})
}

// OfType returns only the values of the given out type.
func OfType[TIn, TOut any](it collections.Iterator[TIn]) collections.Iterator[TOut] {
	return New(func() (TOut, bool) {
		for it.Next() {
			if target, ok := any(it.Current()).(TOut); ok {
				return target, true
			}
		}
		return utils.Zero[TOut](), false
	})
}

// tryCast will attempt to convert the given value into the given target type.
// The target type object must match the generic TOut type.
func tryCast[TIn, TOut any](value TIn, target reflect.Type) (result TOut, ok bool) {
	defer func() {
		if r := recover(); r != nil {
			ok = false
		}
	}()

	result, ok = reflect.ValueOf(value).Convert(target).Interface().(TOut)
	return result, ok
}

// Cast casts the given type into the given out type
// If a cast isn't possible then zero is returned.
func Cast[TIn, TOut any](it collections.Iterator[TIn]) collections.Iterator[TOut] {
	target := utils.TypeOf[TOut]()
	return New(func() (TOut, bool) {
		if it.Next() {
			value, ok := tryCast[TIn, TOut](it.Current(), target)
			if !ok {
				return utils.Zero[TOut](), true
			}
			return value, true
		}
		return utils.Zero[TOut](), false
	})
}

// Expand creates an iterator which iterates all the values from all the iterators
// selected from the values in the given iterator.
func Expand[TIn, TOut any, TEnum collections.Iterable[TOut]](
	it collections.Iterator[TIn],
	expander collections.Selector[TIn, TEnum],
) collections.Iterator[TOut] {
	var current collections.Iterator[TOut]
	return New(func() (TOut, bool) {
		for {
			if current != nil {
				if current.Next() {
					return current.Current(), true
				}
				current = nil
			}
			if it.Next() {
				iterFactory := expander(it.Current())
				if !utils.IsNil(iterFactory) {
					current = iterFactory()
				}
				continue
			}
			return utils.Zero[TOut](), false
		}
	})
}

// Reduce reads all the values from the given iterator.
// The reduce method is called with the prior returned value from the previous call.
// The first call is given the initial value.
// The last returned value from reduce is returned. or init if no values.
func Reduce[TIn, TOut any](it collections.Iterator[TIn], init TOut, reducer collections.Reducer[TIn, TOut]) TOut {
	prior := init
	for it.Next() {
		prior = reducer(it.Current(), prior)
	}
	return prior
}

// Merge reads all the values from the given iterator together.
// The merge method is called with the prior returned value from the previous call.
// The first value is used as the prior value with the second value in the merger.
// The last returned value from merge is returned, the first value if there
// is only one value, or the zero value if no values.
func Merge[T any](it collections.Iterator[T], merger collections.Reducer[T, T]) T {
	prior := utils.Zero[T]()
	first := true
	for it.Next() {
		if first {
			first = false
			prior = it.Current()
		} else {
			prior = merger(it.Current(), prior)
		}
	}
	return prior
}

// SlidingWindow creates an iterator which handles a sliding window over
// the values from the given iterator.
//
// The sliding widow is always the specified size. The size must be greater than zero.
// The given stride is how far to advance the window, it must be [1..size].
// The same slice is reused for the window so modifying it could cause problems and
// it should be copied if the window is kept.
func SlidingWindow[TIn, TOut any](
	it collections.Iterator[TIn],
	size, stride int,
	window collections.Window[TIn, TOut],
) collections.Iterator[TOut] {
	if size <= 0 {
		panic(terror.New(`the given window size must be greater than zero`).
			With(`size`, size))
	}
	if stride <= 0 || stride > size {
		panic(terror.New(`the given window stride must be greater than zero and less than or equal to the size`).
			With(`size`, size).
			With(`stride`, stride))
	}

	frame := make([]TIn, size)
	loadIndex := 0
	return New(func() (TOut, bool) {
		for loadIndex < size {
			if !it.Next() {
				return utils.Zero[TOut](), false
			}
			frame[loadIndex] = it.Current()
			loadIndex++
		}

		result := window(frame)
		if stride < size {
			copy(frame, frame[stride:])
		}
		loadIndex -= stride
		return result, true
	})
}

// Chunk creates an iterator which has the values grouped into chunks of the given size.
// Each chunk is the given size of values put into a slice. The chunks have already been copied.
// The last chunk may not be the given size, it may be the remaining values.
func Chunk[T any](it collections.Iterator[T], size int) collections.Iterator[[]T] {
	if size <= 0 {
		panic(terror.New(`the given chunk size must be greater than zero`).
			With(`size`, size))
	}

	frame := make([]T, size)
	loadIndex := 0
	return New(func() ([]T, bool) {
		for loadIndex < size {
			if !it.Next() {
				if loadIndex > 0 {
					result := slices.Clone(frame[:loadIndex])
					loadIndex = 0
					return result, true
				}
				return utils.Zero[[]T](), false
			}
			frame[loadIndex] = it.Current()
			loadIndex++
		}

		result := slices.Clone(frame)
		loadIndex = 0
		return result, true
	})
}

// Sum gets the sum of all value in the given iterator
// and the number of values that were summed.
func Sum[T utils.NumConstraint](it collections.Iterator[T]) (T, int) {
	var sum T
	count := 0
	for it.Next() {
		sum += it.Current()
		count++
	}
	return sum, count
}

// IsUnique determines if there are only unique items in the iterator.
// Returns false if there are duplicate values in the iterator.
func IsUnique[T comparable](it collections.Iterator[T]) bool {
	touched := simpleSet.New[T]()
	return All(it, touched.SetTest)
}

// Unique creates an iterator that returns only the unique items in the iterator.
func Unique[T comparable](it collections.Iterator[T]) collections.Iterator[T] {
	touched := simpleSet.New[T]()
	return Where(it, touched.SetTest)
}

// Intersection creates an iterator that returns only the values which exists in both iterators.
//
// The right iterator takes precedence over the result such that it
// determines the order and if there are repeats in the result.
func Intersection[T comparable](left, right collections.Iterator[T]) collections.Iterator[T] {
	inLeft := simpleSet.New[T]()
	return Where(right, func(value T) bool {
		if inLeft.Has(value) {
			return true
		}

		if left != nil {
			for left.Next() {
				cur := left.Current()
				inLeft.Set(cur)
				if value == cur {
					return true
				}
			}
			left = nil
		}

		return false
	})
}

// Subtract creates an iterator that returns only the values which exists
// in the right iterator but not in the left.
// This will subtract the left set from the right set.
//
// The right iterator takes precedence over the result such that it
// determines the order and if there are repeats in the result.
func Subtract[T comparable](left, right collections.Iterator[T]) collections.Iterator[T] {
	inLeft := simpleSet.New[T]()
	return Where(right, func(value T) bool {
		if inLeft.Has(value) {
			return false
		}

		if left != nil {
			for left.Next() {
				cur := left.Current()
				inLeft.Set(cur)
				if value == cur {
					return false
				}
			}
			left = nil
		}

		return true
	})
}

// Zip merges two iterators together while both iterators have values
// and returns an iterator with a tuple containing values from both iterators.
func Zip[TFirst, TSecond, TOut any](
	firsts collections.Iterator[TFirst],
	seconds collections.Iterator[TSecond],
	combiner collections.Combiner[TFirst, TSecond, TOut],
) collections.Iterator[TOut] {
	zipping := true
	return New(func() (TOut, bool) {
		if zipping && firsts.Next() && seconds.Next() {
			return combiner(firsts.Current(), seconds.Current()), true
		}
		zipping = false
		return utils.Zero[TOut](), false
	})
}

// Interweave will pull values from each iterator, one at a time,
// and return them in the cycling order as an iterator.
// When an iterator runs out the remaining will interweave until all are empty.
func Interweave[T any](its []collections.Iterator[T]) collections.Iterator[T] {
	index, count := -1, len(its)
	return New(func() (T, bool) {
		for i := count - 1; i >= 0; i-- {
			index++
			if index >= count {
				index = 0
			}
			if it := its[index]; it != nil {
				if it.Next() {
					return it.Current(), true
				}
				its[index] = nil
			}
		}
		return utils.Zero[T](), false
	})
}

// SortInterweave creates an iterator that is the two given iterators interwoven
// such that both lists keep their order but lowest value from each list is used first.
// If the two iterators are sorted, this will effectively merge sort the values.
// This can take an optional comparer to override the default or
// if this type doesn't have a default comparer.
func SortInterweave[T any](left, right collections.Iterator[T], comparer ...comp.Comparer[T]) collections.Iterator[T] {
	cmp := optional.Comparer(comparer)
	hasLeft, hasRight := false, false
	var leftValue, rightValue T
	return New(func() (T, bool) {
		if !hasLeft && left.Next() {
			hasLeft = true
			leftValue = left.Current()
		}

		if !hasRight && right.Next() {
			hasRight = true
			rightValue = right.Current()
		}

		if hasLeft {
			if hasRight && cmp(leftValue, rightValue) > 0 {
				hasRight = false
				return rightValue, true
			}
			hasLeft = false
			return leftValue, true
		}

		if hasRight {
			hasRight = false
			return rightValue, true
		}

		return utils.Zero[T](), false
	})
}

// Sort iterates the values from the given iterator in sorted order.
// The values are sorted by the given comparer function or the default comparer.
// This can take an optional comparer to override the default or
// if this type doesn't have a default comparer.
func Sort[T any](it collections.Iterator[T], comparer ...comp.Comparer[T]) collections.Iterator[T] {
	cmp := optional.Comparer(comparer)
	var index, count int
	var sortedValues []T
	first := true
	return New(func() (T, bool) {
		if first {
			first = false
			sortedValues = ToSlice(it)
			slices.SortFunc(sortedValues, cmp)
			index = -1
			count = len(sortedValues)
		}

		if count > 0 {
			count--
			index++
			return sortedValues[index], true
		}
		return utils.Zero[T](), false
	})
}

// Sorted returns true if the given values in the iterator are sorted
// based on the given comparer or the default comparer.
// This can take an optional comparer to override the default or
// if this type doesn't have a default comparer.
func Sorted[T any](it collections.Iterator[T], comparer ...comp.Comparer[T]) bool {
	cmp := optional.Comparer(comparer)
	if it.Next() {
		prev := it.Current()
		for it.Next() {
			cur := it.Current()
			if cmp(prev, cur) > 0 {
				return false
			}
			prev = cur
		}
	}
	return true
}

// Indexed returns an iterator that returns a tuple containing the value
// from the given iterator and an index of the value starting with zero.
func Indexed[T any](it collections.Iterator[T]) collections.Iterator[collections.Tuple2[int, T]] {
	index := -1
	return New[collections.Tuple2[int, T]](func() (collections.Tuple2[int, T], bool) {
		if it.Next() {
			index++
			return tuple2.New(index, it.Current()), true
		}
		return utils.Zero[collections.Tuple2[int, T]](), false
	})
}
