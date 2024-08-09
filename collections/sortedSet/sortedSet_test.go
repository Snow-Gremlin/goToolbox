package sortedSet

import (
	"bytes"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/collections/predicate"
	"github.com/Snow-Gremlin/goToolbox/comp"
	"github.com/Snow-Gremlin/goToolbox/events/listener"
	"github.com/Snow-Gremlin/goToolbox/testers/check"
)

func Test_SortedSet(t *testing.T) {
	s := With([]int{1, 2, 3})
	check.Length(t, 3).Assert(s)
	check.String(t, `1, 2, 3`).Assert(s)
	check.False(t).Assert(s.Empty())

	p := s.ToSlice()
	check.Equal(t, []int{1, 2, 3}).Assert(p)
	check.Length(t, 3).Assert(s.ToList())

	p = make([]int, 1)
	s.CopyToSlice(p) // Didn't panic
	check.Equal(t, []int{1}).Assert(p)

	p = make([]int, 5)
	s.CopyToSlice(p)
	check.Equal(t, []int{1, 2, 3, 0, 0}).Assert(p)

	check.True(t).Assert(s.Contains(1))
	check.False(t).Assert(s.Contains(4))

	check.False(t).Assert(s.Add(1, 2))
	check.True(t).Assert(s.Add(3, 5))
	check.String(t, `1, 2, 3, 5`).Assert(s)
	check.Length(t, 4).Assert(s)

	check.String(t, `1, 2, 3, 5`).Assert(s.Readonly())

	s2 := s.Clone()
	check.Equal(t, s2).Assert(s)
	check.String(t, `1, 2, 3, 5`).Assert(s2)

	check.True(t).Assert(s2.Add(4))
	check.True(t).Assert(s2.Remove(5))
	check.String(t, `1, 2, 3, 4`).Assert(s2)
	check.NotEqual(t, s2).Assert(s)

	s2.Clear()
	check.Empty(t).Assert(s2)
	check.True(t).Assert(s2.Empty())
	check.String(t, ``).Assert(s2)
	check.NotEqual(t, s2).Assert(s)
	check.MatchError(t, `^collection contains no values \{action: First\}$`).Panic(func() { s2.First() })
	check.MatchError(t, `^collection contains no values \{action: Last\}$`).Panic(func() { s2.Last() })
	check.MatchError(t, `^collection contains no values \{action: TakeFirst\}$`).Panic(func() { s2.TakeFirst() })
	check.MatchError(t, `^collection contains no values \{action: TakeLast\}$`).Panic(func() { s2.TakeLast() })

	check.True(t).Assert(s.Remove(4, 5))
	check.False(t).Assert(s.Remove(4, 5))
	check.String(t, `1, 2, 3`).Assert(s)

	check.True(t).Assert(s.Add(4, 5, 6, 7, 8))
	check.False(t).Assert(s.RemoveIf(nil)) // no effect
	check.False(t).Assert(s.RemoveIf(predicate.IsZero[int]()))
	check.True(t).Assert(s.RemoveIf(predicate.LessThan(5)))
	check.String(t, `5, 6, 7, 8`).Assert(s)

	check.False(t).Assert(s.AddFrom(nil))
	check.False(t).Assert(s.AddFrom(enumerator.Range(5, 3)))
	check.True(t).Assert(s.AddFrom(enumerator.Range(9, 3)))
	check.String(t, `5, 6, 7, 8, 9, 10, 11`).Assert(s)
	s.RemoveRange(3, 0) // no effect
	s.RemoveRange(3, 2)
	check.String(t, `5, 6, 7, 10, 11`).Assert(s)

	check.Equal(t, 5).Assert(s.Get(0))
	check.Equal(t, 6).Assert(s.Get(1))
	check.Equal(t, 7).Assert(s.Get(2))
	check.Equal(t, 10).Assert(s.Get(3))
	check.Equal(t, 11).Assert(s.Get(4))
	check.MatchError(t, `^index out of bounds \{count: 5, index: -1\}$`).Panic(func() { s.Get(-1) })
	check.MatchError(t, `^index out of bounds \{count: 5, index: 5\}$`).Panic(func() { s.Get(5) })

	v, ok := s.TryGet(2)
	check.True(t).Assert(ok)
	check.Equal(t, 7).Assert(v)

	v, ok = s.TryGet(-1)
	check.False(t).Assert(ok)
	check.Zero(t).Assert(v)

	check.Equal(t, 5).Assert(s.First())
	check.Equal(t, 11).Assert(s.Last())

	check.Equal(t, -1).Assert(s.IndexOf(4))
	check.Equal(t, 0).Assert(s.IndexOf(5))
	check.Equal(t, 1).Assert(s.IndexOf(6))
	check.Equal(t, 2).Assert(s.IndexOf(7))
	check.Equal(t, -1).Assert(s.IndexOf(8))
	check.Equal(t, -1).Assert(s.IndexOf(9))
	check.Equal(t, 3).Assert(s.IndexOf(10))
	check.Equal(t, 4).Assert(s.IndexOf(11))
	check.Equal(t, -1).Assert(s.IndexOf(12))

	s3 := s.Clone()
	check.Equal(t, 5).Assert(s3.TakeFirst())
	check.String(t, `6, 7, 10, 11`).Assert(s3)
	check.Equal(t, 11).Assert(s3.TakeLast())
	check.String(t, `6, 7, 10`).Assert(s3)

	s3 = s.Clone()
	check.String(t, ``).Assert(s3.TakeFront(0))
	check.String(t, `5, 6`).Assert(s3.TakeFront(2))
	check.String(t, `7, 10, 11`).Assert(s3)
	check.String(t, ``).Assert(s3.TakeBack(0))
	check.String(t, `10, 11`).Assert(s3.TakeBack(2))
	check.String(t, `7`).Assert(s3)
}

