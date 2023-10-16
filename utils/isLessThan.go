package utils

// IsLessThan returns true if the given x is less than the given y, otherwise false.
type IsLessThan[T any] func(x, y T) bool
