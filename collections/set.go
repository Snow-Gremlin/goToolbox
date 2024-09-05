package collections

// Set is a collection of values in random order which has no repeat values.
//
// For sets, the `ToSlice`, `ToList`, and `Enumerate` methods do not guarantee
// any specific order and must be considered returning values in random order.
type Set[T any] interface {
	ReadonlySet[T]

	// Add inserts the given values into the set.
	// Returns true if any value was added, false if all values already existed.
	Add(values ...T) bool

	// AddFrom inserts the values from the given enumerator into the set.
	// Returns true if any value was added, false if all values already existed.
	AddFrom(e Enumerator[T]) bool

	// TakeAny removes one value from the set and returns it.
	// The set is in random order so this will be a random value.
	// If the set is empty, this will panic.
	TakeAny() T

	// TakeMany removes the given number of values from the set.
	// The values be in random order.
	// It will return less values if the set is shorter than the count.
	TakeMany(count int) []T

	// Remove removes all the given values from the set.
	// Returns true if any values were removed.
	Remove(values ...T) bool

	// RemoveIf removes all the values which return true for the given predicate.
	// Returns true if any values were removed.
	RemoveIf(handle Predicate[T]) bool

	// Clear removes all the values from the set.
	Clear()

	// Clones this set.
	Clone() Set[T]

	// Readonly gets a readonly version of this set.
	//
	// The readonly version points back to this set
	// but is not able to be cast into this set.
	Readonly() ReadonlySet[T]
}
