package readonlyDictionary

import "goToolbox/collections"

// New wraps another dictionary in a readonly shell.
func New[TKey comparable, TValue any](dic collections.ReadonlyDictionary[TKey, TValue]) collections.ReadonlyDictionary[TKey, TValue] {
	return &readonlyDictionaryImp[TKey, TValue]{
		dic: dic,
	}
}
