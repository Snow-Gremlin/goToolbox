package collections

// ReadonlySet is a readonly version of a set.
type ReadonlySet[T comparable] interface {
	Collection[T]
	Sliceable[T]
	Listable[T]
	Container[T]
}
