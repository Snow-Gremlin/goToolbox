package sortedDictionary

import (
	"bytes"
	"strings"
	"testing"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/collections/predicate"
	"github.com/Snow-Gremlin/goToolbox/collections/tuple2"
	"github.com/Snow-Gremlin/goToolbox/comp"
	"github.com/Snow-Gremlin/goToolbox/events/listener"
	"github.com/Snow-Gremlin/goToolbox/testers/check"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

func validate[TKey comparable, TValue any](t *testing.T, dic collections.Dictionary[TKey, TValue]) {
	d, ok := dic.(*sortedDictionaryImp[TKey, TValue])
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
	check.True(t).Assert(d1.Add(123, 456))
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

func Test_SortedDictionary_New(t *testing.T) {
	d1 := New[int, string]().(*sortedDictionaryImp[int, string])
	check.Empty(t).Assert(d1.keys)
	check.Zero(t).Assert(cap(d1.keys))

	d1 = CapNew[int, string](5).(*sortedDictionaryImp[int, string])
	check.Empty(t).Assert(d1.keys)
	check.Equal(t, 5).Assert(cap(d1.keys))

	check.MatchError(t, `^invalid number of arguments \{count: 2, maximum: 1, usage: comparer\}$`).
		Panic(func() { New[int, string](comp.Ordered[int](), comp.Ordered[int]()) })

	d2 := CapNew[*pseudoComparable, int](7).(*sortedDictionaryImp[*pseudoComparable, int])
	check.Empty(t).Assert(d2.keys)
	check.Equal(t, 7).Assert(cap(d2.keys))

	d2 = With[*pseudoComparable, int]((map[*pseudoComparable]int)(nil)).(*sortedDictionaryImp[*pseudoComparable, int])
	check.Empty(t).Assert(d2.keys)
	check.Zero(t).Assert(cap(d2.keys))

	d3 := From[*pseudoComparable, int](d2.Enumerate()).(*sortedDictionaryImp[*pseudoComparable, int])
	check.Empty(t).Assert(d3.keys)
	check.Zero(t).Assert(cap(d3.keys))

	d3 = CapFrom[*pseudoComparable, int](d2.Enumerate(), 6).(*sortedDictionaryImp[*pseudoComparable, int])
	check.Empty(t).Assert(d3.keys)
	check.Equal(t, 6).Assert(cap(d3.keys))
}

func Test_SortedDictionary_UnstableIteration(t *testing.T) {
	s := New[int, string]()
	s.Add(1, `one`)
	s.Add(2, `two`)
	s.Add(3, `three`)
	s.Add(4, `four`)
	s.Add(5, `five`)
	check.String(t, "1: one\n2: two\n3: three\n4: four\n5: five").Assert(s)

	it1 := s.Enumerate().Iterate()
	check.True(t).Assert(it1.Next())
	check.String(t, `[1, one]`).Assert(it1.Current())

	s.Add(8, `eight`)
	s.Add(9, `nine`)
	check.String(t, "1: one\n2: two\n3: three\n4: four\n5: five\n8: eight\n9: nine").Assert(s)

	check.True(t).Assert(it1.Next())
	check.String(t, `[2, two]`).Assert(it1.Current())

	s.Remove(2, 3, 4, 5)
	check.String(t, "1: one\n8: eight\n9: nine").Assert(s)

	check.True(t).Assert(it1.Next())
	check.String(t, `[9, nine]`).Assert(it1.Current())
	check.False(t).Assert(it1.Next())
}

func Test_SortedDictionary_OnChange(t *testing.T) {
	buf := &bytes.Buffer{}
	d := New[int, string]()
	lis := listener.New(func(args collections.ChangeArgs) {
		_, _ = buf.WriteString(args.Type().String())
	})
	defer lis.Cancel()
	check.True(t).Assert(lis.Subscribe(d.OnChange()))
	check.StringAndReset(t, ``).Assert(buf)

	check.True(t).Assert(d.Add(1, `one`))
	check.StringAndReset(t, `Added`).Assert(buf)
	check.False(t).Assert(d.Add(1, `one`))
	check.StringAndReset(t, ``).Assert(buf)
	check.True(t).Assert(d.Add(1, `uno`))
	check.StringAndReset(t, `Replaced`).Assert(buf)
	check.False(t).Assert(d.AddIfNotSet(1, `one`))
	check.StringAndReset(t, ``).Assert(buf)
	check.True(t).Assert(d.AddIfNotSet(2, `two`))
	check.StringAndReset(t, `Added`).Assert(buf)

	check.False(t).Assert(d.AddFrom(nil))
	check.StringAndReset(t, ``).Assert(buf)
	check.False(t).Assert(d.AddFrom(enumerator.Enumerate[collections.Tuple2[int, string]]()))
	check.StringAndReset(t, ``).Assert(buf)
	check.True(t).Assert(d.AddFrom(enumerator.Enumerate(tuple2.New(3, `three`))))
	check.StringAndReset(t, `Added`).Assert(buf)
	check.True(t).Assert(d.AddFrom(enumerator.Enumerate(tuple2.New(1, `one`), tuple2.New(4, `four`))))
	check.StringAndReset(t, `Replaced`).Assert(buf)
	check.False(t).Assert(d.AddIfNotSetFrom(enumerator.Enumerate(tuple2.New(1, `I`), tuple2.New(4, `IV`))))
	check.StringAndReset(t, ``).Assert(buf)
	check.True(t).Assert(d.AddIfNotSetFrom(enumerator.Enumerate(tuple2.New(5, `five`))))
	check.StringAndReset(t, `Added`).Assert(buf)

	check.False(t).Assert(d.AddMap(map[int]string{}))
	check.StringAndReset(t, ``).Assert(buf)
	check.False(t).Assert(d.AddMap(map[int]string{1: `one`}))
	check.StringAndReset(t, ``).Assert(buf)
	check.True(t).Assert(d.AddMap(map[int]string{6: `six`}))
	check.StringAndReset(t, `Added`).Assert(buf)
	check.True(t).Assert(d.AddMap(map[int]string{6: `VI`}))
	check.StringAndReset(t, `Replaced`).Assert(buf)
	check.False(t).Assert(d.AddMapIfNotSet(map[int]string{6: `six`}))
	check.StringAndReset(t, ``).Assert(buf)
	check.True(t).Assert(d.AddMapIfNotSet(map[int]string{7: `VII`}))
	check.StringAndReset(t, `Added`).Assert(buf)

	check.True(t).Assert(d.Remove(4, 7))
	check.StringAndReset(t, `Removed`).Assert(buf)
	check.False(t).Assert(d.Remove(4, 7))
	check.StringAndReset(t, ``).Assert(buf)
	check.False(t).Assert(d.RemoveIf(nil))
	check.StringAndReset(t, ``).Assert(buf)
	check.True(t).Assert(d.RemoveIf(predicate.GreaterThan(4)))
	check.StringAndReset(t, `Removed`).Assert(buf)
	check.False(t).Assert(d.RemoveIf(predicate.GreaterThan(4)))
	check.StringAndReset(t, ``).Assert(buf)

	d.Clear()
	check.StringAndReset(t, `Removed`).Assert(buf)
	d.Clear()
	check.StringAndReset(t, ``).Assert(buf)
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
