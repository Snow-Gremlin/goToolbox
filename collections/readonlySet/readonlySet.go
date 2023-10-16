package readonlySet

import "goToolbox/collections"

// New wraps another set in a readonly shell.
func New[T comparable](s collections.ReadonlySet[T]) collections.ReadonlySet[T] {
	return &readonlySetImp[T]{
		s: s,
	}
}
