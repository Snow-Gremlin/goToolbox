package simpleSet

import (
	"maps"
	"sort"
	"strings"

	"github.com/Snow-Gremlin/goToolbox/internal/liteUtils"
)

var slug = struct{}{}

// Set is a simple set using an underlying map
type Set[T comparable] map[T]struct{}

// New creates a simple set without capacity.
func New[T comparable]() Set[T] {
	return make(Set[T])
}

// Cap creates a simple set with capacity.
func Cap[T comparable](capacity int) Set[T] {
	return make(Set[T], capacity)
}

// With creates a simple set with the given values as the keys in the map.
func With[T comparable](values ...T) Set[T] {
	s := Cap[T](len(values))
	for _, key := range values {
		s.Set(key)
	}
	return s
}

// Has determines if the given value is in the set.
func (s Set[T]) Has(value T) bool {
	_, has := s[value]
	return has
}

// Count gets the number of values in the set.
func (s Set[T]) Count() int {
	return len(s)
}

// ToSlice gets all the values in the set in random order.
func (s Set[T]) ToSlice() []T {
	return liteUtils.Keys(s)
}

// Clone creates a shallow clone of the set.
// This will not deep copying the values.
func (s Set[T]) Clone() Set[T] {
	return maps.Clone(s)
}

// Set the given value in the set.
// Has no effect if the value is already set.
func (s Set[T]) Set(value T) {
	s[value] = slug
}

// SetTest checks if the value exists or not before being set.
// Returns true if value is new, otherwise false if already set.
func (s Set[T]) SetTest(value T) bool {
	if _, has := s[value]; has {
		return false
	}
	s[value] = slug
	return true
}

// Remove removes the given value from the set.
// Has no effect if the value is not already set.
func (s Set[T]) Remove(value T) {
	delete(s, value)
}

// RemoveTest removes the given value from the set if it exists.
// Returns true if the value was removed, otherwise false if not set.
func (s Set[T]) RemoveTest(value T) bool {
	if _, has := s[value]; !has {
		return false
	}
	delete(s, value)
	return true
}

// RemoveIf removes any value which the given predicate accepts.
// Returns true if any value was remove, otherwise false if none were removed.
func (s Set[T]) RemoveIf(predicate func(T) bool) bool {
	removed := false
	maps.DeleteFunc(s, func(key T, _ struct{}) bool {
		if predicate(key) {
			removed = true
			return true
		}
		return false
	})
	return removed
}

// ToString gets a comma separated list of the values
// sorted by the strings of the values.
func (s Set[T]) ToString() string {
	keys := liteUtils.Strings(s.ToSlice())
	sort.Strings(keys)
	return strings.Join(keys, `, `)
}
