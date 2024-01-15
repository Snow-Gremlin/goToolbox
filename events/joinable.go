package events

// Joinable is an observer which gets notified when joined to an event.
type Joinable[T any] interface {
	// Joined is called by the given event when the observation
	// has been added to the event.
	Joined(event Event[T])
}
