package tuple4

import "github.com/Snow-Gremlin/goToolbox/collections"

// New constructs a new tuple with four values.
func New[T1, T2, T3, T4 any](value1 T1, value2 T2, value3 T3, value4 T4) collections.Tuple4[T1, T2, T3, T4] {
	return tuple4Imp[T1, T2, T3, T4]{
		value1: value1,
		value2: value2,
		value3: value3,
		value4: value4,
	}
}
