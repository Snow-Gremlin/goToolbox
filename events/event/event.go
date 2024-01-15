package event

import "github.com/Snow-Gremlin/goToolbox/events"

// New creates a new event instance.
func New[T any]() events.Event[T] {
	return &eventImp[T]{obs: nil}
}

// Empty creates a new event instance which will not error
// but will not allow observers to be added nor invoked.
//
// This should be used when an event is not available so that attempting
// to add an observer to an event will not require the event assignment
// to first check if the event is nil or not.
func Empty[T any]() events.Event[T] {
	return (*eventImp[T])(nil)
}
