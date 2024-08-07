package collections

import (
	"github.com/Snow-Gremlin/goToolbox/comp"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

// Tuple is an object containing several values.
type Tuple interface {
	Countable
	Getter[int, any]
	Sliceable[any]
	utils.Stringer
	comp.Equatable
}
