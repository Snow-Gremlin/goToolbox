package queue

import (
	"goToolbox/collections"
	"goToolbox/collections/enumerator"
	"goToolbox/internal/optional"
	"goToolbox/utils"
)

// New creates a new linked queue.
//
// This may optionally have an initial size to
// pre-populate the queue with that number of zero values.
func New[T any](size ...int) collections.Queue[T] {
	return Fill(utils.Zero[T](), optional.Size(size))
}

// Fill creates a new linked queue filled with the given
// value repeated the given number of times.
func Fill[T any](value T, count int) collections.Queue[T] {
	return From(enumerator.Repeat(value, max(count, 0)))
}

// With creates a queue with the given values.
func With[T any](values ...T) collections.Queue[T] {
	q := &queueImp[T]{}
	q.Enqueue(values...)
	return q
}

// From creates a new queue from the given enumerator.
func From[T any](e collections.Enumerator[T]) collections.Queue[T] {
	q := &queueImp[T]{}
	q.EnqueueFrom(e)
	return q
}
