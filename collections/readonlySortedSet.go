package collections

// ReadonlySortedSet is a readonly version of a sorted set.
type ReadonlySortedSet[T any] interface {
	ReadonlySet[T]
	Getter[int, T]

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
	IndexOf(value T) int
}