func Test_SortedSet_CustomCompare(t *testing.T) {
	revStr := func(v int) string {
		digits := []byte(strconv.Itoa(v))
		slices.Reverse(digits)
		return string(digits)
	}
	s := New(func(x, y int) int {
		return strings.Compare(revStr(x), revStr(y))
	})
	s.Add(48, 22, 123, 43, 33, 2, 20, 25)
	check.String(t, `20, 2, 22, 123, 33, 43, 25, 48`).Assert(s)

	s2 := s.Clone()
	s2.Add(1, 2, 3, 4, 5, 6, 10, 30)
	check.String(t, `10, 20, 30, 1, 2, 22, 3, 123, 33, 43, 4, 5, 25, 6, 48`).Assert(s2)

	check.True(t).Assert(s2.Contains(22))
	check.True(t).Assert(s2.Contains(123))
	check.True(t).Assert(s2.Contains(4))
	check.False(t).Assert(s2.Contains(52))
}

func Test_SortedSet_New(t *testing.T) {
	s := New[int]()
	check.Empty(t).Assert(s)

	s = CapNew[int](10)
	check.Empty(t).Assert(s)

	check.MatchError(t, `^invalid number of arguments \{count: 2, maximum: 1, usage: comparer\}$`).
		Panic(func() { From(nil, comp.Ordered[int](), comp.Ordered[int]()) })

	s = With([]int{1, 2, 3})
	check.Length(t, 3).Assert(s)
	check.String(t, `1, 2, 3`).Assert(s)

	s = From[int](nil)
	check.Empty(t).Assert(s)

	s = CapFrom[int](nil, 10)
	check.Empty(t).Assert(s)

	s = CapFrom(enumerator.Range(1, 5), 10)
	check.Length(t, 5).Assert(s)
	check.String(t, `1, 2, 3, 4, 5`).Assert(s)
}

func Test_SortedSet_TryAddAndOverwrite(t *testing.T) {
	type person struct {
		first, last string
	}

	compareOnlyLast := func(p, q person) int {
		return strings.Compare(p.last, q.last)
	}

	s := New(compareOnlyLast)
	v, ok := s.TryAdd(person{first: `Jill`, last: `Smith`})
	check.True(t).Assert(ok)
	check.String(t, `{Jill Smith}`).Assert(v)
	check.String(t, `{Jill Smith}`).Assert(s)

	v, ok = s.TryAdd(person{first: `Jill`, last: `Johnson`})
	check.True(t).Assert(ok)
	check.String(t, `{Jill Johnson}`).Assert(v)
	check.String(t, `{Jill Johnson}, {Jill Smith}`).Assert(s)

	// "Smith" already exists so don't overwrite and return.
	v, ok = s.TryAdd(person{first: `Tom`, last: `Smith`})
	check.False(t).Assert(ok)
	check.String(t, `{Jill Smith}`).Assert(v)
	check.String(t, `{Jill Johnson}, {Jill Smith}`).Assert(s)

	// Try to add but don't replace the original.
	ok = s.Add(person{first: `Tom`, last: `Smith`})
	check.False(t).Assert(ok)
	check.String(t, `{Jill Johnson}, {Jill Smith}`).Assert(s)

	// Try again but overwrite this time.
	ok = s.Overwrite(person{first: `Tom`, last: `Smith`})
	check.False(t).Assert(ok)
	check.String(t, `{Jill Johnson}, {Tom Smith}`).Assert(s)

	ok = s.Overwrite(person{first: `Bill`, last: `Wolf`})
	check.True(t).Assert(ok)
	check.String(t, `{Jill Johnson}, {Tom Smith}, {Bill Wolf}`).Assert(s)

	ok = s.OverwriteFrom(enumerator.Enumerate(
		person{first: `Mark`, last: `Wolf`},
		person{first: `Mark`, last: `Smith`},
		person{first: `Mark`, last: `Gram`}))
	check.True(t).Assert(ok)
	check.String(t, `{Mark Gram}, {Jill Johnson}, {Mark Smith}, {Mark Wolf}`).Assert(s)
}

