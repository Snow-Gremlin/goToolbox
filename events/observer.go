package events

// Observer is an object that can be added to an event to listen for that
// event to be invoked. On the event invocation this update will be called.
type Observer[T any] interface {
	// Update is called when an event this is listening to is invoked.
	Update(value T)
}
