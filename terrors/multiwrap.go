package terrors

// MultiWrap is an error which wraps a single error.
type MultiWrap interface {
	// Unwraps all of the internal errors.
	Unwrap() []error
}
