package testers

// Check performs a test of an expectation on some value.
type Check[T any] interface {

	// With adds context to the check being performed.
	// The new check with the additional context is returned.
	With(key string, args ...any) Check[T]

	// Withf adds formatted context to the check being performed.
	// The new check with the additional context is returned.
	Withf(key, format string, args ...any) Check[T]

	// WithType adds the type of the given value as context
	// to the check being performed.
	// The new check with the additional context is returned.
	WithType(key string, valueForType any) Check[T]

	// WithValue adds the given value as context to the check being performed.
	// The value is formatted based on test configuration.
	// The new check with the additional context is returned.
	WithValue(key string, value any) Check[T]

	// Name adds context with the key "name" to the check being performed.
	// This is a quick way to add unique information to a check.
	// The new check with the additional context is returned.
	Name(name string) Check[T]

	// AsText indicates that chars (byte) and runes (uint32)
	// will be formatted as unicode instead of integers.
	// Text formatting will take precedence over integer formats
	// for chars (byte) and runes (uint32).
	AsText() Check[T]

	// AsHex indicates that integers will be formatted
	// as hexadecimal numbers when outputted. Hexadecimal numbers
	// will be formatted in groups of four nibbles in code form,
	// e.g. `0xFFFF_FFFF_FFFF_FFFF`.
	AsHex() Check[T]

	// AsOct indicates that integers will be formatted
	// as octal numbers when outputted. Octal numbers
	// will be formatted in groups of four octets in code form,
	// e.g. `0o7777_7777_7777_7777`.
	AsOct() Check[T]

	// AsBin indicates that integers will be formatted
	// as binary numbers when outputted. Binary numbers
	// will be formatted in groups of four bits in code form,
	// e.g. `0b0000_1111_0000_1111_0000_1111_0000_1111`.
	AsBin() Check[T]

	// TimeAs indicates the string to format time with.
	// This will not affect the way time is compared, only how time
	// is printed on error.
	TimeAs(timeFormat string) Check[T]

	// Required sets all following asserts as required
	// such that if a check fails then the test is halted.
	Required() Check[T]

	// Assert performs this check on the given actual value.
	Assert(actual T) Check[T]

	// Require performs this check on the given actual value.
	//
	// If the check fails then the test is halted.
	// This is a short-cut for `Required().Assert(actual)`.
	Require(actual T) Check[T]

	// AssertAll performs this check on all the given
	// actual values in the given collection.
	//
	// The collection is read must have every element
	// the current expected type parameter.
	// If a map is given to either the expected or actual values.
	// The values being matched will be key/value tuples.
	AssertAll(actual any) Check[T]

	// RequireAll performs this check on all the given
	// actual values in the given enumeration.
	//
	// If the check fails then the test is halted.
	// This is a short-cut for `Required().AssertAll(actual)`.
	RequireAll(actual any) Check[T]

	// Panic performs this check on the value panicked from the given function.
	// This will fail if the function doesn't panic or the correct type wasn't recovered.
	// If the test fails an error is printed but the test will continue, unless required.
	Panic(handle func()) Check[T]
}
