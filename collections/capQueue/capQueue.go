package capQueue

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/internal/optional"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

// New creates a new linked queue with capacity.
//
// This may optionally have an initial size or capacity to
// pre-populate the queue with that number of zero values.
func New[T any](sizes ...int) collections.Queue[T] {
	size, initCap := optional.SizeAndCapacity(sizes)
	q := newImp[T]()
	q.growCap(initCap)
	q.EnqueueFrom(enumerator.Repeat(utils.Zero[T](), size))
	return q
}

// Fill creates a new linked queue filled with the given
// value repeated the given number of times.
// This may include an optional capacity.
// The capacity must be larger than the count to have any effect.
func Fill[T any](value T, count int, capacity ...int) collections.Queue[T] {
	count = max(count, 0)
	initCap := max(count, optional.Capacity(capacity))
	q := newImp[T]()
	q.growCap(initCap)
	q.EnqueueFrom(enumerator.Repeat(value, count))
	return q
}

// With creates a queue with the given values.
func With[T any](values ...T) collections.Queue[T] {
	q := newImp[T]()
	q.Enqueue(values...)
	return q
}

// From creates a new queue from the given enumerator.
func From[T any](e collections.Enumerator[T], capacity ...int) collections.Queue[T] {
	initCap := optional.Capacity(capacity)
	q := newImp[T]()
	q.growCap(initCap)
	q.EnqueueFrom(e)
	return q
}
