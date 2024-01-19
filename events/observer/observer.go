package observer

import "github.com/Snow-Gremlin/goToolbox/events"

// New creates a new function observer.
//
// This calls a function in the same thread as the event invoke.
func New[T any](handle func(T)) events.Observer[T] {
	return &observerImp[T]{
		handle: handle,
	}
}
