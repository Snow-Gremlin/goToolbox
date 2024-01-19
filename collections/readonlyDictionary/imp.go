package readonlyDictionary

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/events"
)

type readonlyDictionaryImp[TKey comparable, TValue any] struct {
	dic collections.ReadonlyDictionary[TKey, TValue]
}

func (r readonlyDictionaryImp[TKey, TValue]) Get(key TKey) TValue {
	return r.dic.Get(key)
}

func (r readonlyDictionaryImp[TKey, TValue]) TryGet(key TKey) (TValue, bool) {
	return r.dic.TryGet(key)
}

func (r readonlyDictionaryImp[TKey, TValue]) Enumerate() collections.Enumerator[collections.Tuple2[TKey, TValue]] {
	return r.dic.Enumerate()
}

func (r readonlyDictionaryImp[TKey, TValue]) Keys() collections.Enumerator[TKey] {
	return r.dic.Keys()
}

func (r readonlyDictionaryImp[TKey, TValue]) Values() collections.Enumerator[TValue] {
	return r.dic.Values()
}

func (r readonlyDictionaryImp[TKey, TValue]) ToMap() map[TKey]TValue {
	return r.dic.ToMap()
}

func (r readonlyDictionaryImp[TKey, TValue]) Empty() bool {
	return r.dic.Empty()
}

func (r readonlyDictionaryImp[TKey, TValue]) Count() int {
	return r.dic.Count()
}

func (r readonlyDictionaryImp[TKey, TValue]) Contains(key TKey) bool {
	return r.dic.Contains(key)
}

func (r readonlyDictionaryImp[TKey, TValue]) String() string {
	return r.dic.String()
}

func (r readonlyDictionaryImp[TKey, TValue]) Equals(other any) bool {
	return r.dic.Equals(other)
}

func (r readonlyDictionaryImp[TKey, TValue]) OnChange() events.Event[collections.ChangeArgs] {
	return r.dic.OnChange()
}
