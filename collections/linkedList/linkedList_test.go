package linkedList

import (
	"testing"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/collections/iterator"
	"github.com/Snow-Gremlin/goToolbox/collections/predicate"
	"github.com/Snow-Gremlin/goToolbox/testers/check"
)

func validate[T any](t *testing.T, list collections.List[T]) {
	s, ok := list.(*linkedListImp[T])
	check.True(t).Assert(ok)

	if s.head == nil {
		check.Nil(t).Assert(s.tail)
		check.Zero(t).Assert(s.count)
		return
	}

	n := s.head
	var prior *node[T]
	count := 0
	for {
		check.Same(t, prior).Name(`Checking Previous Pointer`).
			With(`index`, count).Assert(n.prev)
		count++
		if n.next == nil {
			break
		}
		prior = n
		n = n.next
	}

	check.Same(t, n).Name(`Checking Last`).Assert(s.tail)
	check.Equal(t, count).Name(`checking Count`).Assert(s.count)
}

func Test_LinkedList(t *testing.T) {
	s := With(1, 2, 3, 4, 5)
	validate(t, s)
	check.String(t, `1, 2, 3, 4, 5`).Assert(s)
	check.String(t, `1, 2, 3, 4, 5`).Assert(s.Readonly())
	check.Length(t, 5).Assert(s)
	check.False(t).Assert(s.Empty())
	check.Equal(t, []int{1, 2, 3, 4, 5}).Assert(s.ToSlice())
	check.Equal(t, []int{1, 2, 3, 4, 5}).Assert(s.Enumerate().ToSlice())
	check.Equal(t, []int{5, 4, 3, 2, 1}).Assert(s.Backwards().ToSlice())

	sc := make([]int, 3)
	s.CopyToSlice(sc)
	check.Equal(t, []int{1, 2, 3}).Assert(sc)

	sc = make([]int, 7)
	s.CopyToSlice(sc)
	check.Equal(t, []int{1, 2, 3, 4, 5, 0, 0}).Assert(sc)

	check.Equal(t, 3).Assert(s.IndexOf(4))
	validate(t, s)
	check.Equal(t, -1).Assert(s.IndexOf(6))
	validate(t, s)
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
	check.String(t, `1, 2, 3, 4, 5`).Assert(s2)

	s2.Set(2, 11)
	validate(t, s2)
	check.String(t, `1, 2, 11, 4, 5`).Assert(s2)
	s2.Set(2, 24, 42)
	validate(t, s2)
	check.String(t, `1, 2, 24, 42, 5`).Assert(s2)
	s2.Set(2)
	validate(t, s2)
	check.String(t, `1, 2, 24, 42, 5`).Assert(s2)
	check.NotEqual(t, s).Assert(s2)
	check.MatchError(t, `^index out of bounds \{count: 5, index: -1\}$`).Panic(func() { s2.Set(-1, 45) })
	check.MatchError(t, `^index out of bounds \{count: 5, index: 6\}$`).Panic(func() { s2.Set(6, 45) })
	check.Length(t, 5).Assert(s2)

	s2.Remove(1, 3)
	validate(t, s2)
	check.String(t, `1, 5`).Assert(s2)
	check.NotEqual(t, s).Assert(s2)
	check.Length(t, 2).Assert(s2)

	s2.Clear()
	validate(t, s2)
	check.String(t, ``).Assert(s2)
	check.NotEqual(t, s).Assert(s2)
	check.MatchError(t, `^collection contains no values \{action: First\}$`).Panic(func() { s2.First() })
	check.MatchError(t, `^collection contains no values \{action: Last\}$`).Panic(func() { s2.Last() })
	check.MatchError(t, `^collection contains no values \{action: TakeFirst\}$`).Panic(func() { s2.TakeFirst() })
	check.MatchError(t, `^collection contains no values \{action: TakeLast\}$`).Panic(func() { s2.TakeLast() })

	s2.Prepend(1)
	validate(t, s2)
	s2.Append(2)
	validate(t, s2)
	s2.Prepend(3)
	validate(t, s2)
	check.String(t, `3, 1, 2`).Assert(s2)
	check.Equal(t, 3).Assert(s2.TakeFirst())
	validate(t, s2)
	check.String(t, `1, 2`).Assert(s2)
	check.Equal(t, 2).Assert(s2.TakeLast())
	validate(t, s2)
	check.String(t, `1`).Assert(s2)
	s2.Append(2, 3)
	validate(t, s2)
	check.String(t, `1, 2, 3`).Assert(s2)

	s2.Insert(1, 7, 2, 8, 9)
	validate(t, s2)
	check.String(t, `1, 7, 2, 8, 9, 2, 3`).Assert(s2)
	check.Equal(t, 2).Assert(s2.IndexOf(2, -10))
	check.Equal(t, 2).Assert(s2.IndexOf(2, 1))
	check.Equal(t, 5).Assert(s2.IndexOf(2, 2))
	check.Equal(t, -1).Assert(s2.IndexOf(2, 5))
	check.False(t).Assert(s2.RemoveIf(predicate.LessThan(0)))
	validate(t, s2)
	check.String(t, `1, 7, 2, 8, 9, 2, 3`).Assert(s2)
	check.True(t).Assert(s2.RemoveIf(predicate.LessThan(4)))
	validate(t, s2)
	check.String(t, `7, 8, 9`).Assert(s2)

	s3 := From(s.Enumerate().Strings())
	check.Equal(t, []string{`1`, `2`, `3`, `4`, `5`}).Assert(s3.ToSlice())
	validate(t, s2)
}

