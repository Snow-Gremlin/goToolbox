package collections

// Peeker is an object which can peek data from the location that data is usually taken from in the object.
type Peeker[T any] interface {

	// Peek peeks at the next value, without removing it.
	//
	// This will panic if there are no values that can be peeked.
	Peek() T

	// TryPeek peeks at the next value, without removing it.
	//
	// Returns zero and false if there are no values that can be peeked.
	TryPeek() (T, bool)
}