func Test_SortedSet_UnstableIteration(t *testing.T) {
	s := With([]int{2, 4, 6})
	it := s.Enumerate().Iterate()

	check.True(t).Assert(it.Next())
	check.Equal(t, 2).Assert(it.Current())

	check.True(t).Assert(it.Next())
	check.Equal(t, 4).Assert(it.Current())

	check.True(t).Assert(s.Add(3))
	check.True(t).Assert(it.Next()) // repeat 4 since 3 inserted before it
	check.Equal(t, 4).Assert(it.Current())

	check.True(t).Assert(s.Remove(2, 3, 4)) // removes everything but 6
	check.False(t).Assert(it.Next())
	check.Zero(t).Assert(it.Current())
}

func Test_SortedSet_UnstableIteration_Backwards(t *testing.T) {
	s := With([]int{2, 4, 6})
	it := s.Backwards().Iterate()

	check.True(t).Assert(it.Next())
	check.Equal(t, 6).Assert(it.Current())

	check.True(t).Assert(it.Next())
	check.Equal(t, 4).Assert(it.Current())

	check.True(t).Assert(s.Add(3))
	check.True(t).Assert(it.Next()) // skip 3 and goto 2 since 3 inserted where 4 was
	check.Equal(t, 2).Assert(it.Current())

	check.True(t).Assert(s.Remove(2, 3, 4)) // removes everything but 6
	check.False(t).Assert(it.Next())
	check.Zero(t).Assert(it.Current())
}

func Test_SortedSet_OnChange(t *testing.T) {
	buf := &bytes.Buffer{}
	s := New[int]()
	lis := listener.New(func(args collections.ChangeArgs) {
		_, _ = buf.WriteString(args.Type().String())
	})
	defer lis.Cancel()
	check.True(t).Assert(lis.Subscribe(s.OnChange()))
	check.StringAndReset(t, ``).Assert(buf)

	check.False(t).Assert(s.Add())
	check.StringAndReset(t, ``).Assert(buf)
	check.True(t).Assert(s.Add(1, 5))
	check.StringAndReset(t, `Added`).Assert(buf)
	check.False(t).Assert(s.Add(1, 5))
	check.StringAndReset(t, ``).Assert(buf)
	check.False(t).Assert(s.AddFrom(nil))
	check.StringAndReset(t, ``).Assert(buf)
	check.False(t).Assert(s.AddFrom(enumerator.Enumerate[int]()))
	check.StringAndReset(t, ``).Assert(buf)
	check.False(t).Assert(s.AddFrom(enumerator.Enumerate(1, 5)))
	check.StringAndReset(t, ``).Assert(buf)
	check.True(t).Assert(s.AddFrom(enumerator.Enumerate(3, 2)))
	check.StringAndReset(t, `Added`).Assert(buf)
	check.String(t, `1, 2, 3, 5`).Assert(s)

	check.False(t).Assert(s.Remove())
	check.StringAndReset(t, ``).Assert(buf)
	check.False(t).Assert(s.Remove(4, 6))
	check.StringAndReset(t, ``).Assert(buf)
	check.True(t).Assert(s.Remove(2))
	check.StringAndReset(t, `Removed`).Assert(buf)
	check.False(t).Assert(s.RemoveIf(nil))
	check.StringAndReset(t, ``).Assert(buf)
	check.True(t).Assert(s.RemoveIf(predicate.GreaterEq(3)))
	check.StringAndReset(t, `Removed`).Assert(buf)
	check.False(t).Assert(s.RemoveIf(predicate.GreaterEq(3)))
	check.StringAndReset(t, ``).Assert(buf)
	check.String(t, `1`).Assert(s)

	s.Clear()
	check.StringAndReset(t, `Removed`).Assert(buf)
	s.Clear()
	check.StringAndReset(t, ``).Assert(buf)
}