func Test_LinkedList_MoreModifiers(t *testing.T) {
	s := New[int]()
	validate(t, s)
	s.Append(1, 2, 3)
	validate(t, s)
	check.String(t, `1, 2, 3`).Assert(s)
	check.String(t, ``).Assert(s.TakeFront(0))
	validate(t, s)
	check.String(t, `1`).Assert(s.TakeFront(1))
	validate(t, s)
	check.String(t, `2, 3`).Assert(s.TakeFront(42))
	validate(t, s)
	check.String(t, ``).Assert(s)

	s.Prepend(1, 2, 3)
	validate(t, s)
	check.String(t, `1, 2, 3`).Assert(s)
	check.String(t, ``).Assert(s.TakeBack(0))
	validate(t, s)
	check.String(t, `3`).Assert(s.TakeBack(1))
	validate(t, s)
	check.String(t, `1, 2`).Assert(s.TakeBack(42))
	validate(t, s)
	check.String(t, ``).Assert(s)

	s.Insert(0)
	validate(t, s)
	s.Insert(0, 5)
	validate(t, s)
	check.Equal(t, 5).Assert(s.TakeFirst())
	validate(t, s)
	s.Insert(0, 7)
	validate(t, s)
	check.Equal(t, 7).Assert(s.TakeLast())
	validate(t, s)

	s.Insert(0, 4)
	validate(t, s)
	s.Insert(1, 5)
	validate(t, s)
	check.String(t, `4, 5`).Assert(s)
	check.MatchError(t, `^index out of bounds \{count: 2, index: -1\}$`).Panic(func() { s.Insert(-1, 3) })
	check.MatchError(t, `^index out of bounds \{count: 2, index: 3\}$`).Panic(func() { s.Insert(3, 3) })

	s.Remove(0, 0)
	check.MatchError(t, `^index out of bounds \{count: 2, index: -1\}$`).Panic(func() { s.Remove(-1, 1) })
	check.MatchError(t, `^index out of bounds \{count: 2, index: 3\}$`).Panic(func() { s.Remove(3, 1) })
	s.Insert(2, 6)
	validate(t, s)
	check.String(t, `4, 5, 6`).Assert(s)
	s.Remove(0, 3)
	check.String(t, ``).Assert(s)
	validate(t, s)
}

func Test_LinkedList_New(t *testing.T) {
	s := New[int]()
	check.Empty(t).Assert(s)
	check.String(t, ``).Assert(s)
	validate(t, s)

	s = New[int](5)
	check.Length(t, 5).Assert(s)
	check.String(t, `0, 0, 0, 0, 0`).Assert(s)
	validate(t, s)

	check.MatchError(t, `^invalid number of arguments \{count: 2, maximum: 1, usage: size\}$`).
		Panic(func() { New[int](1, 2) })

	s = New[int](-1)
	check.Empty(t).Assert(s)
	check.String(t, ``).Assert(s)
	validate(t, s)

	s = From[int](nil)
	check.Empty(t).Assert(s)
	check.String(t, ``).Assert(s)
	validate(t, s)

	s = Fill[int](8, 5)
	check.Length(t, 5).Assert(s)
	check.String(t, `8, 8, 8, 8, 8`).Assert(s)
	validate(t, s)

	s = Fill[int](8, -1)
	check.Empty(t).Assert(s)
	check.String(t, ``).Assert(s)
	validate(t, s)
}

