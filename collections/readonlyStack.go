package collections

// ReadonlyStack is the readonly version of a stack.
type ReadonlyStack[T any] interface {
	Collection[T]
	Sliceable[T]
	Listable[T]
	Peeker[T]
}
