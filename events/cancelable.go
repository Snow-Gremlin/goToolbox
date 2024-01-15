package events

// Cancelable is an object which can be cancelled.
type Cancelable interface {
	// Cancel will cancel this object.s
	Cancel()
}
