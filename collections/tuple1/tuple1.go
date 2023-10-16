package tuple1

import "goToolbox/collections"

// New constructs a new tuple with one value.
func New[T1 any](value1 T1) collections.Tuple1[T1] {
	return tuple1Imp[T1]{
		value1: value1,
	}
}
