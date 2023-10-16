package collections

// Predicate is a function that takes a value and returns a boolean.
//
// Typically predicates are used to answer some question like "is equal" or "contains".
type Predicate[T any] func(value T) bool
