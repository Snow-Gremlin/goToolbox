package capQueue

import (
	"goToolbox/collections"
	"goToolbox/collections/enumerator"
	"goToolbox/internal/optional"
	"goToolbox/utils"
)

// New creates a new linked queue with capacity.
//
// This may optionally have an initial size or capacity to
// pre-populate the queue with that number of zero values.
func New[T any](sizes ...int) collections.Queue[T] {
	size, cap := optional.SizeAndCapacity(sizes)
	q := &capQueueImp[T]{}
	q.growCap(cap)
	q.EnqueueFrom(enumerator.Repeat(utils.Zero[T](), size))
	return q
}

// Fill creates a new linked queue filled with the given
// value repeated the given number of times.
// This may include an optional capacity.
// The capacity must be larger than the count to have any effect.
func Fill[T any](value T, count int, capacity ...int) collections.Queue[T] {
	count = max(count, 0)
	cap := max(count, optional.Capacity(capacity))
	q := &capQueueImp[T]{}
	q.growCap(cap)
	q.EnqueueFrom(enumerator.Repeat(value, count))
	return q
}

// With creates a queue with the given values.
func With[T any](values ...T) collections.Queue[T] {
	q := &capQueueImp[T]{}
	q.Enqueue(values...)
	return q
}

// From creates a new queue from the given enumerator.
func From[T any](e collections.Enumerator[T], capacity ...int) collections.Queue[T] {
	cap := optional.Capacity(capacity)
	q := &capQueueImp[T]{}
	q.growCap(cap)
	q.EnqueueFrom(e)
	return q
}
