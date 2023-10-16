package collections

// Tuple1 is an object containing one value.
type Tuple1[T1 any] interface {
	Tuple

	// Value1 gets the value.
	Value1() T1
}
