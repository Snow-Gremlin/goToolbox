package stack

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
func New[T any](size ...int) collections.Stack[T] {
	return Fill(utils.Zero[T](), optional.Size(size))
}

// Fill creates a new stack filled with the given
// value repeated the given number of times.
func Fill[T any](value T, count int) collections.Stack[T] {
	return From(enumerator.Repeat(value, max(count, 0)))
}

// With creates a stack with the given values.
func With[T any](values ...T) collections.Stack[T] {
	s := newImp[T]()
	s.Push(values...)
	return s
}

// From creates a new stack from the given enumerator.
func From[T any](e collections.Enumerator[T]) collections.Stack[T] {
	s := newImp[T]()
	s.PushFrom(e)
	return s
}
