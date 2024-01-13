package collections

import "github.com/Snow-Gremlin/goToolbox/utils"

// Enumerator is a tool for walking through a collection of data.
type Enumerator[T any] interface {
	Countable
	Sliceable[T]
	utils.Equatable

	// Iterate creates a new iterator
	Iterate() Iterator[T]

	// Where filters the enumeration to only values which satisfy the predicate.
	Where(p Predicate[T]) Enumerator[T]

	// WhereNot filters the enumeration to only values which do not satisfy the predicate.
	WhereNot(p Predicate[T]) Enumerator[T]

	// NotNil filters the enumeration to only values which are not nil.
	NotNil() Enumerator[T]

	// NotZero filters the enumeration to only values which are not zero.
	NotZero() Enumerator[T]

	// Any determines if any value in the enumerator satisfies the given predicate.
	Any(p Predicate[T]) bool

	// All determines if all of the values in the enumerator satisfies the given predicate.
	// This will return true if the enumerator is empty.
	All(p Predicate[T]) bool

	// Foreach runs the given function for each values from the given enumerator.
	//
	// Use an `All` if your foreach needs to escape early,
	Foreach(m func(value T))

	// DoUntilError runs the given function for each value from the given enumerator.
	// If any error occurs, the error will be returned right away,
	// if no error occurs, then nil is returned.
	DoUntilError(m func(value T) error) error

	// StepsUntil determines the number of values in the enumerator are read until
	// a value satisfies the given predicate.
	//
	// The count will not include the value which satisfied the predicate such that
	// if the first value satisfies the predicate then this will return zero.
	// If no value satisfies the predicate then -1 is returned.
	StepsUntil(p Predicate[T]) int

	// Empty determines if the enumerator has no values.
	Empty() bool

	// AtLeast determines if there are at least the given number of values.
	// This is faster than using count so use when an exact count isn't needed.
	AtLeast(min int) bool

	// AtMost determines if there are at most the given number of values.
	// This is faster than using count so use when an exact count isn't needed.
	AtMost(max int) bool

	// First returns the first value in the enumerator with true,
	// or zero value with false if the enumerator is empty.
	First() (T, bool)

	// Last returns the last value in the enumerator with true,
	// or zero value with false if the enumerator is empty.
	Last() (T, bool)

	// Single returns the only value if there is only one value.
	// If there are no values or more than one then zero and false is returned.
	Single() (T, bool)

	// Skip skips over the given count of values before returning the rest.
	Skip(count int) Enumerator[T]

	// SkipWhile skips over values until the given predicate returns false.
	// The values from the given enumerator are returned after and including
	// the first false from the predicate.
	SkipWhile(p Predicate[T]) Enumerator[T]

	// Take enumerates the given number of values before stopping enumeration.
	Take(count int) Enumerator[T]

	// TakeWhile enumerates the values until the given predicate returns false.
	// The values are returned until and excluding the first false from the predicate.
	TakeWhile(p Predicate[T]) Enumerator[T]

	// Replace replaces or returns values using the given replacer function.
	Replace(replacer Selector[T, T]) Enumerator[T]

	// Reverse enumerates the values in the opposite order.
	Reverse() Enumerator[T]

	// Strings enumerates the string of each value.
	Strings() Enumerator[string]

	// Quotes enumerates the string of each value quoted.
	Quotes() Enumerator[string]

	// Trim converts all the values into strings and trims any whitespace
	// from the front and end of the strings.
	Trim() Enumerator[string]

	// Join converts all the values into strings and joins
	// them with the given separator.
	Join(sep string) string

	// Append enumerates all the values in the enumerator
	// followed by the given tail value.
	Append(tails ...T) Enumerator[T]

	// Concat enumerates all the values in the enumerator
	// followed by the values enumerators from the given tail.
	Concat(tails ...Enumerator[T]) Enumerator[T]

	// SortInterweave creates an enumerator that is the two given enumerators interwoven
	// such that both lists keep their order but lowest value from each list is used first.
	//
	// If the two enumerators are sorted, this will effectively merge sort the values.
	// This can take an optional comparer to override the default comparer
	// or to give a comparer if there is no default comparer for this type.
	SortInterweave(other Enumerator[T], comparer ...utils.Comparer[T]) Enumerator[T]

	// Sorted determines if the values in the enumerator are already sorted.
	//
	// This can take an optional comparer to override the default comparer
	// or to give a comparer if there is no default comparer for this type.
	Sorted(comparer ...utils.Comparer[T]) bool

	// Sort enumerates the values in sorted order.
	//
	// This can take an optional comparer to override the default comparer
	// or to give a comparer if there is no default comparer for this type.
	Sort(comparer ...utils.Comparer[T]) Enumerator[T]

	// Merge preforms a merge of the values in the given enumerator.
	//
	// The merge method is called with the prior returned value from the previous call.
	// The first value is used as the prior value with the second value in the merger.
	// The last returned value from merge is returned, the first value if there
	// is only one value, or the zero value if no values.
	Merge(merger Reducer[T, T]) T

	// Max gets the maximum value from all the values.
	//
	// This can take an optional comparer to override the default comparer
	// or to give a comparer if there is no default comparer for this type.
	Max(comparer ...utils.Comparer[T]) T

	// Min gets the minimum value from all the values.
	//
	// This can take an optional comparer to override the default comparer
	// or to give a comparer if there is no default comparer for this type.
	Min(comparer ...utils.Comparer[T]) T

	// Buffered stores the result of an enumeration and repeats it back
	// in the returned enumerator. Uses memory to reduce calculations.
	// This will not read from this enumerator only when needed.
	Buffered() Enumerator[T]

	// StartsWith determines if the first enumerator starts with the given prefix.
	StartsWith(other Enumerator[T]) bool
}
