package list

import (
	"testing"

	"goToolbox/collections/enumerator"
	"goToolbox/collections/predicate"
	"goToolbox/testers/check"
)

func Test_List(t *testing.T) {
	s := With(1, 2, 3, 4, 5)
	check.String(t, `1, 2, 3, 4, 5`).
		Assert(s).
		Assert(s.Readonly())
	check.Length(t, 5).Assert(s)
	check.False(t).Assert(s.Empty())
	check.Equal(t, []int{1, 2, 3, 4, 5}).
		Assert(s.ToSlice()).
		Assert(s.Enumerate().ToSlice())
	check.Equal(t, []int{5, 4, 3, 2, 1}).Assert(s.Backwards().ToSlice())

	sc := make([]int, 3)
	s.CopyToSlice(sc)
	check.Equal(t, []int{1, 2, 3}).Assert(sc)

	sc = make([]int, 8)
	s.CopyToSlice(sc)
	check.Equal(t, []int{1, 2, 3, 4, 5, 0, 0, 0}).Assert(sc)

	check.Equal(t, 3).Assert(s.IndexOf(4))
	check.Equal(t, -1).Assert(s.IndexOf(6))
	check.True(t).Assert(s.Contains(3))
	check.False(t).Assert(s.Contains(-1))
	check.Equal(t, 1).Assert(s.First())
	check.Equal(t, 5).Assert(s.Last())
	check.Equal(t, 1).Assert(s.Get(0))
	check.Equal(t, 2).Assert(s.Get(1))
	check.Equal(t, 3).Assert(s.Get(2))
	check.Equal(t, 4).Assert(s.Get(3))
	check.Equal(t, 5).Assert(s.Get(4))
	check.MatchError(t, `^index out of bounds \{count: 5, index: -1\}$`).Panic(func() { s.Get(-1) })
	check.MatchError(t, `^index out of bounds \{count: 5, index: 5\}$`).Panic(func() { s.Get(5) })

	v, ok := s.TryGet(2)
	check.True(t).Assert(ok)
	check.Equal(t, 3).Assert(v)

	v, ok = s.TryGet(-1)
	check.False(t).Assert(ok)
	check.Zero(t).Assert(v)

	check.True(t).Withf(`[%s].StartsWith([1, 2, 3])`, s.String()).Assert(s.StartsWith(With(1, 2, 3)))
	check.False(t).Withf(`[%s].StartsWith([1, 2, 4])`, s.String()).Assert(s.StartsWith(With(1, 2, 4)))
	check.True(t).Withf(`[%s].EndsWith([3, 4, 5])`, s.String()).Assert(s.EndsWith(With(3, 4, 5)))
	check.False(t).Withf(`[%s].EndsWith([3, 2, 5])`, s.String()).Assert(s.EndsWith(With(3, 2, 5)))

	s2 := s.Clone()
	check.Equal(t, s).Assert(s2)
	check.Equal(t, s2).Assert(s)
	check.String(t, `1, 2, 3, 4, 5`).Assert(s2)

	s2.Set(2, 11)
	check.String(t, `1, 2, 11, 4, 5`).Assert(s2)
	s2.Set(2, 24, 42)
	check.String(t, `1, 2, 24, 42, 5`).Assert(s2)
	s2.Set(2) // no effect
	check.String(t, `1, 2, 24, 42, 5`).Assert(s2)
	check.NotEqual(t, s).Assert(s2)
	check.MatchError(t, `^index out of bounds \{count: 5, index: -1\}$`).Panic(func() { s2.Set(-1, 45) })
	check.MatchError(t, `^index out of bounds \{count: 5, index: 6\}$`).Panic(func() { s2.Set(6, 45) })

	s2.Remove(1, 3)
	check.String(t, `1, 5`).Assert(s2)
	check.NotEqual(t, s).Assert(s2)

	s2.Clear()
	check.String(t, ``).Assert(s2)
	check.NotEqual(t, s).Assert(s2)
	check.MatchError(t, `^collection contains no values \{action: First\}$`).Panic(func() { s2.First() })
	check.MatchError(t, `^collection contains no values \{action: Last\}$`).Panic(func() { s2.Last() })
	check.MatchError(t, `^collection contains no values \{action: TakeFirst\}$`).Panic(func() { s2.TakeFirst() })
	check.MatchError(t, `^collection contains no values \{action: TakeLast\}$`).Panic(func() { s2.TakeLast() })

	s2.Prepend(1)
	s2.Append(2)
	s2.Prepend(3)
	check.String(t, `3, 1, 2`).Assert(s2)
	check.Equal(t, 3).Assert(s2.TakeFirst())
	check.String(t, `1, 2`).Assert(s2)
	check.Equal(t, 2).Assert(s2.TakeLast())
	check.String(t, `1`).Assert(s2)
	s2.Append(2, 3)
	check.String(t, `1, 2, 3`).Assert(s2)
	s2.Insert(1, 7, 2, 8, 9)
	check.String(t, `1, 7, 2, 8, 9, 2, 3`).Assert(s2)
	check.Equal(t, 2).Assert(s2.IndexOf(2, -10))
	check.Equal(t, 2).Assert(s2.IndexOf(2, 1))
	check.Equal(t, 5).Assert(s2.IndexOf(2, 2))
	check.Equal(t, -1).Assert(s2.IndexOf(2, 5))
	check.False(t).Assert(s2.RemoveIf(predicate.LessThan(0)))
	check.String(t, `1, 7, 2, 8, 9, 2, 3`).Assert(s2)
	check.True(t).Assert(s2.RemoveIf(predicate.LessThan(4)))
	check.String(t, `7, 8, 9`).Assert(s2)

	s3 := From(s.Enumerate().Strings())
	check.Equal(t, []string{`1`, `2`, `3`, `4`, `5`}).Assert(s3.ToSlice())
}

