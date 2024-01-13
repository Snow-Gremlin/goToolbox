package readonlyDictionary

import "github.com/Snow-Gremlin/goToolbox/collections"

type readonlyDictionaryImp[TKey comparable, TValue any] struct {
	dic collections.ReadonlyDictionary[TKey, TValue]
}

func (m *readonlyDictionaryImp[TKey, TValue]) Get(key TKey) TValue {
	return m.dic.Get(key)
}

func (m *readonlyDictionaryImp[TKey, TValue]) TryGet(key TKey) (TValue, bool) {
	return m.dic.TryGet(key)
}

func (m *readonlyDictionaryImp[TKey, TValue]) Enumerate() collections.Enumerator[collections.Tuple2[TKey, TValue]] {
	return m.dic.Enumerate()
}

func (m *readonlyDictionaryImp[TKey, TValue]) Keys() collections.Enumerator[TKey] {
	return m.dic.Keys()
}

func (m *readonlyDictionaryImp[TKey, TValue]) Values() collections.Enumerator[TValue] {
	return m.dic.Values()
}

func (m *readonlyDictionaryImp[TKey, TValue]) ToMap() map[TKey]TValue {
	return m.dic.ToMap()
}

func (m *readonlyDictionaryImp[TKey, TValue]) Empty() bool {
	return m.dic.Empty()
}

func (m *readonlyDictionaryImp[TKey, TValue]) Count() int {
	return m.dic.Count()
}

func (m *readonlyDictionaryImp[TKey, TValue]) Contains(key TKey) bool {
	return m.dic.Contains(key)
}

func (m *readonlyDictionaryImp[TKey, TValue]) String() string {
	return m.dic.String()
}

func (m *readonlyDictionaryImp[TKey, TValue]) Equals(other any) bool {
	return m.dic.Equals(other)
}
