package collections

// Iterable is a function which constructs a new instance of an iterator.
type Iterable[T any] func() Iterator[T]
