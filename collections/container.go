package collections

// Container is a data collection that can be queried
// to determine if a specific value is in the collection.
type Container[T any] interface {
	// Contains determines if the given value exists in the collection.
	// This method works as a predicate.
	Contains(value T) bool
}
