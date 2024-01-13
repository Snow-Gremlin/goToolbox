package terrors

// Messager is an object, typically an error, which contains a message.
// For errors this message is the raw (constant) string message
// for the error with is used in the more detailed `Error() string`.
type Messager interface {
	// Message gets the raw message for this object.
	Message() string
}
