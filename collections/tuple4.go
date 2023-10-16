package collections

// Triple4 is an object containing four values.
type Tuple4[T1, T2, T3, T4 any] interface {
	Tuple

	// Value1 gets the first value.
	Value1() T1

	// Value2 gets the second value.
	Value2() T2

	// Value3 gets the third value.
	Value3() T3

	// Value4 gets the fourth value.
	Value4() T4

	// Values gets all the values in the tuple.
	Values() (T1, T2, T3, T4)
}
