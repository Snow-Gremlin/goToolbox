package collections

// Queue is a linear collection of values which are FIFO (first in, first out).
type Queue[T any] interface {
	ReadonlyQueue[T]
	Clippable

	// Enqueue adds all the given values into
	// the queue in the order that they were given in.
	Enqueue(values ...T)

	// EnqueueFrom adds all the values from the given enumerator
	// onto the queue in the order that they were given in.
	EnqueueFrom(e Enumerator[T])

	// Take dequeues the given number of values from the queue.
	// It will return less values if the queue is shorter than the count.
	Take(count int) []T

	// Dequeue removes and returns one value from the queue.
	// If there are no values in the queue, this will panic.
	Dequeue() T

	// TryDequeue removes and returns one value from the queue.
	// Returns zero and false if there are no values in the queue.
	TryDequeue() (T, bool)

	// Clear removes all the values from the queue.
	Clear()

	// Clone makes a copy of this queue.
	Clone() Queue[T]

	// Readonly gets a readonly version of this queue that will stay up-to-date
	// with this queue but will not allow changes itself.
	Readonly() ReadonlyQueue[T]
}
