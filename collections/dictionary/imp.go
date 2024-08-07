package dictionary

import (
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
	"github.com/Snow-Gremlin/goToolbox/comp"
	"github.com/Snow-Gremlin/goToolbox/events"
	"github.com/Snow-Gremlin/goToolbox/events/event"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

type changeFlag int

const (
	noChange      changeFlag = 0
	addChange     changeFlag = 1
	removeChange  changeFlag = 2
	replaceChange changeFlag = addChange | removeChange
)

type dictionaryImp[TKey comparable, TValue any] struct {
	m     map[TKey]TValue
	event events.Event[collections.ChangeArgs]
}

func (d *dictionaryImp[TKey, TValue]) onChanged(cf changeFlag) bool {
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

func (d *dictionaryImp[TKey, TValue]) addOne(key TKey, val TValue) changeFlag {
	if v2, exists := d.m[key]; exists {
		if comp.Equal(val, v2) {
			return noChange
		}

		d.m[key] = val
		return replaceChange
	}

	d.m[key] = val
	return addChange
}

func (d *dictionaryImp[TKey, TValue]) addOneIfNotSet(key TKey, val TValue) changeFlag {
	if _, exists := d.m[key]; exists {
		return noChange
	}
	d.m[key] = val
	return addChange
}

func (d *dictionaryImp[TKey, TValue]) Add(key TKey, val TValue) bool {
	return d.onChanged(d.addOne(key, val))
}

func (d *dictionaryImp[TKey, TValue]) AddIfNotSet(key TKey, val TValue) bool {
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

func (d *dictionaryImp[TKey, TValue]) AddFrom(e collections.Enumerator[collections.Tuple2[TKey, TValue]]) bool {
	return d.onChanged(addFromTo(e, d.addOne))
}

func (d *dictionaryImp[TKey, TValue]) AddIfNotSetFrom(e collections.Enumerator[collections.Tuple2[TKey, TValue]]) bool {
	return d.onChanged(addFromTo(e, d.addOneIfNotSet))
}

func addMapTo[TKey comparable, TValue any](m map[TKey]TValue, addHandle func(key TKey, val TValue) changeFlag) changeFlag {
	result := noChange
	for key, value := range m {
		result |= addHandle(key, value)
	}
	return result
}

func (d *dictionaryImp[TKey, TValue]) AddMap(m map[TKey]TValue) bool {
	return d.onChanged(addMapTo(m, d.addOne))
}

func (d *dictionaryImp[TKey, TValue]) AddMapIfNotSet(m map[TKey]TValue) bool {
	return d.onChanged(addMapTo(m, d.addOneIfNotSet))
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
	result := noChange
	for _, key := range keys {
		if _, exists := d.m[key]; exists {
			delete(d.m, key)
			result = removeChange
		}
	}
	return d.onChanged(result)
}

func (d *dictionaryImp[TKey, TValue]) RemoveIf(p collections.Predicate[TKey]) bool {
	if utils.IsNil(p) {
		return false
	}
	result := noChange
	maps.DeleteFunc(d.m, func(key TKey, _ TValue) bool {
		if p(key) {
			result = removeChange
			return true
		}
		return false
	})
	return d.onChanged(result)
}

func (d *dictionaryImp[TKey, TValue]) Clear() {
	if len(d.m) > 0 {
		d.m = make(map[TKey]TValue)
		d.onChanged(removeChange)
	}
}

func (d *dictionaryImp[TKey, TValue]) Clone() collections.Dictionary[TKey, TValue] {
	return &dictionaryImp[TKey, TValue]{
		m:     maps.Clone(d.m),
		event: nil,
	}
}

func (d *dictionaryImp[TKey, TValue]) Readonly() collections.ReadonlyDictionary[TKey, TValue] {
	return readonlyDictionary.New(d)
}

func (d *dictionaryImp[TKey, TValue]) OnChange() events.Event[collections.ChangeArgs] {
	if d.event == nil {
		d.event = event.New[collections.ChangeArgs]()
	}
	return d.event
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
		if !ok || !comp.Equal(v2, value) {
			return false
		}
	}
	return true
}
