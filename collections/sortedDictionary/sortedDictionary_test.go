package sortedDictionary

import (
	"strings"
	"testing"

	"goToolbox/collections"
	"goToolbox/collections/enumerator"
	"goToolbox/collections/predicate"
	"goToolbox/collections/tuple2"
	"goToolbox/testers/check"
	"goToolbox/utils"
)

func validate[TKey comparable, TValue any](t *testing.T, dic collections.Dictionary[TKey, TValue]) {
	d, ok := dic.(*sortedImp[TKey, TValue])
	check.True(t).Name(`Convert to sortedImp for validation`).Assert(ok)
	check.NotNil(t).Assert(d.keys)
	check.NotNil(t).Assert(d.data)
	check.NotNil(t).Assert(d.comparer)
	check.Equal(t, d.keys).Assert(utils.SortedKeys(d.data, d.comparer))
}

func Test_SortedDictionary(t *testing.T) {
	d1 := New[int, int]()
	check.Empty(t).Assert(d1)
	check.True(t).Assert(d1.Empty())
	validate(t, d1)

	check.True(t).Assert(d1.Add(123, 321))
	validate(t, d1)
	check.False(t).Assert(d1.Add(123, 456))
	validate(t, d1)
	check.Length(t, 1).Assert(d1)
	check.False(t).Assert(d1.Empty())
	check.Equal(t, 456).Assert(d1.Get(123))
	check.True(t).Assert(d1.Contains(123))
	check.False(t).Assert(d1.Contains(765))

	check.False(t).Assert(d1.AddIfNotSet(123, 555))
	validate(t, d1)
	check.True(t).Assert(d1.AddIfNotSet(222, 333))
	validate(t, d1)
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
	validate(t, d2)
	check.String(t, "123: 456\n222: 333").Assert(d2)
	check.Equal(t, d1).Name(`d1.Equals(d2)`).Assert(d2)
	check.Equal(t, d2).Name(`d2.Equals(d1)`).Assert(d1)

	check.False(t).Assert(d1.Remove(833))
	validate(t, d1)
	check.True(t).Assert(d1.Remove(222))
	validate(t, d1)
	check.Length(t, 1).Assert(d1)
	check.NotEqual(t, d1).Name(`d1.Equals(d2)`).Assert(d2)
	check.NotEqual(t, d2).Name(`d2.Equals(d1)`).Assert(d1)

	check.True(t).Assert(d1.Add(833, 411))
	validate(t, d1)
	check.Length(t, 2).Assert(d1)
	check.NotEqual(t, d1).Name(`d1.Equals(d2)`).Assert(d2)
	check.NotEqual(t, d2).Name(`d2.Equals(d1)`).Assert(d1)
	d1.Clear()
	check.Empty(t).Assert(d1)
	validate(t, d1)

	check.Equal(t, `[123, 456]|[222, 333]`).Assert(d2.Enumerate().Join(`|`))
	check.Equal(t, `123|222`).Assert(d2.Keys().Join(`|`))
	check.Equal(t, `456|333`).Assert(d2.Values().Join(`|`))

	check.String(t, "123: 456\n222: 333").Assert(d2.Readonly())

	check.False(t).Assert(d2.RemoveIf(predicate.LessThan(10)))
	check.String(t, "123: 456\n222: 333").Assert(d2)
	validate(t, d2)
	check.False(t).Assert(d2.RemoveIf(nil))
	check.String(t, "123: 456\n222: 333").Assert(d2)
	validate(t, d2)
	check.True(t).Assert(d2.RemoveIf(predicate.GreaterThan(200)))
	check.String(t, "123: 456").Assert(d2)
	validate(t, d2)
}

func Test_SortedDictionary_FromAndWith(t *testing.T) {
	d1 := With(map[string]string{"One": "I", "Two": "II", "Three": "III"})
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

	check.True(t).Assert(d1.AddMap(map[string]string{"One": "Uno", "Three": "3", "Four": "0x04"}))
	check.String(t, "Four:  0x04\nOne:   Uno\nThree: 3\nTwo:   II").Assert(d1)

	check.False(t).Assert(d1.AddMap(map[string]string{"One": "Uno"}))
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

type pseudoComparable struct {
	name string
}

func (c *pseudoComparable) CompareTo(other *pseudoComparable) int {
	if c == nil {
		if other == nil {
			return 0
		}
		return -1
	}
	if other == nil {
		return 1
	}
	return strings.Compare(c.name, other.name)
}

func Test_SortedDictionary_New(t *testing.T) {
	d1 := New[int, string]().(*sortedImp[int, string])
	check.Empty(t).Assert(d1.keys)
	check.Zero(t).Assert(cap(d1.keys))

	d1 = CapNew[int, string](5).(*sortedImp[int, string])
	check.Empty(t).Assert(d1.keys)
	check.Equal(t, 5).Assert(cap(d1.keys))

	check.MatchError(t, `^invalid number of arguments \{count: 2, maximum: 1, usage: comparer\}$`).
		Panic(func() { New[int, string](utils.OrderedComparer[int](), utils.OrderedComparer[int]()) })

	d2 := CapNew[*pseudoComparable, int](7).(*sortedImp[*pseudoComparable, int])
	check.Empty(t).Assert(d2.keys)
	check.Equal(t, 7).Assert(cap(d2.keys))

	d2 = With[*pseudoComparable, int]((map[*pseudoComparable]int)(nil)).(*sortedImp[*pseudoComparable, int])
	check.Empty(t).Assert(d2.keys)
	check.Zero(t).Assert(cap(d2.keys))

	d3 := From[*pseudoComparable, int](d2.Enumerate()).(*sortedImp[*pseudoComparable, int])
	check.Empty(t).Assert(d3.keys)
	check.Zero(t).Assert(cap(d3.keys))

	d3 = CapFrom[*pseudoComparable, int](d2.Enumerate(), 6).(*sortedImp[*pseudoComparable, int])
	check.Empty(t).Assert(d3.keys)
	check.Equal(t, 6).Assert(cap(d3.keys))
}
