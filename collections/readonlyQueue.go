package collections

// ReadonlyQueue is the readonly version of a queue.
type ReadonlyQueue[T any] interface {
	Collection[T]
	Sliceable[T]
	Listable[T]
	Peeker[T]
}
