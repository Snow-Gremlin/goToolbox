package events

// Unjoinable is an observer which gets notified when unjoined from an event.
type Unjoinable[T any] interface {
	// Unjoined is called by the given event when the observation
	// has been removed from the event.
	Unjoined(event Event[T])
}
