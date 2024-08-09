package sortedSet

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/comp"
	"github.com/Snow-Gremlin/goToolbox/internal/optional"
)

// New creates a new list on an underlying array
//
// The sizes specifies the optional initial length and capacity.
// If only one size if given, that is the initial size of this list
// as well as the initial capacity. If a second integer argument is
// provided it will specify a different capacity from the length.
// The capacity will never be smaller than the list's length.
func New[T any](comparer ...comp.Comparer[T]) collections.SortedSet[T] {
	return CapNew[T](0, comparer...)
}

// CapNew creates a new dictionary with sorted keys and initial capacity
// by the optional given comparer function or the default comparer.
func CapNew[T any](capacity int, comparer ...comp.Comparer[T]) collections.SortedSet[T] {
	cmp := optional.Comparer(comparer)
	capacity = max(capacity, 0)
	return &sortedSetImp[T]{
		data:     make([]T, 0, capacity),
		comparer: cmp,
		event:    nil,
	}
}

// With creates a new list with the given values.
func With[T any](s []T, comparer ...comp.Comparer[T]) collections.SortedSet[T] {
	return CapFrom[T](enumerator.Enumerate(s...), len(s), comparer...)
}

// From creates a new list from the given enumerator.
//
// This may have an optional capacity for the list's initial capacity.
// Giving it a capacity will help when the enumerator contains a lot of values.
func From[T any](e collections.Enumerator[T], comparer ...comp.Comparer[T]) collections.SortedSet[T] {
	return CapFrom(e, 0, comparer...)
}

// CapFrom creates a new dictionary with sorted keys and an initial capacity
// populated with key/value pairs from the given tuple enumerator.
//
// The keys are sorted with the optional given comparer function
// or the default comparer if no comparer was given.
func CapFrom[T any](e collections.Enumerator[T], capacity int, comparer ...comp.Comparer[T]) collections.SortedSet[T] {
	s := CapNew(capacity, comparer...)
	s.AddFrom(e)
	return s
}
