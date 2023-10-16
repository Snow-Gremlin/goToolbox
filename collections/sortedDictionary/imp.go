package sortedDictionary

import (
	"bytes"
	"fmt"
	"maps"
	"slices"
	"strings"

	"goToolbox/collections"
	"goToolbox/collections/enumerator"
	"goToolbox/collections/iterator"
	"goToolbox/collections/readonlyDictionary"
	"goToolbox/collections/tuple2"
	"goToolbox/internal/simpleSet"
	"goToolbox/utils"
)

type sortedImp[TKey comparable, TValue any] struct {
	data     map[TKey]TValue
	keys     []TKey
	comparer utils.Comparer[TKey]
}

func (m *sortedImp[TKey, TValue]) insertKey(key TKey) {
	if index, found := slices.BinarySearchFunc(m.keys, key, m.comparer); !found {
		m.keys = slices.Insert(m.keys, index, key)
	}
}

func (m *sortedImp[TKey, TValue]) removeKeys(keyToRemove simpleSet.Set[TKey]) {
	newKeys := slices.DeleteFunc(m.keys, keyToRemove.Has)
	zero := utils.Zero[TKey]()
	for i, count := len(newKeys), len(m.keys); i < count; i++ {
		m.keys[i] = zero
	}
	m.keys = newKeys
}

func (m *sortedImp[TKey, TValue]) Add(key TKey, val TValue) bool {
	_, exists := m.data[key]
	m.data[key] = val
	if !exists {
		m.insertKey(key)
	}
	return !exists
}

func (m *sortedImp[TKey, TValue]) AddIfNotSet(key TKey, val TValue) bool {
	if _, exists := m.data[key]; exists {
		return false
	}
	m.data[key] = val
	m.insertKey(key)
	return true
}

func addFromTo[TKey comparable, TValue any](e collections.Enumerator[collections.Tuple2[TKey, TValue]], addHandle func(key TKey, val TValue) bool) bool {
	if utils.IsNil(e) {
		return false
	}
	result := false
	e.All(func(t collections.Tuple2[TKey, TValue]) bool {
		result = addHandle(t.Values()) || result
		return true
	})
	return result
}

func (m *sortedImp[TKey, TValue]) AddFrom(e collections.Enumerator[collections.Tuple2[TKey, TValue]]) bool {
	return addFromTo(e, m.Add)
}

func (m *sortedImp[TKey, TValue]) AddIfNotSetFrom(e collections.Enumerator[collections.Tuple2[TKey, TValue]]) bool {
	return addFromTo(e, m.AddIfNotSet)
}

func addMapTo[TKey comparable, TValue any](data map[TKey]TValue, addHandle func(key TKey, val TValue) bool) bool {
	result := false
	for key, value := range data {
		result = addHandle(key, value) || result
	}
	return result
}

func (m *sortedImp[TKey, TValue]) AddMap(data map[TKey]TValue) bool {
	return addMapTo(data, m.Add)
}

func (m *sortedImp[TKey, TValue]) AddMapIfNotSet(data map[TKey]TValue) bool {
	return addMapTo(data, m.AddIfNotSet)
}

func (m *sortedImp[TKey, TValue]) Get(key TKey) TValue {
	return m.data[key]
}

func (m *sortedImp[TKey, TValue]) TryGet(key TKey) (TValue, bool) {
	value, exists := m.data[key]
	return value, exists
}

func (m *sortedImp[TKey, TValue]) ToMap() map[TKey]TValue {
	return maps.Clone(m.data)
}

func (m *sortedImp[TKey, TValue]) Remove(keys ...TKey) bool {
	removed := simpleSet.New[TKey]()
	for _, key := range keys {
		if _, exists := m.data[key]; exists {
			delete(m.data, key)
			removed.Set(key)
		}
	}
	if removed.Count() > 0 {
		m.removeKeys(removed)
		return true
	}
	return false
}

