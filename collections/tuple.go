package collections

import "goToolbox/utils"

// Tuple is an object containing several values.
type Tuple interface {
	Countable
	Getter[int, any]
	Sliceable[any]
	utils.Stringer
	utils.Equatable
}
