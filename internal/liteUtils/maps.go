package liteUtils

// Keys gets all the keys for the given map in random order.
func Keys[TKey comparable, TValue any, TMap ~map[TKey]TValue](m TMap) []TKey {
	keys := make([]TKey, len(m))
	index := 0
	for key := range m {
		keys[index] = key
		index++
	}
	return keys
}

// Values gets all the values for the given map in random order.
func Values[TKey comparable, TValue any, TMap ~map[TKey]TValue](m TMap) []TValue {
	values := make([]TValue, len(m))
	index := 0
	for _, value := range m {
		values[index] = value
		index++
	}
	return values
}
