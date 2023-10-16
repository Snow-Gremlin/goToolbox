package capStack

import (
	"goToolbox/collections"
	"goToolbox/collections/enumerator"
	"goToolbox/internal/optional"
	"goToolbox/utils"
)

// New creates a new stack.
//
// This may optionally have an initial size to
// pre-populate the stack with that number of zero values.
func New[T any](sizes ...int) collections.Stack[T] {
	size, cap := optional.SizeAndCapacity(sizes)
	s := &capStackImp[T]{}
	s.growCap(cap)
	s.PushFrom(enumerator.Repeat(utils.Zero[T](), size))
	return s
}

// Fill creates a new stack filled with the given
// value repeated the given number of times.
func Fill[T any](value T, count int, capacity ...int) collections.Stack[T] {
	count = max(count, 0)
	cap := max(count, optional.Capacity(capacity))
	s := &capStackImp[T]{}
	s.growCap(cap)
	s.PushFrom(enumerator.Repeat(value, count))
	return s
}

// With creates a stack with the given values.
func With[T any](values ...T) collections.Stack[T] {
	s := &capStackImp[T]{}
	s.Push(values...)
	return s
}

// From creates a new stack from the given enumerator.
func From[T any](e collections.Enumerator[T], capacity ...int) collections.Stack[T] {
	cap := optional.Capacity(capacity)
	s := &capStackImp[T]{}
	s.growCap(cap)
	s.PushFrom(e)
	return s
}
