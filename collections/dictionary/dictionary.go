package dictionary

import (
	"maps"

	"goToolbox/collections"
	"goToolbox/internal/optional"
)

// New creates a new dictionary with unsorted keys.
//
// The keys will be returned in random order when enumeration
// and may have different orders per enumeration.
//
// If one capacity value is given, an empty underlying map is allocated
// with enough space to hold the specified number of elements.
// The capacity may be omitted, in which case a small starting size is allocated.
func New[TKey comparable, TValue any](capacity ...int) collections.Dictionary[TKey, TValue] {
	cap := optional.Capacity(capacity)
	d := &dictionaryImp[TKey, TValue]{}
	d.m = make(map[TKey]TValue, cap)
	return d
}

// With creates a new dictionary with unsorted keys
// populated with key/value pairs from the given map.
func With[TKey comparable, TValue any](m map[TKey]TValue) collections.Dictionary[TKey, TValue] {
	return &dictionaryImp[TKey, TValue]{
		m: maps.Clone(m),
	}
}

// From creates a new dictionary with unsorted keys
// populated with key/value pairs from the given tuple enumerator.
func From[TKey comparable, TValue any](e collections.Enumerator[collections.Tuple2[TKey, TValue]], capacity ...int) collections.Dictionary[TKey, TValue] {
	d := New[TKey, TValue](capacity...)
	d.AddFrom(e)
	return d
}
