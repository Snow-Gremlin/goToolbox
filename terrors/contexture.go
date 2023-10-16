package terrors

// Contexture is an object, typically an error, which
// contains contextual information.
type Contexture interface {

	// Context gets a copy of the context for this object.
	Context() map[string]any
}