func Test_List_New(t *testing.T) {
	s := New[int]().(*listImp[int])
	check.Zero(t).Assert(len(s.s))
	check.Zero(t).Assert(cap(s.s))
	check.String(t, ``).Assert(s)

	s = New[int](5).(*listImp[int])
	check.Equal(t, 5).Assert(len(s.s))
	check.Equal(t, 5).Assert(cap(s.s))
	check.String(t, `0, 0, 0, 0, 0`).Assert(s)

	s = New[int](5, 10).(*listImp[int])
	check.Equal(t, 5).Assert(len(s.s))
	check.Equal(t, 10).Assert(cap(s.s))
	check.String(t, `0, 0, 0, 0, 0`).Assert(s)

	s = New[int](5, 2).(*listImp[int])
	check.Equal(t, 5).Assert(len(s.s))
	check.Equal(t, 5).Assert(cap(s.s))
	check.String(t, `0, 0, 0, 0, 0`).Assert(s)

	check.MatchError(t, `^invalid number of arguments \{count: 3, maximum: 2, usage: size and capacity\}$`).
		Panic(func() { New[int](1, 2, 3) })

	s = New[int](-1).(*listImp[int])
	check.Zero(t).Assert(len(s.s))
	check.Zero(t).Assert(cap(s.s))
	check.String(t, ``).Assert(s)

	s = New[int](-1, -1).(*listImp[int])
	check.Zero(t).Assert(len(s.s))
	check.Zero(t).Assert(cap(s.s))
	check.String(t, ``).Assert(s)
}

func Test_List_Fill(t *testing.T) {
	s := Fill(8, 3).(*listImp[int])
	check.Equal(t, 3).Assert(len(s.s))
	check.Equal(t, 3).Assert(cap(s.s))
	check.String(t, `8, 8, 8`).Assert(s)

	s = Fill(8, 3, 6).(*listImp[int])
	check.Equal(t, 3).Assert(len(s.s))
	check.Equal(t, 6).Assert(cap(s.s))
	check.String(t, `8, 8, 8`).Assert(s)

	check.MatchError(t, `^invalid number of arguments \{count: 2, maximum: 1, usage: capacity\}$`).
		Panic(func() { Fill[int](8, 2, 3, 4) })

	s = Fill(8, -1).(*listImp[int])
	check.Zero(t).Assert(len(s.s))
	check.Zero(t).Assert(cap(s.s))
	check.String(t, ``).Assert(s)

	s = Fill(8, -1, -1).(*listImp[int])
	check.Zero(t).Assert(len(s.s))
	check.Zero(t).Assert(cap(s.s))
	check.String(t, ``).Assert(s)
}