func Test_LinkedList_MoreMethods(t *testing.T) {
	s1 := From(enumerator.Range(0, 3))
	check.Length(t, 3).Assert(s1)
	validate(t, s1)

	s1.PrependFrom(nil)
	s1.AppendFrom(nil)
	check.Length(t, 3).Assert(s1)
	check.String(t, `0, 1, 2`).Assert(s1)
	validate(t, s1)

	s1.PrependFrom(enumerator.Range(5, 3))
	check.Length(t, 6).Assert(s1)
	check.String(t, `5, 6, 7, 0, 1, 2`).Assert(s1)
	validate(t, s1)

	s2 := s1.TakeFront(2)
	check.Length(t, 2).Assert(s2)
	check.String(t, `5, 6`).Assert(s2)
	validate(t, s2)
	check.Length(t, 4).Assert(s1)
	check.String(t, `7, 0, 1, 2`).Assert(s1)
	validate(t, s1)

	s3 := s1.TakeBack(2)
	check.Length(t, 2).Assert(s3)
	check.String(t, `1, 2`).Assert(s3)
	validate(t, s3)
	check.Length(t, 2).Assert(s1)
	check.String(t, `7, 0`).Assert(s1)
	validate(t, s1)

	s4 := s1.TakeFront(12)
	check.Length(t, 2).Assert(s4)
	check.String(t, `7, 0`).Assert(s4)
	validate(t, s4)
	check.Empty(t).Assert(s1)
	check.String(t, ``).Assert(s1)
	validate(t, s1)

	s4 = s1.TakeFront(12)
	check.Empty(t).Assert(s4)
	validate(t, s4)

	s1.InsertFrom(0, enumerator.Enumerate(9, 8))
	check.Length(t, 2).Assert(s1)
	check.String(t, `9, 8`).Assert(s1)
	validate(t, s1)

	s5 := s1.TakeBack(12)
	check.Length(t, 2).Assert(s5)
	check.String(t, `9, 8`).Assert(s5)
	validate(t, s5)
	check.Empty(t).Assert(s1)
	check.String(t, ``).Assert(s1)
	validate(t, s1)

	s5 = s1.TakeBack(12)
	check.Empty(t).Assert(s5)
	validate(t, s5)

	s1.Append(9, 7, 5)
	check.String(t, `9, 7, 5`).Assert(s1)
	s1.SetFrom(1, nil)
	s1.SetFrom(1, enumerator.Enumerate(77, 55, 33, 11))
	check.Length(t, 5).Assert(s1)
	check.String(t, `9, 77, 55, 33, 11`).Assert(s1)
	validate(t, s1)

	s1.RemoveIf(predicate.GreaterThan(-1))
	check.Empty(t).Assert(s1)
	check.String(t, ``).Assert(s1)
}

func Test_LinkedList_StableIteration(t *testing.T) {
	s := With(1, 2, 3, 4, 5)

	it1 := s.Enumerate().Iterate()
	it2 := s.Backwards().Iterate()
	check.Equal(t, []int{1, 2}).Assert(iterator.ToSlice(iterator.Take(it1, 2)))
	check.Equal(t, []int{5, 4}).Assert(iterator.ToSlice(iterator.Take(it2, 2)))
	s.Insert(1, 11, 22)
	s.Insert(6, 44, 55)
	check.String(t, `1, 11, 22, 2, 3, 4, 44, 55, 5`).Assert(s)
	check.Equal(t, []int{3, 4, 44, 55, 5}).Assert(iterator.ToSlice(it1))
	check.Equal(t, []int{3, 2, 22, 11, 1}).Assert(iterator.ToSlice(it2))

	it1 = s.Enumerate().Iterate()
	it2 = s.Backwards().Iterate()
	check.Equal(t, []int{1, 11}).Assert(iterator.ToSlice(iterator.Take(it1, 2)))
	check.Equal(t, []int{5, 55}).Assert(iterator.ToSlice(iterator.Take(it2, 2)))
	s.Remove(2, 2)
	s.Remove(3, 2)
	check.String(t, `1, 11, 3, 55, 5`).Assert(s)

	check.MatchError(t, `Collection was modified; iteration may not continue`).
		Panic(func() { it1.Next() })
	check.MatchError(t, `Collection was modified; iteration may not continue`).
		Panic(func() { it2.Next() })
}
