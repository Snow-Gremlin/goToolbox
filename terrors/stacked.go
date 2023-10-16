package terrors

// Stacked is an object, typically an error, which contains a
// stack trace for where this object was created.
type Stacked interface {

	// Gets the stack trace for where this object was created.
	Stack() string
}
