package collections

// List is a linear collection of values.
type List[T any] interface {
	ReadonlyList[T]

	// Prepend adds a new values to the front of the list.
	// The values will end up in the list in the same order they are given.
	Prepend(values ...T)

	// PrependFrom adds a new values to the front of the list.
	// The values will end up in the list in the same order they are given.
	PrependFrom(e Enumerator[T])

	// Append adds a new  values to the back of the list.
	// The values will end up in the list in the same order they are given.
	Append(values ...T)

	// AppendFrom adds a new  values to the back of the list.
	// The values will end up in the list in the same order they are given.
	AppendFrom(e Enumerator[T])

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

	// Insert adds the given values into the list at the given location.
	// If the index is zero, then these values will be added to the front of the list.
	// If the index is the length of the list, then the values will be added to the back of the list.
	Insert(index int, values ...T)

	// InsertFrom adds the values from the given enumerator into the list at the given location.
	// If the index is zero, then these values will be added to the front of the list.
	// If the index is the length of the list, then the values will be added to the back of the list.
	InsertFrom(index int, e Enumerator[T])

	// Remove removes the given number of values from the given index.
	Remove(index, count int)

	// RemoveIf removes all the values which return true for the given predicate.
	RemoveIf(handle Predicate[T]) bool

	// Set sets the values starting with the given index.
	// The index must [0..count] to be valid.
	// If there are more values given than already exist in the list,
	// the remaining will be appended.
	Set(index int, values ...T)

	// SetFrom sets the values starting with the given index.
	// If there are more values given than already exist in the list,
	// the remaining will be appended.
	SetFrom(index int, e Enumerator[T])

	// Clear removes all the values from the whole list, leaving the list empty.
	Clear()

	// Clones this list.
	Clone() List[T]

	// Readonly gets a readonly version of this list that will stay up-to-date
	// with this list but will not allow changes itself.
	Readonly() ReadonlyList[T]
}
