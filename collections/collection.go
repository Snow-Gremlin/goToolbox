package collections

import (
	"github.com/Snow-Gremlin/goToolbox/comp"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

// Collection is a collection of values.
type Collection[T any] interface {
	Enumerable[T]
	Countable
	utils.Stringer
	comp.Equatable

	// Empty indicates if the collection is empty.
	Empty() bool
}
