package list

import (
	"slices"

	"goToolbox/collections"
	"goToolbox/internal/optional"
	"goToolbox/utils"
)

// New creates a new list on an underlying array
//
// The sizes specifies the optional initial length and capacity.
// If only one size if given, that is the initial size of this list
// as well as the initial capacity. If a second integer argument is
// provided it will specify a different capacity from the length.
// The capacity will never be smaller than the list's length.
func New[T any](sizes ...int) collections.List[T] {
	size, cap := optional.SizeAndCapacity(sizes)
	return &listImp[T]{
		s: make([]T, size, cap),
	}
}

// Fill creates a new list initialized with a repeated value.
//
// This may have an optional capacity for the list's initial capacity.
func Fill[T any](value T, count int, capacity ...int) collections.List[T] {
	count = max(count, 0)
	cap := max(count, optional.Capacity(capacity))
	list := &listImp[T]{
		s: make([]T, count, cap),
	}
	for i := 0; i < count; i++ {
		list.s[i] = value
	}
	return list
}

// With creates a new list with the given values.
func With[T any](s ...T) collections.List[T] {
	return &listImp[T]{
		s: slices.Clone(s),
	}
}

// From creates a new list from the given enumerator.
//
// This may have an optional capacity for the list's initial capacity.
// Giving it a capacity will help when the enumerator contains a lot of values.
func From[T any](e collections.Enumerator[T], capacity ...int) collections.List[T] {
	s := Fill(utils.Zero[T](), 0, capacity...)
	s.AppendFrom(e)
	return s
}
