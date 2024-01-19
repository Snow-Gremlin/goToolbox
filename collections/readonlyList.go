package collections

// ReadonlyList is a readonly linear collection of values.
type ReadonlyList[T any] interface {
	Collection[T]
	Sliceable[T]
	Container[T]
	Getter[int, T]
	OnChanger

	// First attempts to get the first value from the list.
	// If the list is empty, this will panic.
	First() T

	// Last attempts to get the last value from the list.
	// If the list is empty, this will panic.
	Last() T

	// Backwards gets an enumerator for this list that
	// goes from the end to the front.
	Backwards() Enumerator[T]

	// IndexOf gets the index of the given value type,
	// -1 is returned if the value is not in the list.
	//
	// May have one optional after index to start searching after.
	// If the after index is negative the search starts from the beginning.
	// If the after index is greater or equal to the length then -1 is returned.
	IndexOf(value T, after ...int) int

	// StartsWith determines if the given list of values
	// is at the start of this list.
	StartsWith(other ReadonlyList[T]) bool

	// EndsWith determines if the given list of values
	// is at the end of this list.
	EndsWith(other ReadonlyList[T]) bool
}
