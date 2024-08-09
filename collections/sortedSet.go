package collections

// Set is a collection of values in random order which has no repeat values.
//
// For sets, the `ToSlice`, `ToList`, and `Enumerate` methods do not guarantee
// any specific order and must be considered returning values in random order.
type SortedSet[T any] interface {
	ReadonlySortedSet[T]

	// Add inserts the given values into the set.
	// If any value already exists, it will not be replaced.
	// Returns true if any value was added, false if all values already existed.
	Add(values ...T) bool

	// Overwrite inserts the given values into the set.
	// If any value already exists, it will be replaced with the new instance.
	// Returns true if any value was added, false if all values already existed.
	Overwrite(values ...T) bool

	// AddFrom inserts the values from the given enumerator into the set.
	// If any value already exists, it will not be replaced.
	// Returns true if any value was added, false if all values already existed.
	AddFrom(e Enumerator[T]) bool

	// OverwriteFrom inserts the values from the given enumerator into the set.
	// If any value already exists, it will be replaced with the new instance.
	// Returns true if any value was added, false if all values already existed.
	OverwriteFrom(e Enumerator[T]) bool

	// TryAdd inserts the given value into the set and returns the inserted
	// value with true, unless the value already exists. If the given value
	// already exists then the existing value is returned with false.
	TryAdd(value T) (T, bool)

	// TakeFirst removes one value from the front of the list.
	// If the list is empty, this will panic.
	TakeFirst() T

	// TakeLast removes one value from the back of the list.
	// If the list is empty, this will panic.
	TakeLast() T

	// TakeFront remove the given number of values from the front of list.
	// It will return less values if the list is shorter than the count.
	TakeFront(count int) List[T]

	// TakeBack remove the given number of values from the back of list.
	// It will return less values if the list is shorter than the count.
	TakeBack(count int) List[T]

	// Remove removes all the given values from the set.
	// Returns true if any values were removed.
	Remove(values ...T) bool

	// RemoveIf removes all the values which return true for the given predicate.
	// Returns true if any values were removed.
	RemoveIf(handle Predicate[T]) bool

	// RemoveRange removes the given number of values from the given index.
	RemoveRange(index, count int)

	// Clear removes all the values from the set.
	Clear()

	// Clones this set.
	Clone() SortedSet[T]

	// Readonly gets a readonly version of this set.
	//
	// The readonly version points back to this set
	// but is not able to be cast into this set.
	Readonly() ReadonlySortedSet[T]
}
