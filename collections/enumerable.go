package collections

// Enumerable is an interface for an object which can have its data enumerated.
type Enumerable[T any] interface {

	// Enumerate gets an enumerator for this objects data.
	Enumerate() Enumerator[T]
}
