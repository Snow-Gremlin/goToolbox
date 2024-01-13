package capStack

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/internal/optional"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

// New creates a new stack.
//
// This may optionally have an initial size to
// pre-populate the stack with that number of zero values.
func New[T any](sizes ...int) collections.Stack[T] {
	size, initCap := optional.SizeAndCapacity(sizes)
	s := newImp[T]()
	s.growCap(initCap)
	s.PushFrom(enumerator.Repeat(utils.Zero[T](), size))
	return s
}

// Fill creates a new stack filled with the given
// value repeated the given number of times.
func Fill[T any](value T, count int, capacity ...int) collections.Stack[T] {
	count = max(count, 0)
	initCap := max(count, optional.Capacity(capacity))
	s := newImp[T]()
	s.growCap(initCap)
	s.PushFrom(enumerator.Repeat(value, count))
	return s
}

// With creates a stack with the given values.
func With[T any](values ...T) collections.Stack[T] {
	s := newImp[T]()
	s.Push(values...)
	return s
}

// From creates a new stack from the given enumerator.
func From[T any](e collections.Enumerator[T], capacity ...int) collections.Stack[T] {
	initCap := optional.Capacity(capacity)
	s := newImp[T]()
	s.growCap(initCap)
	s.PushFrom(e)
	return s
}
