package dictionary

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"goToolbox/collections"
	"goToolbox/collections/enumerator"
	"goToolbox/collections/iterator"
	"goToolbox/collections/readonlyDictionary"
	"goToolbox/collections/tuple2"
	"goToolbox/utils"
)

type dictionaryImp[TKey comparable, TValue any] struct {
	m map[TKey]TValue
}

func (d *dictionaryImp[TKey, TValue]) Add(key TKey, val TValue) bool {
	_, exists := d.m[key]
	d.m[key] = val
	return !exists
}

func (d *dictionaryImp[TKey, TValue]) AddIfNotSet(key TKey, val TValue) bool {
	if _, exists := d.m[key]; exists {
		return false
	}
	d.m[key] = val
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

func (d *dictionaryImp[TKey, TValue]) AddFrom(e collections.Enumerator[collections.Tuple2[TKey, TValue]]) bool {
	return addFromTo(e, d.Add)
}

func (d *dictionaryImp[TKey, TValue]) AddIfNotSetFrom(e collections.Enumerator[collections.Tuple2[TKey, TValue]]) bool {
	return addFromTo(e, d.AddIfNotSet)
}

func addMapTo[TKey comparable, TValue any](m map[TKey]TValue, addHandle func(key TKey, val TValue) bool) bool {
	result := false
	for key, value := range m {
		result = addHandle(key, value) || result
	}
	return result
}

func (d *dictionaryImp[TKey, TValue]) AddMap(m map[TKey]TValue) bool {
	return addMapTo(m, d.Add)
}

func (d *dictionaryImp[TKey, TValue]) AddMapIfNotSet(m map[TKey]TValue) bool {
	return addMapTo(m, d.AddIfNotSet)
}

func (d *dictionaryImp[TKey, TValue]) Get(key TKey) TValue {
	return d.m[key]
}

func (d *dictionaryImp[TKey, TValue]) TryGet(key TKey) (TValue, bool) {
	value, exists := d.m[key]
	return value, exists
}

func (d *dictionaryImp[TKey, TValue]) ToMap() map[TKey]TValue {
	return maps.Clone(d.m)
}

func (d *dictionaryImp[TKey, TValue]) Remove(keys ...TKey) bool {
	removed := false
	for _, key := range keys {
		if _, exists := d.m[key]; exists {
			delete(d.m, key)
			removed = true
		}
	}
	return removed
}

func (d *dictionaryImp[TKey, TValue]) RemoveIf(p collections.Predicate[TKey]) bool {
	if utils.IsNil(p) {
		return false
	}
	count := len(d.m)
	maps.DeleteFunc(d.m, func(key TKey, _ TValue) bool { return p(key) })
	return len(d.m) != count
}

func (d *dictionaryImp[TKey, TValue]) Clear() {
	d.m = make(map[TKey]TValue)
}

func (d *dictionaryImp[TKey, TValue]) Clone() collections.Dictionary[TKey, TValue] {
	return &dictionaryImp[TKey, TValue]{
		m: maps.Clone(d.m),
	}
}

func (d *dictionaryImp[TKey, TValue]) Readonly() collections.ReadonlyDictionary[TKey, TValue] {
	return readonlyDictionary.New(d)
}

func (d *dictionaryImp[TKey, TValue]) Enumerate() collections.Enumerator[collections.Tuple2[TKey, TValue]] {
	// Since Go randomizes the order of values, to keep a consistent
	// iteration, all the keys must be collected once before iteration.
	// The keys will still be in random order but consistent.
	// Because the keys are collected ahead of iterations, changes to
	// the dictionary may just cause the enumeration to be unstable
	// but doesn't require it to be stopped.
	return enumerator.New(func() collections.Iterator[collections.Tuple2[TKey, TValue]] {
		keys := utils.Keys(d.m)
		index, count := -1, len(keys)-1
		return iterator.New(func() (collections.Tuple2[TKey, TValue], bool) {
			for index < count {
				index++
				key := keys[index]
				if value, ok := d.m[key]; ok {
					return tuple2.New(key, value), true
				}
			}
			return utils.Zero[collections.Tuple2[TKey, TValue]](), false
		})
	})
}

func (d *dictionaryImp[TKey, TValue]) Keys() collections.Enumerator[TKey] {
	// See comment in Enumerate
	return enumerator.New(func() collections.Iterator[TKey] {
		keys := utils.Keys(d.m)
		index, count := -1, len(keys)-1
		return iterator.New(func() (TKey, bool) {
			if index < count {
				index++
				return keys[index], true
			}
			return utils.Zero[TKey](), false
		})
	})
}

func (d *dictionaryImp[TKey, TValue]) Values() collections.Enumerator[TValue] {
	// See comment in Enumerate
	return enumerator.New(func() collections.Iterator[TValue] {
		values := utils.Values(d.m)
		index, count := -1, len(values)-1
		return iterator.New(func() (TValue, bool) {
			if index < count {
				index++
				return values[index], true
			}
			return utils.Zero[TValue](), false
		})
	})
}

func (d *dictionaryImp[TKey, TValue]) Empty() bool {
	return len(d.m) <= 0
}

func (d *dictionaryImp[TKey, TValue]) Count() int {
	return len(d.m)
}

func (d *dictionaryImp[TKey, TValue]) Contains(key TKey) bool {
	_, contains := d.m[key]
	return contains
}

func (d *dictionaryImp[TKey, TValue]) String() string {
	const newline = "\n"
	keys := utils.Keys(d.m)
	keyStr := utils.Strings(keys)
	maxWidth := utils.GetMaxStringLen(keyStr) + 2
	padding := newline + strings.Repeat(` `, maxWidth)
	lines := make([]string, len(keys))
	for i, key := range keys {
		value := utils.String(d.m[key])
		value = strings.ReplaceAll(value, newline, padding)
		lines[i] = fmt.Sprintf(`%-*s%s`, maxWidth, keyStr[i]+`: `, value)
	}
	slices.Sort(lines)
	return strings.Join(lines, newline)
}

func (d *dictionaryImp[TKey, TValue]) Equals(other any) bool {
	d2, ok := other.(collections.Collection[collections.Tuple2[TKey, TValue]])
	if !ok || d.Count() != d2.Count() {
		return false
	}

	it := d2.Enumerate().Iterate()
	for it.Next() {
		key, value := it.Current().Values()
		v2, ok := d.TryGet(key)
		if !ok || !utils.Equal(v2, value) {
			return false
		}
	}
	return true
}
