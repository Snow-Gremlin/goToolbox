package tuple2

import "github.com/Snow-Gremlin/goToolbox/collections"

// New constructs a new tuple with two values.
func New[T1, T2 any](value1 T1, value2 T2) collections.Tuple2[T1, T2] {
	return tuple2Imp[T1, T2]{
		value1: value1,
		value2: value2,
	}
}
