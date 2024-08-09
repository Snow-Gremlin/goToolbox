package readonlySortedSet

import "github.com/Snow-Gremlin/goToolbox/collections"

// New wraps another sorted set in a readonly shell.
func New[T any](s collections.ReadonlySortedSet[T]) collections.ReadonlySortedSet[T] {
	return readonlySortedSetImp[T]{s: s}
}
