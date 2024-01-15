package collections

// Dictionary is the interface for key/value pairs.
//
// The keys are unique. Depending on the implementation
// the keys may be in sorted order or not.
type Dictionary[TKey comparable, TValue any] interface {
	ReadonlyDictionary[TKey, TValue]

	// Add will add or overwrite the key with the given value.
	// Returns true if the key was added or, if the key
	// existed but the value is different, otherwise returns false.
	Add(key TKey, value TValue) bool

	// AddIfNotSet will add the given key with the given value if the
	// given key doesn't exist. If the key exists the value is not overwritten.
	// Returns true if the key was added or false if not added.
	AddIfNotSet(key TKey, value TValue) bool

	// AddFrom adds all the key/value pairs from the tuples.
	// This will overwrite any existing value with the same key.
	// Returns true if any key/value was added or overwritten,
	// false if none were changed.
	AddFrom(e Enumerator[Tuple2[TKey, TValue]]) bool

	// AddIfNotSetFrom adds all the key/value pairs
	// from the tuples for each key that doesn't exist.
	// If the key exists the value is not overwritten.
	// Returns true if any key was added, false if none were added.
	AddIfNotSetFrom(e Enumerator[Tuple2[TKey, TValue]]) bool

	// AddMap adds all the key/value pairs from the map.
	// This will overwrite any existing value with the same key.
	// Returns true if any key/value was added or overwritten,
	// false if none were changed.
	AddMap(m map[TKey]TValue) bool

	// AddIfNotSetMap adds all the key/value pairs from the map
	// for each key that doesn't exist.
	// If the key exists the value is not overwritten.
	// Returns true if any key was added, false if none were added.
	AddMapIfNotSet(m map[TKey]TValue) bool

	// Remove removes the given keys.
	// Returns true if any key existed and was removed, false if none of the keys existed.
	Remove(keys ...TKey) bool

	// RemoveIf removes the keys that the predicate returns true for.
	// Returns true if any key was removed, false if nothing was removed.
	RemoveIf(p Predicate[TKey]) bool

	// Clear removes all the values from the dictionary.
	Clear()

	// Clones this dictionary.
	Clone() Dictionary[TKey, TValue]

	// Readonly gets a readonly version of this dictionary.
	//
	// The readonly version points back to this dictionary
	// but is not able to be cast into this dictionary.
	Readonly() ReadonlyDictionary[TKey, TValue]
}
