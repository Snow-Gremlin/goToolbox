package sortedDictionary

import (
	"maps"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/internal/optional"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

// New creates a new dictionary with sorted keys by the
// optional given comparer function or the default comparer.
func New[TKey comparable, TValue any](comparer ...utils.Comparer[TKey]) collections.Dictionary[TKey, TValue] {
	return CapNew[TKey, TValue](0, comparer...)
}

// CapNew creates a new dictionary with sorted keys and initial capacity
// by the optional given comparer function or the default comparer.
func CapNew[TKey comparable, TValue any](capacity int, comparer ...utils.Comparer[TKey]) collections.Dictionary[TKey, TValue] {
	cmp := optional.Comparer(comparer)
	capacity = max(capacity, 0)
	return &sortedDictionaryImp[TKey, TValue]{
		data:     make(map[TKey]TValue, capacity),
		keys:     make([]TKey, 0, capacity),
		comparer: cmp,
		event:    nil,
	}
}

// With creates a new dictionary with sorted keys
// populated with key/value pairs from the given map.
//
// The keys are sorted with the optional given comparer function
// or the default comparer if no comparer was given.
func With[TKey comparable, TValue any, M ~map[TKey]TValue](m M, comparer ...utils.Comparer[TKey]) collections.Dictionary[TKey, TValue] {
	cmp := optional.Comparer(comparer)
	data := maps.Clone(m)
	if data == nil {
		data = make(map[TKey]TValue)
	}
	return &sortedDictionaryImp[TKey, TValue]{
		data:     data,
		keys:     utils.SortedKeys(m, cmp),
		comparer: cmp,
		event:    nil,
	}
}

// From creates a new dictionary with sorted keys
// populated with key/value pairs from the given tuple enumerator.
//
// The keys are sorted with the optional given comparer function
// or the default comparer if no comparer was given.
func From[TKey comparable, TValue any](e collections.Enumerator[collections.Tuple2[TKey, TValue]], comparer ...utils.Comparer[TKey]) collections.Dictionary[TKey, TValue] {
	d := CapNew[TKey, TValue](0, comparer...)
	d.AddFrom(e)
	return d
}

// CapFrom creates a new dictionary with sorted keys and an initial capacity
// populated with key/value pairs from the given tuple enumerator.
//
// The keys are sorted with the optional given comparer function
// or the default comparer if no comparer was given.
func CapFrom[TKey comparable, TValue any](e collections.Enumerator[collections.Tuple2[TKey, TValue]], capacity int, comparer ...utils.Comparer[TKey]) collections.Dictionary[TKey, TValue] {
	d := CapNew[TKey, TValue](capacity, comparer...)
	d.AddFrom(e)
	return d
}
