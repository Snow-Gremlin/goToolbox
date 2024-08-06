package readonlySet

import "github.com/Snow-Gremlin/goToolbox/collections"

// New wraps another set in a readonly shell.
func New[T any](s collections.ReadonlySet[T]) collections.ReadonlySet[T] {
	return readonlySetImp[T]{s: s}
}
