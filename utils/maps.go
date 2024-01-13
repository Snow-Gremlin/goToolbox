package utils

import (
	"slices"

	"github.com/Snow-Gremlin/goToolbox/internal/liteUtils"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
)

// Keys gets all the keys for the given map in random order.
func Keys[TKey comparable, TValue any, TMap ~map[TKey]TValue](m TMap) []TKey {
	return liteUtils.Keys(m)
}

// SortedKeys gets all the keys from the given map in sorted order.
//
// An optional comparer maybe added to override the default comparer
// or for types that don't have a default comparer.
func SortedKeys[TKey comparable, TValue any, TMap ~map[TKey]TValue](m TMap, comparer ...Comparer[TKey]) []TKey {
	var cmp Comparer[TKey]
	if count := len(comparer); count > 0 {
		if count > 1 {
			panic(terror.InvalidArgCount(1, count, `comparer`))
		}
		cmp = comparer[0]
	}

	if IsNil(cmp) {
		cmp = DefaultComparer[TKey]()
		if IsNil(cmp) {
			panic(terror.New(`must provide a comparer to compare this type`).
				With(`type`, TypeOf[TKey]()))
		}
	}

	keys := Keys(m)
	slices.SortFunc(keys, cmp)
	return keys
}

// Values gets all the values for the given map in random order.
func Values[TKey comparable, TValue any, TMap ~map[TKey]TValue](m TMap) []TValue {
	return liteUtils.Values(m)
}
