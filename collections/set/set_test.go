package set

import (
	"slices"
	"testing"

	"goToolbox/collections/enumerator"
	"goToolbox/collections/predicate"
	"goToolbox/testers/check"
)

func Test_Set(t *testing.T) {
	s := With(1, 2, 3)
	check.Length(t, 3).Assert(s)
	check.String(t, `1, 2, 3`).Assert(s)
	check.False(t).Assert(s.Empty())

	p := s.ToSlice()
	slices.Sort(p)
	check.Equal(t, []int{1, 2, 3}).Assert(p)
	check.Length(t, 3).Assert(s.ToList())

	p = make([]int, 1)
	s.CopyToSlice(p) // Didn't panic

	p = make([]int, 5)
	s.CopyToSlice(p)
	slices.Sort(p)
	check.Equal(t, []int{0, 0, 1, 2, 3}).Assert(p)

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
	check.String(t, `10, 11, 5, 6, 7, 8, 9`).Assert(s)
}

func Test_Set_New(t *testing.T) {
	s := New[int]()
	check.Empty(t).Assert(s)

	s = New[int](10)
	check.Empty(t).Assert(s)

	check.MatchError(t, `^invalid number of arguments \{count: 2, maximum: 1, usage: capacity\}$`).
		Panic(func() { New[int](15, 5) })

	s = With(1, 2, 3)
	check.Length(t, 3).Assert(s)
	check.String(t, `1, 2, 3`).Assert(s)

	s = From[int](nil)
	check.Empty(t).Assert(s)

	s = From[int](nil, 10)
	check.Empty(t).Assert(s)

	check.MatchError(t, `^invalid number of arguments \{count: 2, maximum: 1, usage: capacity\}$`).
		Panic(func() { From[int](nil, 1, 5) })

	s = From[int](enumerator.Range(1, 5), 10)
	check.Length(t, 5).Assert(s)
	check.String(t, `1, 2, 3, 4, 5`).Assert(s)
}

func Test_Set_UnstableIteration(t *testing.T) {
	s := With(1, 2, 3)
	it := s.Enumerate().Iterate()

	check.True(t).Assert(it.Next())
	value1 := it.Current()
	values := []int{value1}

	check.True(t).Assert(s.Add(4))
	check.True(t).Assert(it.Next())
	values = append(values, it.Current())

	check.True(t).Assert(s.Remove(value1))
	check.True(t).Assert(it.Next())
	values = append(values, it.Current())

	// Didn't pick up Four because it was added after the keys were captured.
	check.False(t).Assert(it.Next())
	check.SameElems(t, []int{1, 2, 3}).Assert(values)
}
