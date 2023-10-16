package testers

// Check performs a test of an expectation on some value.
type Check[T any] interface {

	// With adds context to the check being performed.
	// The new check with the additional context is returned.
	With(key string, args ...any) Check[T]

	// Withf adds formatted context to the check being performed.
	// The new check with the additional context is returned.
	Withf(key, format string, args ...any) Check[T]

	// Name adds context with the key "name" to the check being performed.
	// This is a quick way to add unique information to a check.
	// The new check with the additional context is returned.
	Name(name string) Check[T]

	// Required sets all following asserts as required
	// such that if a check fails then the test is halted.
	Required() Check[T]

	// Assert performs this check on the given actual value.
	Assert(actual T) Check[T]

	// Require performs this check on the given actual value.
	// If the check fails then the test is halted.
	// This is a short-cut for `Required().Assert(actual)`.
	Require(actual T) Check[T]

	// Panic performs this check on the value panicked from the given function.
	// This will fail if the function doesn't panic or the correct type wasn't recovered.
	// If the test fails an error is printed but the test will continue, unless required.
	Panic(handle func()) Check[T]
}
