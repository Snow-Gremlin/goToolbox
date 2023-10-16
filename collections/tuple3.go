package collections

// Triple3 is an object containing three values.
type Tuple3[T1, T2, T3 any] interface {
	Tuple

	// Value1 gets the first value.
	Value1() T1

	// Value2 gets the second value.
	Value2() T2

	// Value3 gets the third value.
	Value3() T3

	// Values gets all the values in the tuple.
	Values() (T1, T2, T3)
}
