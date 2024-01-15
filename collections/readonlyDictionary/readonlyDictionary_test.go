package readonlyDictionary

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/collections/tuple2"
	"github.com/Snow-Gremlin/goToolbox/events"
	"github.com/Snow-Gremlin/goToolbox/events/event"
	"github.com/Snow-Gremlin/goToolbox/testers/check"
)

type pseudoDic[TKey comparable, TValue any] struct {
	m map[TKey]TValue
	e events.Event[collections.ChangeArgs]
}

func (d *pseudoDic[TKey, TValue]) Get(key TKey) TValue {
	return d.m[key]
}

func (d *pseudoDic[TKey, TValue]) TryGet(key TKey) (TValue, bool) {
	v, ok := d.m[key]
	return v, ok
}

func (d *pseudoDic[TKey, TValue]) Enumerate() collections.Enumerator[collections.Tuple2[TKey, TValue]] {
	list := make([]collections.Tuple2[TKey, TValue], 0, len(d.m))
	for key, value := range d.m {
		list = append(list, tuple2.New(key, value))
	}
	return enumerator.Enumerate(list...)
}

func (d *pseudoDic[TKey, TValue]) Keys() collections.Enumerator[TKey] {
	list := make([]TKey, 0, len(d.m))
	for key := range d.m {
		list = append(list, key)
	}
	return enumerator.Enumerate(list...)
}

func (d *pseudoDic[TKey, TValue]) Values() collections.Enumerator[TValue] {
	list := make([]TValue, 0, len(d.m))
	for _, value := range d.m {
		list = append(list, value)
	}
	return enumerator.Enumerate(list...)
}

func (d *pseudoDic[TKey, TValue]) ToMap() map[TKey]TValue {
	return d.m
}

func (d *pseudoDic[TKey, TValue]) Empty() bool {
	return len(d.m) <= 0
}

func (d *pseudoDic[TKey, TValue]) Count() int {
	return len(d.m)
}

func (d *pseudoDic[TKey, TValue]) Contains(key TKey) bool {
	_, ok := d.m[key]
	return ok
}

func (d *pseudoDic[TKey, TValue]) String() string {
	return fmt.Sprintf(`%v`, d.m)
}

func (d *pseudoDic[TKey, TValue]) Equals(other any) bool {
	d2, ok := other.(collections.ReadonlyDictionary[TKey, TValue])
	return ok && reflect.DeepEqual(d.ToMap(), d2.ToMap())
}

func (d *pseudoDic[TKey, TValue]) OnChange() events.Event[collections.ChangeArgs] {
	return d.e
}

func Test_Dictionary_Readonly(t *testing.T) {
	d1 := &pseudoDic[int, int]{
		m: map[int]int{},
		e: event.New[collections.ChangeArgs](),
	}
	r1 := New(d1)
	check.Empty(t).Assert(r1)
	check.True(t).Assert(r1.Empty())

	d1.m[123] = 456
	check.Length(t, 1).Assert(r1)
	check.False(t).Assert(r1.Empty())
	check.Equal(t, 456).Assert(r1.Get(123))
	check.True(t).Assert(r1.Contains(123))
	check.False(t).Assert(r1.Contains(765))

	d1.m[222] = 333
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

	check.Same(t, d1.OnChange()).Assert(r1.OnChange())
}