func Test_List_MoreMethods(t *testing.T) {
	s1 := From(enumerator.Range(0, 3))
	check.Length(t, 3).Assert(s1)

	s1.PrependFrom(nil)
	s1.AppendFrom(nil)
	check.Length(t, 3).Assert(s1)
	check.String(t, `0, 1, 2`).Assert(s1)

	s1.PrependFrom(enumerator.Range(5, 3))
	check.Length(t, 6).Assert(s1)
	check.String(t, `5, 6, 7, 0, 1, 2`).Assert(s1)

	s2 := s1.TakeFront(2)
	check.Length(t, 2).Assert(s2)
	check.String(t, `5, 6`).Assert(s2)
	check.Length(t, 4).Assert(s1)
	check.String(t, `7, 0, 1, 2`).Assert(s1)

	s3 := s1.TakeBack(2)
	check.Length(t, 2).Assert(s3)
	check.String(t, `1, 2`).Assert(s3)
	check.Length(t, 2).Assert(s1)
	check.String(t, `7, 0`).Assert(s1)

	s4 := s1.TakeFront(12)
	check.Length(t, 2).Assert(s4)
	check.String(t, `7, 0`).Assert(s4)
	check.Empty(t).Assert(s1)
	check.String(t, ``).Assert(s1)
	s4 = s1.TakeFront(12)
	check.Empty(t).Assert(s4)

	s1.InsertFrom(0, enumerator.Enumerate(9, 8))
	check.Length(t, 2).Assert(s1)
	check.String(t, `9, 8`).Assert(s1)

	s5 := s1.TakeBack(12)
	check.Length(t, 2).Assert(s5)
	check.String(t, `9, 8`).Assert(s5)
	check.Empty(t).Assert(s1)
	check.String(t, ``).Assert(s1)
	s5 = s1.TakeBack(12)
	check.Empty(t).Assert(s5)

	s1.Append(9, 7, 5)
	check.String(t, `9, 7, 5`).Assert(s1)
	s1.SetFrom(1, enumerator.Enumerate(77, 55, 33, 11))
	check.Length(t, 5).Assert(s1)
	check.String(t, `9, 77, 55, 33, 11`).Assert(s1)
}

func Test_List_UnstableIteration(t *testing.T) {
	s := From(enumerator.Range(1, 5))
	check.String(t, `1, 2, 3, 4, 5`).Assert(s)

	it1 := s.Enumerate().Iterate()
	check.True(t).Assert(it1.Next())
	check.Equal(t, 1).Assert(it1.Current())

	it2 := s.Backwards().Iterate()
	check.True(t).Assert(it2.Next())
	check.Equal(t, 5).Assert(it2.Current())

	s.Append(8, 9)
	check.String(t, `1, 2, 3, 4, 5, 8, 9`).Assert(s)

	check.True(t).Assert(it1.Next())
	check.Equal(t, 2).Assert(it1.Current())

	check.True(t).Assert(it2.Next())
	check.Equal(t, 4).Assert(it2.Current())

	s.Remove(1, 4)
	check.String(t, `1, 8, 9`).Assert(s)

	check.True(t).Assert(it1.Next())
	check.Equal(t, 9).Assert(it1.Current())
	check.False(t).Assert(it1.Next())

	check.True(t).Assert(it2.Next())
	check.Equal(t, 9).Assert(it2.Current())
	check.True(t).Assert(it2.Next())
	check.Equal(t, 8).Assert(it2.Current())
	check.True(t).Assert(it2.Next())
	check.Equal(t, 1).Assert(it2.Current())
	check.False(t).Assert(it2.Next())
}
