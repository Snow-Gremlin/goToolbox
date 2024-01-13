package tuple3

import "github.com/Snow-Gremlin/goToolbox/collections"

// New constructs a new tuple with three values.
func New[T1, T2, T3 any](value1 T1, value2 T2, value3 T3) collections.Tuple3[T1, T2, T3] {
	return tuple3Imp[T1, T2, T3]{
		value1: value1,
		value2: value2,
		value3: value3,
	}
}
