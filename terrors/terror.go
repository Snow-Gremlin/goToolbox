package terrors

// Error is an error object for the toolbox.
type TError interface {
	error
	Stacked
	Contexture
	MultiWrap
	Messager

	// Equals returns true if this error and the given error are equal.
	//
	// This will not compare stacks.
	Equals(other any) bool

	// String gets a detailed string for this error.
	String() string

	// With adds more context to this error.
	// Returns the receiver for method chaining.
	With(key string, value any) TError

	// WithError adds another wrapped error to this error.
	// Returns the receiver for method chaining.
	WithError(err error) TError

	// Resets the stack trace for the error.
	// Offset is the number of frames to leave off of the top of the stack.
	// Returns the receiver for method chaining.
	ResetStack(offset int) TError

	// Clone creates a copy of this error.
	Clone() TError
}
