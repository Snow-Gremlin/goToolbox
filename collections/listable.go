package collections

// Listable is an object which can get the data as a List.
type Listable[T any] interface {
	// ToList returns the values as a list.
	//
	// Typically the list will be an array or sliced
	// based list but may, for some cases, be a linked list.
	ToList() List[T]
}
