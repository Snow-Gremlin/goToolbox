package stack

import (
	"testing"

	"goToolbox/collections"
	"goToolbox/collections/enumerator"
	"goToolbox/collections/list"
	"goToolbox/testers/check"
)

func validate[T any](t *testing.T, stack collections.Stack[T]) {
	s := stack.(*stackImp[T])

	n := s.head
	touched := map[*node[T]]bool{}
	count := 0
	for n != nil {
		n = n.prev
		_, exists := touched[n]
		check.False(t).
			Name(`loop check in stack`).
			With(`count`, count).
			Require(exists)
		touched[n] = true
		count++
	}

	check.Equal(t, count).Name(`Count`).Assert(s.count)
}

func Test_Stack(t *testing.T) {
	s := With(1, 2, 3)
	validate(t, s)
	check.False(t).Assert(s.Empty())
	check.Length(t, 3).Assert(s)
	check.Equal(t, []int{1, 2, 3}).Assert(s.Enumerate().ToSlice())
	check.Equal(t, []int{1, 2, 3}).Assert(s.Readonly().ToSlice())
	check.String(t, `1, 2, 3`).Assert(s)
	check.Equal(t, list.With(1, 2, 3)).Assert(s.ToList())

	s.Push() // no effect
	s.Push(4, 5)
	validate(t, s)
	check.Length(t, 5).Assert(s)
	check.String(t, `4, 5, 1, 2, 3`).Assert(s)
	check.Equal(t, []int{4, 5, 1, 2, 3}).Assert(s.ToSlice())

	sc := make([]int, 4)
	s.CopyToSlice(sc)
	check.Equal(t, []int{4, 5, 1, 2}).Assert(sc)

	check.Equal(t, 4).Assert(s.Peek())
	v, ok := s.TryPeek()
	check.Equal(t, 4).Assert(v)
	check.True(t).Assert(ok)

	s2 := s.Clone()
	validate(t, s2)
	check.Equal(t, s).Assert(s2)
	check.String(t, `4, 5, 1, 2, 3`).Assert(s2)

	check.Equal(t, 4).Assert(s.Pop())
	validate(t, s)
	v, ok = s.TryPop()
	check.Equal(t, 5).Assert(v)
	check.True(t).Assert(ok)
	validate(t, s)
	check.String(t, `1, 2, 3`).Assert(s)

	check.Equal(t, []int{}).Assert(s.Take(0))
	check.Equal(t, []int{1}).Assert(s.Take(1))
	validate(t, s)
	check.Equal(t, []int{2, 3}).Assert(s.Take(3))
	validate(t, s)
	check.String(t, ``).Assert(s)
	check.True(t).Assert(s.Empty())
	check.NotEqual(t, s).Assert(s2)

	check.MatchError(t, `^collection contains no values \{action: Pop\}$`).Panic(func() { s.Pop() })
	check.MatchError(t, `^collection contains no values \{action: Peek\}$`).Panic(func() { s.Peek() })
	v, ok = s.TryPop()
	check.Zero(t).Assert(v)
	check.False(t).Assert(ok)
	v, ok = s.TryPeek()
	check.Zero(t).Assert(v)
	check.False(t).Assert(ok)

	s2.Clear()
	check.String(t, ``).Assert(s2)
	validate(t, s2)
	s2.Push(42)
	check.String(t, `42`).Assert(s2)
	validate(t, s2)
	check.Equal(t, 42).Assert(s2.Pop())
	check.String(t, ``).Assert(s2)
	validate(t, s2)

	s2.PushFrom(enumerator.Range(3, 5))
	check.String(t, `3, 4, 5, 6, 7`).Assert(s2)
	validate(t, s2)

	s2.PushFrom(enumerator.Range(14, 3))
	check.String(t, `14, 15, 16, 3, 4, 5, 6, 7`).Assert(s2)
	validate(t, s2)

	s2.PushFrom(nil) // No effect
	check.String(t, `14, 15, 16, 3, 4, 5, 6, 7`).Assert(s2)
	validate(t, s2)

	s2.TrimTo(10)
	check.String(t, `14, 15, 16, 3, 4, 5, 6, 7`).Assert(s2)
	validate(t, s2)

	s2.TrimTo(6)
	check.String(t, `14, 15, 16, 3, 4, 5`).Assert(s2)
	validate(t, s2)

	s2.Clip() // no effect
	check.String(t, `14, 15, 16, 3, 4, 5`).Assert(s2)
	validate(t, s2)

	s2.TrimTo(0)
	check.String(t, ``).Assert(s2)
	validate(t, s2)
}

func Test_Stack_New(t *testing.T) {
	s := New[int]()
	check.Empty(t).Assert(s)
	check.String(t, ``).Assert(s)

	s = New[int](5)
	check.Length(t, 5).Assert(s)
	check.String(t, `0, 0, 0, 0, 0`).Assert(s)

	check.MatchError(t, `^invalid number of arguments \{count: 2, maximum: 1, usage: size\}$`).
		Panic(func() { New[int](1, 2) })

	s = New[int](-1)
	check.Empty(t).Assert(s)
	check.String(t, ``).Assert(s)

	s = Fill[int](8, 5)
	check.Length(t, 5).Assert(s)
	check.String(t, `8, 8, 8, 8, 8`).Assert(s)

	s = Fill[int](8, -1)
	check.Empty(t).Assert(s)
	check.String(t, ``).Assert(s)
}

func Test_Stack_UnstableIteration(t *testing.T) {
	s := With(3, 4, 5, 6, 7)
	it := s.Enumerate().Iterate()

	check.True(t).Assert(it.Next())
	check.Equal(t, 3).Assert(it.Current())

	s.Push(8)
	check.String(t, `8, 3, 4, 5, 6, 7`).Assert(s)

	check.True(t).Assert(it.Next())
	check.Equal(t, 4).Assert(it.Current())

	check.Equal(t, 8).Assert(s.Pop())
	check.MatchError(t, `Collection was modified; iteration may not continue`).
		Panic(func() { it.Next() })
}
