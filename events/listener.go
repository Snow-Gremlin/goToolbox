package events

// Listener is a specialized observer that keeps track of the events
// that it is subscribed to, so that it can be cancelled.
//
// When cancelled, this listener will remove itself from all the events it is
// subscribed to. If an event clears out it's observations, this listener will
// automatically unsubscribe from that event.
//
// Use with caution since this could hold onto an event while the event holds
// onto this listener, meaning neither will be garbage collect until they
// both can be. Be sure to unsubscribe, cancel, or clear the event when done.
type Listener[T any] interface {
	Cancelable

	// Subscribe adds this listener to the given event.
	Subscribe(event Event[T]) bool

	// Unsubscribe removes this listener from the given event.
	Unsubscribe(event Event[T]) bool
}
