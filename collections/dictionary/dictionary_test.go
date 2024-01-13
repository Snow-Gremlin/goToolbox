package dictionary

import (
	"strings"
	"testing"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/collections/predicate"
	"github.com/Snow-Gremlin/goToolbox/collections/tuple2"
	"github.com/Snow-Gremlin/goToolbox/testers/check"
)

func Test_Dictionary(t *testing.T) {
	d1 := New[int, int]()
	check.Empty(t).Assert(d1)
	check.True(t).Assert(d1.Empty())

	check.True(t).Assert(d1.Add(123, 321))
	check.False(t).Assert(d1.Add(123, 456))
	check.Length(t, 1).Assert(d1)
	check.False(t).Assert(d1.Empty())
	check.Equal(t, 456).Assert(d1.Get(123))
	check.True(t).Assert(d1.Contains(123))
	check.False(t).Assert(d1.Contains(765))

	check.False(t).Assert(d1.AddIfNotSet(123, 555))
	check.True(t).Assert(d1.AddIfNotSet(222, 333))
	check.Length(t, 2).Assert(d1)
	check.Equal(t, 456).Assert(d1.Get(123))
	check.Equal(t, 333).Assert(d1.Get(222))

	v, ok := d1.TryGet(222)
	check.Equal(t, 333).Assert(v)
	check.True(t).Assert(ok)

	v, ok = d1.TryGet(251)
	check.Zero(t).Assert(v)
	check.False(t).Assert(ok)

	check.String(t, "123: 456\n222: 333").Assert(d1)
	check.Equal(t, map[int]int{123: 456, 222: 333}).Assert(d1.ToMap())
	d2 := d1.Clone()
	check.String(t, "123: 456\n222: 333").Assert(d2)
	check.Equal(t, d1).Name(`d1.Equals(d2)`).Assert(d2)
	check.Equal(t, d2).Name(`d2.Equals(d1)`).Assert(d1)

	check.False(t).Assert(d1.Remove(833))
	check.True(t).Assert(d1.Remove(222))
	check.Length(t, 1).Assert(d1)
	check.NotEqual(t, d1).Name(`d1.Equals(d2)`).Assert(d2)
	check.NotEqual(t, d2).Name(`d2.Equals(d1)`).Assert(d1)

	check.True(t).Assert(d1.Add(833, 411))
	check.Length(t, 2).Assert(d1)
	check.NotEqual(t, d1).Name(`d1.Equals(d2)`).Assert(d2)
	check.NotEqual(t, d2).Name(`d2.Equals(d1)`).Assert(d1)
	d1.Clear()
	check.Empty(t).Assert(d1)

	check.Equal(t, `[123, 456]|[222, 333]`).Assert(d2.Enumerate().Strings().Sort().Join(`|`))
	check.Equal(t, `123|222`).Assert(d2.Keys().Sort().Join(`|`))
	check.Equal(t, `333|456`).Assert(d2.Values().Sort().Join(`|`))

	check.String(t, "123: 456\n222: 333").Assert(d2.Readonly())

	check.False(t).Assert(d2.RemoveIf(predicate.LessThan(10)))
	check.String(t, "123: 456\n222: 333").Assert(d2)
	check.False(t).Assert(d2.RemoveIf(nil))
	check.String(t, "123: 456\n222: 333").Assert(d2)
	check.True(t).Assert(d2.RemoveIf(predicate.GreaterThan(200)))
	check.String(t, "123: 456").Assert(d2)
}

func Test_Dictionary_FromAndMap(t *testing.T) {
	d1 := With(map[string]string{`One`: `I`, `Two`: `II`, `Three`: `III`})
	check.String(t, "One:   I\nThree: III\nTwo:   II").Assert(d1)

	d2 := From(d1.Enumerate().
		Where(func(t collections.Tuple2[string, string]) bool {
			return !strings.Contains(t.Value1(), `e`)
		}))
	check.String(t, `Two: II`).Assert(d2)

	d3 := From(enumerator.Select(d1.Enumerate(),
		func(t collections.Tuple2[string, string]) collections.Tuple2[string, int] {
			return tuple2.New(t.Value1(), len(t.Value1()))
		}))
	check.String(t, "One:   3\nThree: 5\nTwo:   3").Assert(d3)

	check.True(t).Assert(d1.AddMap(map[string]string{`One`: `Uno`, `Three`: `3`, `Four`: `0x04`}))
	check.String(t, "Four:  0x04\nOne:   Uno\nThree: 3\nTwo:   II").Assert(d1)

	check.False(t).Assert(d1.AddMap(map[string]string{`One`: `Uno`}))
	check.String(t, "Four:  0x04\nOne:   Uno\nThree: 3\nTwo:   II").Assert(d1)

	check.False(t).Assert(d1.AddMap(nil))
	check.String(t, "Four:  0x04\nOne:   Uno\nThree: 3\nTwo:   II").Assert(d1)

	d4 := From[string, string](nil)
	check.String(t, ``).Assert(d4)

	check.True(t).Assert(d4.AddMapIfNotSet(map[string]string{`ij`: `k`, `jk`: `i`, `ki`: `j`}))
	check.False(t).Assert(d4.AddMapIfNotSet(map[string]string{`ij`: `-ji`, `jk`: `-kj`, `ki`: `-ik`}))
	check.String(t, "ij: k\njk: i\nki: j").Assert(d4)

	d5 := With(map[string]string{`ij`: `-ji`, `jk`: `-kj`, `ki`: `-ik`, `ijk`: `-1`})
	check.True(t).Assert(d4.AddIfNotSetFrom(d5.Enumerate()))
	check.String(t, "ij:  k\nijk: -1\njk:  i\nki:  j").Assert(d4)
}

func Test_Dictionary_Capacity(t *testing.T) {
	d1 := New[string, int]()
	check.Empty(t).Assert(d1)

	d2 := New[string, int](10)
	check.Empty(t).Assert(d2)

	check.MatchError(t, `^invalid number of arguments \{count: 2, maximum: 1, usage: capacity\}$`).
		Panic(func() { New[string, int](1, 3) })

	d2 = New[string, int](-1)
	check.Empty(t).Assert(d2)
}

func Test_Dictionary_UnstableIteration(t *testing.T) {
	d := With(map[string]string{`One`: `I`, `Two`: `II`, `Three`: `III`})
	it := d.Enumerate().Iterate()

	check.True(t).Assert(it.Next())
	key1 := it.Current().Value1()
	keys := []string{key1}

	check.True(t).Assert(d.Add(`Four`, `IV`))
	check.True(t).Assert(it.Next())
	keys = append(keys, it.Current().Value1())

	check.True(t).Assert(d.Remove(key1))
	check.True(t).Assert(it.Next())
	keys = append(keys, it.Current().Value1())

	// Didn't pick up Four because it was added after the keys were captured.
	check.False(t).Assert(it.Next())
	check.SameElems(t, []string{`One`, `Two`, `Three`}).Assert(keys)
}