func (m *sortedImp[TKey, TValue]) RemoveIf(p collections.Predicate[TKey]) bool {
	if utils.IsNil(p) {
		return false
	}
	removed := simpleSet.New[TKey]()
	for _, key := range m.keys {
		if p(key) {
			delete(m.data, key)
			removed.Set(key)
		}
	}
	if removed.Count() > 0 {
		m.removeKeys(removed)
		return true
	}
	return false
}

func (m *sortedImp[TKey, TValue]) Clear() {
	m.data = make(map[TKey]TValue)
	m.keys = []TKey{}
}

func (m *sortedImp[TKey, TValue]) Clone() collections.Dictionary[TKey, TValue] {
	return &sortedImp[TKey, TValue]{
		data:     maps.Clone(m.data),
		keys:     slices.Clone(m.keys),
		comparer: m.comparer,
	}
}

func (m *sortedImp[TKey, TValue]) Readonly() collections.ReadonlyDictionary[TKey, TValue] {
	return readonlyDictionary.New(m)
}

func (m *sortedImp[TKey, TValue]) Enumerate() collections.Enumerator[collections.Tuple2[TKey, TValue]] {
	return enumerator.New(func() collections.Iterator[collections.Tuple2[TKey, TValue]] {
		index := 0
		return iterator.New(func() (collections.Tuple2[TKey, TValue], bool) {
			for index < len(m.keys) {
				key := m.keys[index]
				index++
				if value, ok := m.data[key]; ok {
					return tuple2.New(key, value), true
				}
			}
			return utils.Zero[collections.Tuple2[TKey, TValue]](), false
		})
	})
}

func (m *sortedImp[TKey, TValue]) Keys() collections.Enumerator[TKey] {
	return enumerator.New(func() collections.Iterator[TKey] {
		index := 0
		return iterator.New(func() (TKey, bool) {
			for index < len(m.keys) {
				key := m.keys[index]
				index++
				return key, true
			}
			return utils.Zero[TKey](), false
		})
	})
}

func (m *sortedImp[TKey, TValue]) Values() collections.Enumerator[TValue] {
	return enumerator.New(func() collections.Iterator[TValue] {
		index := 0
		return iterator.New(func() (TValue, bool) {
			for index < len(m.keys) {
				value := m.data[m.keys[index]]
				index++
				return value, true
			}
			return utils.Zero[TValue](), false
		})
	})
}

func (m *sortedImp[TKey, TValue]) Empty() bool {
	return len(m.data) <= 0
}

func (m *sortedImp[TKey, TValue]) Count() int {
	return len(m.data)
}

func (m *sortedImp[TKey, TValue]) Contains(key TKey) bool {
	_, contains := m.data[key]
	return contains
}

func (m *sortedImp[TKey, TValue]) String() string {
	const newline = "\n"
	keyStr := utils.Strings(m.keys)
	maxWidth := utils.GetMaxStringLen(keyStr) + 2
	padding := newline + strings.Repeat(` `, maxWidth)
	buf := &bytes.Buffer{}
	for i, key := range m.keys {
		if i > 0 {
			_, _ = buf.WriteString(newline)
		}
		value := utils.String(m.data[key])
		value = strings.ReplaceAll(value, newline, padding)
		_, _ = buf.WriteString(fmt.Sprintf(`%-*s`, maxWidth, keyStr[i]+`: `))
		_, _ = buf.WriteString(value)
	}
	return buf.String()
}

func (m *sortedImp[TKey, TValue]) Equals(other any) bool {
	d2, ok := other.(collections.Collection[collections.Tuple2[TKey, TValue]])
	if !ok || m.Count() != d2.Count() {
		return false
	}

	it := d2.Enumerate().Iterate()
	for it.Next() {
		key, value := it.Current().Values()
		v2, ok := m.TryGet(key)
		if !ok || !utils.Equal(v2, value) {
			return false
		}
	}
	return true
}
