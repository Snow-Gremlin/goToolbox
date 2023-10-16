package collections

// ReadonlySet is a readonly version of a set.
//
// For sets, the `ToSlice`, `ToList`, and `Enumerate` methods do not guarantee
// any specific order and must be considered returning values in random order.
type ReadonlySet[T comparable] interface {
	Collection[T]
	Sliceable[T]
	Listable[T]
	Container[T]
}
