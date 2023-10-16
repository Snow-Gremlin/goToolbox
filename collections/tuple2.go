package collections

// Tuple2 is an object containing two values.
type Tuple2[T1, T2 any] interface {
	Tuple

	// Value1 gets the first value.
	Value1() T1

	// Value2 gets the second value.
	Value2() T2

	// Values gets all the values in the tuple.
	Values() (T1, T2)
}
