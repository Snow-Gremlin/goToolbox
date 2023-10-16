package set

import (
	"goToolbox/collections"
	"goToolbox/internal/optional"
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
		m: make(map[T]struct{}, optional.Capacity(capacity)),
	}
}

func With[T comparable](values ...T) collections.Set[T] {
	s := &setImp[T]{
		m: make(map[T]struct{}),
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
