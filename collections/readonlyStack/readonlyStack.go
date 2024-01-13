package readonlyStack

import "github.com/Snow-Gremlin/goToolbox/collections"

// New wraps another stack in a readonly shell.
func New[T any](s collections.ReadonlyStack[T]) collections.ReadonlyStack[T] {
	return &readonlyStackImp[T]{s: s}
}
