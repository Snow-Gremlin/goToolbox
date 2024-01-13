package terrors

// MonoWrap is an error which wraps a single error.
type MonoWrap interface {
	// Unwrap unwraps the single error.
	Unwrap() error
}
