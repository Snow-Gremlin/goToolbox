package linkedList

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/internal/optional"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

// New creates a new linked list.
//
// This may optionally have an initial size to
// pre-populate the list with that number of zero values.
func New[T any](size ...int) collections.List[T] {
	return Fill(utils.Zero[T](), optional.Size(size))
}

// Fill creates a new linked list filled with the given
// value repeated the given number of times.
func Fill[T any](value T, count int) collections.List[T] {
	return impFrom(enumerator.Repeat(value, max(count, 0)))
}

// With creates a new linked list with the given values.
func With[T any](s ...T) collections.List[T] {
	return newImp(s...)
}

// From creates a new linked list from the given enumerator.
func From[T any](e collections.Enumerator[T]) collections.List[T] {
	return impFrom(e)
}
