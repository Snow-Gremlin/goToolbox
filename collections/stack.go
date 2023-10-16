package collections

// Stack is a linear collection of values which are FILO (first in, last out).
type Stack[T any] interface {
	ReadonlyStack[T]
	Clippable

	// Push adds all the given values onto
	// the stack in the order that they were given in.
	Push(values ...T)

	// PushFrom adds all the values from the given enumerator
	// onto the stack in the order that they were given in.
	PushFrom(e Enumerator[T])

	// Take pops the given number of values from the stack.
	// It will return less values if the stack is shorter than the count.
	Take(count int) []T

	// Pop removes and returns one value from the stack.
	// If there are no values in the stack, this will panic.
	Pop() T

	// TryPop removes and returns one value from the stack.
	// Returns zero and false if there are no values in the stack.
	TryPop() (T, bool)

	// TrimTo will remove anything off the back of the stack
	// until the stack is the the given count.
	TrimTo(count int)

	// Clear removes all the values from the stack.
	Clear()

	// Clone makes a copy of this stack.
	Clone() Stack[T]

	// Readonly gets a readonly version of this stack that will stay up-to-date
	// with this stack but will not allow changes itself.
	Readonly() ReadonlyStack[T]
}
