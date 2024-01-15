package events

// Event is a collection of observers which will invoke
// on some event occurrence.
type Event[T any] interface {
	// Add adds a observer to this event.
	Add(observer Observer[T]) bool

	// Remove removes the given observer.
	Remove(observer Observer[T]) bool

	// Clear removes all the observers from this event.
	Clear()

	// Invoke will use the given value to update
	// all the observers attached to this event.
	Invoke(value T)
}
