package readonlyQueue

import "goToolbox/collections"

// New wraps another queue in a readonly shell.
func New[T any](q collections.ReadonlyQueue[T]) collections.ReadonlyQueue[T] {
	return &readonlyQueueImp[T]{q: q}
}
