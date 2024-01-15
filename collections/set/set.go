package set

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/internal/optional"
	"github.com/Snow-Gremlin/goToolbox/internal/simpleSet"
)

// New creates a new set with unsorted values.
//
// The values will be returned in random order when enumeration
// and may have different orders per enumeration.
//
// If one capacity value is given, an empty underlying map is allocated
// with enough space to hold the specified number of elements.
// The capacity may be omitted, in which case a small starting size is allocated.
func New[T comparable](capacity ...int) collections.Set[T] {
	return &setImp[T]{
		m:     simpleSet.Cap[T](optional.Capacity(capacity)),
		event: nil,
	}
}

// With creates a new set initialized with the given values.
//
// The values will be returned in random order when enumeration
// and may have different orders per enumeration.
func With[T comparable](values ...T) collections.Set[T] {
	s := &setImp[T]{
		m:     simpleSet.New[T](),
		event: nil,
	}
	s.Add(values...)
	return s
}

// From creates a new dictionary with unsorted keys
// populated with key/value pairs from the given tuple enumerator.
func From[T comparable](e collections.Enumerator[T], capacity ...int) collections.Set[T] {
	d := New[T](capacity...)
	d.AddFrom(e)
	return d
}
