package readonlyDictionary

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/collections/tuple2"
	"github.com/Snow-Gremlin/goToolbox/testers/check"
)

type pseudoDic[TKey comparable, TValue any] map[TKey]TValue

func (m pseudoDic[TKey, TValue]) Get(key TKey) TValue {
	return m[key]
}

func (m pseudoDic[TKey, TValue]) TryGet(key TKey) (TValue, bool) {
	v, ok := m[key]
	return v, ok
}

func (m pseudoDic[TKey, TValue]) Enumerate() collections.Enumerator[collections.Tuple2[TKey, TValue]] {
	list := make([]collections.Tuple2[TKey, TValue], 0, len(m))
	for key, value := range m {
		list = append(list, tuple2.New(key, value))
	}
	return enumerator.Enumerate(list...)
}

func (m pseudoDic[TKey, TValue]) Keys() collections.Enumerator[TKey] {
	list := make([]TKey, 0, len(m))
	for key := range m {
		list = append(list, key)
	}
	return enumerator.Enumerate(list...)
}

func (m pseudoDic[TKey, TValue]) Values() collections.Enumerator[TValue] {
	list := make([]TValue, 0, len(m))
	for _, value := range m {
		list = append(list, value)
	}
	return enumerator.Enumerate(list...)
}

func (m pseudoDic[TKey, TValue]) ToMap() map[TKey]TValue {
	return m
}

func (m pseudoDic[TKey, TValue]) Empty() bool {
	return len(m) <= 0
}

func (m pseudoDic[TKey, TValue]) Count() int {
	return len(m)
}

func (m pseudoDic[TKey, TValue]) Contains(key TKey) bool {
	_, ok := m[key]
	return ok
}

func (m pseudoDic[TKey, TValue]) String() string {
	return fmt.Sprintf(`%v`, (map[TKey]TValue)(m))
}

func (m pseudoDic[TKey, TValue]) Equals(other any) bool {
	m2, ok := other.(collections.ReadonlyDictionary[TKey, TValue])
	return ok && reflect.DeepEqual(m.ToMap(), m2.ToMap())
}

func Test_Dictionary_Readonly(t *testing.T) {
	d1 := pseudoDic[int, int]{}
	r1 := New(d1)
	check.Empty(t).Assert(r1)
	check.True(t).Assert(r1.Empty())

	d1[123] = 456
	check.Length(t, 1).Assert(r1)
	check.False(t).Assert(r1.Empty())
	check.Equal(t, 456).Assert(r1.Get(123))
	check.True(t).Assert(r1.Contains(123))
	check.False(t).Assert(r1.Contains(765))

	d1[222] = 333
	check.Length(t, 2).Assert(r1)
	check.Equal(t, 456).Assert(r1.Get(123))
	check.Equal(t, 333).Assert(r1.Get(222))

	v, ok := r1.TryGet(222)
	check.Equal(t, 333).Assert(v)
	check.True(t).Assert(ok)

	v, ok = r1.TryGet(251)
	check.Zero(t).Assert(v)
	check.False(t).Assert(ok)

	check.Equal(t, `map[123:456 222:333]`).Assert(r1.String())
	check.Equal(t, map[int]int{123: 456, 222: 333}).Assert(r1.ToMap())
	check.True(t).Name(`d1.Equals(d2)`).Assert(d1.Equals(r1))
	check.True(t).Name(`d2.Equals(d1)`).Assert(r1.Equals(d1))

	check.Equal(t, `[123, 456]|[222, 333]`).Assert(r1.Enumerate().Strings().Sort().Join(`|`))
	check.Equal(t, `123|222`).Assert(r1.Keys().Sort().Join(`|`))
	check.Equal(t, `333|456`).Assert(r1.Values().Sort().Join(`|`))
}
