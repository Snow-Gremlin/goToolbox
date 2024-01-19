package sortedDictionary

import (
	"bytes"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/changeArgs"
	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/collections/iterator"
	"github.com/Snow-Gremlin/goToolbox/collections/readonlyDictionary"
	"github.com/Snow-Gremlin/goToolbox/collections/tuple2"
	"github.com/Snow-Gremlin/goToolbox/events"
	"github.com/Snow-Gremlin/goToolbox/events/event"
	"github.com/Snow-Gremlin/goToolbox/internal/simpleSet"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

type changeFlag int

const (
	noChange      changeFlag = 0
	addChange     changeFlag = 1
	removeChange  changeFlag = 2
	replaceChange changeFlag = addChange | removeChange
)

type sortedDictionaryImp[TKey comparable, TValue any] struct {
	data     map[TKey]TValue
	keys     []TKey
	comparer utils.Comparer[TKey]
	event    events.Event[collections.ChangeArgs]
}

func (d *sortedDictionaryImp[TKey, TValue]) onChanged(cf changeFlag) bool {
	if d.event != nil {
		switch cf {
		case addChange:
			d.event.Invoke(changeArgs.NewAdded())
		case removeChange:
			d.event.Invoke(changeArgs.NewRemoved())
		case replaceChange:
			d.event.Invoke(changeArgs.NewReplaced())
		}
	}
	return cf != noChange
}

func (d *sortedDictionaryImp[TKey, TValue]) insertKey(key TKey) {
	if index, found := slices.BinarySearchFunc(d.keys, key, d.comparer); !found {
		d.keys = slices.Insert(d.keys, index, key)
	}
}

func (d *sortedDictionaryImp[TKey, TValue]) removeKeys(keyToRemove simpleSet.Set[TKey]) {
	newKeys := slices.DeleteFunc(d.keys, keyToRemove.Has)
	zero := utils.Zero[TKey]()
	for i, count := len(newKeys), len(d.keys); i < count; i++ {
		d.keys[i] = zero
	}
	d.keys = newKeys
}

func (d *sortedDictionaryImp[TKey, TValue]) addOne(key TKey, val TValue) changeFlag {
	if v2, exists := d.data[key]; exists {
		if utils.Equal(val, v2) {
			return noChange
		}

		d.data[key] = val
		return replaceChange
	}

	d.data[key] = val
	d.insertKey(key)
	return addChange
}

func (d *sortedDictionaryImp[TKey, TValue]) addOneIfNotSet(key TKey, val TValue) changeFlag {
	if _, exists := d.data[key]; exists {
		return noChange
	}
	d.data[key] = val
	d.insertKey(key)
	return addChange
}

func (d *sortedDictionaryImp[TKey, TValue]) Add(key TKey, val TValue) bool {
	return d.onChanged(d.addOne(key, val))
}

func (d *sortedDictionaryImp[TKey, TValue]) AddIfNotSet(key TKey, val TValue) bool {
	return d.onChanged(d.addOneIfNotSet(key, val))
}

func addFromTo[TKey comparable, TValue any](e collections.Enumerator[collections.Tuple2[TKey, TValue]], addHandle func(key TKey, val TValue) changeFlag) changeFlag {
	if utils.IsNil(e) {
		return noChange
	}
	result := noChange
	e.All(func(t collections.Tuple2[TKey, TValue]) bool {
		result |= addHandle(t.Values())
		return true
	})
	return result
}

func (d *sortedDictionaryImp[TKey, TValue]) AddFrom(e collections.Enumerator[collections.Tuple2[TKey, TValue]]) bool {
	return d.onChanged(addFromTo(e, d.addOne))
}

func (d *sortedDictionaryImp[TKey, TValue]) AddIfNotSetFrom(e collections.Enumerator[collections.Tuple2[TKey, TValue]]) bool {
	return d.onChanged(addFromTo(e, d.addOneIfNotSet))
}

func addMapTo[TKey comparable, TValue any](data map[TKey]TValue, addHandle func(key TKey, val TValue) changeFlag) changeFlag {
	result := noChange
	for key, value := range data {
		result |= addHandle(key, value)
	}
	return result
}

func (d *sortedDictionaryImp[TKey, TValue]) AddMap(m map[TKey]TValue) bool {
	return d.onChanged(addMapTo(m, d.addOne))
}

func (d *sortedDictionaryImp[TKey, TValue]) AddMapIfNotSet(m map[TKey]TValue) bool {
	return d.onChanged(addMapTo(m, d.addOneIfNotSet))
}

func (d *sortedDictionaryImp[TKey, TValue]) Get(key TKey) TValue {
	return d.data[key]
}

func (d *sortedDictionaryImp[TKey, TValue]) TryGet(key TKey) (TValue, bool) {
	value, exists := d.data[key]
	return value, exists
}

func (d *sortedDictionaryImp[TKey, TValue]) ToMap() map[TKey]TValue {
	return maps.Clone(d.data)
}

func (d *sortedDictionaryImp[TKey, TValue]) Remove(keys ...TKey) bool {
	removed := simpleSet.New[TKey]()
	for _, key := range keys {
		if _, exists := d.data[key]; exists {
			delete(d.data, key)
			removed.Set(key)
		}
	}
	if removed.Count() > 0 {
		d.removeKeys(removed)
		d.onChanged(removeChange)
		return true
	}
	return false
}

func (d *sortedDictionaryImp[TKey, TValue]) RemoveIf(p collections.Predicate[TKey]) bool {
	if utils.IsNil(p) {
		return false
	}
	removed := simpleSet.New[TKey]()
	for _, key := range d.keys {
		if p(key) {
			delete(d.data, key)
			removed.Set(key)
		}
	}
	if removed.Count() > 0 {
		d.removeKeys(removed)
		d.onChanged(removeChange)
		return true
	}
	return false
}

func (d *sortedDictionaryImp[TKey, TValue]) Clear() {
	if len(d.data) > 0 {
		d.data = make(map[TKey]TValue)
		d.keys = []TKey{}
		d.onChanged(removeChange)
	}
}

func (d *sortedDictionaryImp[TKey, TValue]) Clone() collections.Dictionary[TKey, TValue] {
	return &sortedDictionaryImp[TKey, TValue]{
		data:     maps.Clone(d.data),
		keys:     slices.Clone(d.keys),
		comparer: d.comparer,
		event:    nil,
	}
}

func (d *sortedDictionaryImp[TKey, TValue]) OnChange() events.Event[collections.ChangeArgs] {
	if d.event == nil {
		d.event = event.New[collections.ChangeArgs]()
	}
	return d.event
}

func (d *sortedDictionaryImp[TKey, TValue]) Readonly() collections.ReadonlyDictionary[TKey, TValue] {
	return readonlyDictionary.New(d)
}

func (d *sortedDictionaryImp[TKey, TValue]) Enumerate() collections.Enumerator[collections.Tuple2[TKey, TValue]] {
	return enumerator.New(func() collections.Iterator[collections.Tuple2[TKey, TValue]] {
		index := 0
		return iterator.New(func() (collections.Tuple2[TKey, TValue], bool) {
			for index < len(d.keys) {
				key := d.keys[index]
				index++
				if value, ok := d.data[key]; ok {
					return tuple2.New(key, value), true
				}
			}
			return utils.Zero[collections.Tuple2[TKey, TValue]](), false
		})
	})
}

func (d *sortedDictionaryImp[TKey, TValue]) Keys() collections.Enumerator[TKey] {
	return enumerator.New(func() collections.Iterator[TKey] {
		index := 0
		return iterator.New(func() (TKey, bool) {
			for index < len(d.keys) {
				key := d.keys[index]
				index++
				return key, true
			}
			return utils.Zero[TKey](), false
		})
	})
}

func (d *sortedDictionaryImp[TKey, TValue]) Values() collections.Enumerator[TValue] {
	return enumerator.New(func() collections.Iterator[TValue] {
		index := 0
		return iterator.New(func() (TValue, bool) {
			for index < len(d.keys) {
				value := d.data[d.keys[index]]
				index++
				return value, true
			}
			return utils.Zero[TValue](), false
		})
	})
}

func (d *sortedDictionaryImp[TKey, TValue]) Empty() bool {
	return len(d.data) <= 0
}

func (d *sortedDictionaryImp[TKey, TValue]) Count() int {
	return len(d.data)
}

func (d *sortedDictionaryImp[TKey, TValue]) Contains(key TKey) bool {
	_, contains := d.data[key]
	return contains
}

func (d *sortedDictionaryImp[TKey, TValue]) String() string {
	const newline = "\n"
	keyStr := utils.Strings(d.keys)
	maxWidth := utils.GetMaxStringLen(keyStr) + 2
	padding := newline + strings.Repeat(` `, maxWidth)
	buf := &bytes.Buffer{}
	for i, key := range d.keys {
		if i > 0 {
			_, _ = buf.WriteString(newline)
		}
		value := utils.String(d.data[key])
		value = strings.ReplaceAll(value, newline, padding)
		_, _ = buf.WriteString(fmt.Sprintf(`%-*s`, maxWidth, keyStr[i]+`: `))
		_, _ = buf.WriteString(value)
	}
	return buf.String()
}

func (d *sortedDictionaryImp[TKey, TValue]) Equals(other any) bool {
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
