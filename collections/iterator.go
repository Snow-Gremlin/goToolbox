package collections

// Iterator is an object to step through a set of values
// as part of an enumeration of values.
//
// Typical usage of an iterator is to use it in a while loop:
//
// ```Go
// var it Iterator[T] = //...
//
//	for it.Next() {
//	   it.Current()
//	}
//
// ```
type Iterator[T any] interface {
	// Next steps this iterator the next value and updates Current.
	//
	// After creation `Next` should be called to prime the iterator
	// to the first value in the set.
	Next() bool

	// Current value in the iterator.
	//
	// This will return the zero value after creation until
	// the first time `Next` is called.
	Current() T
}
