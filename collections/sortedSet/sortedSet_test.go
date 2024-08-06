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

	check.True(t).Assert(s.Remove(4, 5))
	check.False(t).Assert(s.Remove(4, 5))
	check.String(t, `1, 2, 3`).Assert(s)

	check.True(t).Assert(s.Add(4, 5, 6, 7, 8))
	check.False(t).Assert(s.RemoveIf(predicate.IsZero[int]()))
	check.True(t).Assert(s.RemoveIf(predicate.LessThan(5)))
	check.String(t, `5, 6, 7, 8`).Assert(s)

	check.False(t).Assert(s.AddFrom(nil))
	check.False(t).Assert(s.AddFrom(enumerator.Range(5, 3)))
	check.True(t).Assert(s.AddFrom(enumerator.Range(9, 3)))
	check.String(t, `5, 6, 7, 8, 9, 10, 11`).Assert(s)
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

	//check.MatchError(t, `^invalid number of arguments \{count: 2, maximum: 1, usage: capacity\}$`).
	//	Panic(func() { CapNew[int](15, 5) })

	s = With([]int{1, 2, 3})
	check.Length(t, 3).Assert(s)
	check.String(t, `1, 2, 3`).Assert(s)

	s = From[int](nil)
	check.Empty(t).Assert(s)

	s = CapFrom[int](nil, 10)
	check.Empty(t).Assert(s)

	// TODO: Custom compare

	//check.MatchError(t, `^invalid number of arguments \{count: 2, maximum: 1, usage: capacity\}$`).
	//	Panic(func() { From[int](nil, 1, 5) })

	s = CapFrom(enumerator.Range(1, 5), 10)
	check.Length(t, 5).Assert(s)
	check.String(t, `1, 2, 3, 4, 5`).Assert(s)
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
