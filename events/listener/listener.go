package listener

import "github.com/Snow-Gremlin/goToolbox/events"

// New creates a new listener that calls the given handle for when
// any of the subscribed events are invoked.
func New[T any](handle func(value T)) events.Listener[T] {
	return listenerImp[T]{
		obv: &listenerObserver[T]{
			handle: handle,
			events: []events.Event[T]{},
		},
	}
}
