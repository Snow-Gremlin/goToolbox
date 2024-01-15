package collections

// ReadonlyDictionary is the interface for key value pairs
// which can not be directly modified.
//
// The keys are unique. Depending on the implementation
// the keys may be in sorted order or not.
type ReadonlyDictionary[TKey comparable, TValue any] interface {
	Collection[Tuple2[TKey, TValue]]
	Container[TKey]
	Getter[TKey, TValue]
	OnChanger

	// Keys enumerates the keys.
	//
	// Depending on the type of dictionary these may
	// be in random order or be sorted.
	Keys() Enumerator[TKey]

	// Values enumerates the values.
	//
	// Depending on the type of dictionary these may
	// be in random order or ordered to match the sorted keys.
	Values() Enumerator[TValue]

	// ToMap creates a map for this dictionary.
	ToMap() map[TKey]TValue
}
