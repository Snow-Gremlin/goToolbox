package capStack

import (
	"bytes"
	"testing"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/collections/list"
	"github.com/Snow-Gremlin/goToolbox/collections/predicate"
	"github.com/Snow-Gremlin/goToolbox/events/listener"
	"github.com/Snow-Gremlin/goToolbox/testers/check"
)

type nodeCollection map[any]bool

func validate[T any](t *testing.T, stack collections.Stack[T], nc nodeCollection) {
	s := stack.(*capStackImp[T])
	for n := range nc {
		nc[n] = false
	}

	n := s.head
	touched := map[*node[T]]int{}
	touched[n] = 0
	nc[n] = true
	count := 0
	for n != nil {
		n = n.prev
		count++

		old, exists := touched[n]
		check.False(t).
			With(`old index`, old).
			Name(`loop check in stack`).
			With(`index`, count).
			Require(exists)

		touched[n] = count
		nc[n] = true
	}

	check.Equal(t, count).Name(`Count`).Assert(s.count)

	g := s.graveyard
	count = 0
	for g != nil {
		old, exists := touched[g]
		c := check.False(t).
			Name(`stack graveyard check`)
		if old >= 0 {
			c = c.With(`old index in stack`, old)
		} else {
			c = c.With(`old index in graveyard`, -old-1)
		}
		c.Require(exists)

		touched[g] = -1 - count
		nc[g] = true
		check.Zero(t).Assert(g.value)

		g = g.prev
		count++
	}

	nodesLost := false
	for _, found := range nc {
		if !found {
			nodesLost = true
			break
		}
	}
	check.False(t).
		Name(`check for lost nodes`).
		With(`all nodes`, nc).
		Assert(nodesLost)
}

func checkTombCount[T any](t *testing.T, stack collections.Stack[T], expCount int) {
	check.Equal(t, expCount).Assert(stack.(*capStackImp[T]).tombs())
}

func Test_CapStack(t *testing.T) {
	s := With(1, 2, 3)
	snc := nodeCollection{}
	validate(t, s, snc)
	check.False(t).Assert(s.Empty())
	check.Length(t, 3).Assert(s)
	check.Equal(t, []int{1, 2, 3}).Assert(s.Enumerate().ToSlice())
	check.Equal(t, []int{1, 2, 3}).Assert(s.Readonly().ToSlice())
	check.String(t, `1, 2, 3`).Assert(s)
	check.Equal(t, list.With(1, 2, 3)).Assert(s.ToList())

	s.Push() // no effect
	validate(t, s, snc)
	s.Push(4, 5)
	validate(t, s, snc)
	check.Length(t, 5).Assert(s)
	check.String(t, `4, 5, 1, 2, 3`).Assert(s)
	check.Equal(t, []int{4, 5, 1, 2, 3}).Assert(s.ToSlice())
	checkTombCount(t, s, 5)

	sc := make([]int, 4)
	s.CopyToSlice(sc)
	check.Equal(t, []int{4, 5, 1, 2}).Assert(sc)

	check.Equal(t, 4).Assert(s.Peek())
	v, ok := s.TryPeek()
	check.Equal(t, 4).Assert(v)
	check.True(t).Assert(ok)

	s2 := s.Clone()
	s2nc := nodeCollection{}
	validate(t, s2, s2nc)
	check.Equal(t, s).Assert(s2)
	check.String(t, `4, 5, 1, 2, 3`).Assert(s2)

	check.Equal(t, 4).Assert(s.Pop())
	validate(t, s, snc)
	v, ok = s.TryPop()
	check.Equal(t, 5).Assert(v)
	check.True(t).Assert(ok)
	validate(t, s, snc)
	check.String(t, `1, 2, 3`).Assert(s)
	checkTombCount(t, s, 7)

	check.Equal(t, []int{}).Assert(s.Take(0))
	check.Equal(t, []int{1}).Assert(s.Take(1))
	validate(t, s, snc)
	check.Equal(t, []int{2, 3}).Assert(s.Take(3))
	validate(t, s, snc)
	check.String(t, ``).Assert(s)
	check.True(t).Assert(s.Empty())
	check.NotEqual(t, s).Assert(s2)
	checkTombCount(t, s, 10)

	check.MatchError(t, `^collection contains no values \{action: Pop\}$`).
		Panic(func() { s.Pop() })
	check.MatchError(t, `^collection contains no values \{action: Peek\}$`).
		Panic(func() { s.Peek() })
	v, ok = s.TryPop()
	check.Zero(t).Assert(v)
	check.False(t).Assert(ok)
	v, ok = s.TryPeek()
	check.Zero(t).Assert(v)
	check.False(t).Assert(ok)
	checkTombCount(t, s, 10)

	checkTombCount(t, s2, 5)
	s2.Clear()
	check.String(t, ``).Assert(s2)
	validate(t, s2, s2nc)
	checkTombCount(t, s2, 10)

	s2.(*capStackImp[int]).growCap(4) // no effect
	checkTombCount(t, s2, 10)

	s2.Push(42)
	check.String(t, `42`).Assert(s2)
	validate(t, s2, s2nc)
	checkTombCount(t, s2, 9)
	check.Equal(t, 42).Assert(s2.Pop())
	check.String(t, ``).Assert(s2)
	validate(t, s2, s2nc)
	checkTombCount(t, s2, 10)

	s2.PushFrom(enumerator.Range(3, 5))
	check.String(t, `3, 4, 5, 6, 7`).Assert(s2)
	validate(t, s2, s2nc)
	checkTombCount(t, s2, 5)

	s2.PushFrom(enumerator.Range(14, 3))
	check.String(t, `14, 15, 16, 3, 4, 5, 6, 7`).Assert(s2)
	validate(t, s2, s2nc)
	checkTombCount(t, s2, 2)

	s2.PushFrom(nil) // No effect
	check.String(t, `14, 15, 16, 3, 4, 5, 6, 7`).Assert(s2)
	validate(t, s2, s2nc)
	checkTombCount(t, s2, 2)

	s2.TrimTo(10) // no effect
	check.String(t, `14, 15, 16, 3, 4, 5, 6, 7`).Assert(s2)
	validate(t, s2, s2nc)
	checkTombCount(t, s2, 2)

	s2.TrimTo(6)
	check.String(t, `14, 15, 16, 3, 4, 5`).Assert(s2)
	validate(t, s2, s2nc)
	checkTombCount(t, s2, 4)

	s2.TrimTo(0)
	check.String(t, ``).Assert(s2)
	validate(t, s2, s2nc)
	checkTombCount(t, s2, 10)

	s2.Clear()
	check.String(t, ``).Assert(s2)
	validate(t, s2, s2nc)
	checkTombCount(t, s2, 10)

	s2.Clip()
	s2nc = nodeCollection{} // reset for because of Clip
	check.String(t, ``).Assert(s2)
	validate(t, s2, s2nc)
	checkTombCount(t, s2, 0)

	s2.Push(3, 4, 5)
	check.String(t, `3, 4, 5`).Assert(s2)
	validate(t, s2, s2nc)
	checkTombCount(t, s2, 7)

	s2.Clip()
	s2nc = nodeCollection{} // reset for because of Clip
	check.String(t, `3, 4, 5`).Assert(s2)
	validate(t, s2, s2nc)
	checkTombCount(t, s2, 0)

	s2.(*capStackImp[int]).growCap(2) // no effect
	check.String(t, `3, 4, 5`).Assert(s2)
	validate(t, s2, s2nc)
	checkTombCount(t, s2, 0)

	s2.Clip() // no effect
	check.String(t, `3, 4, 5`).Assert(s2)
	validate(t, s2, s2nc)
	checkTombCount(t, s2, 0)
}

func Test_CapStack_New(t *testing.T) {
	s := New[int]()
	check.Empty(t).Assert(s)
	check.String(t, ``).Assert(s)
	checkTombCount(t, s, 0)

	s = New[int](5)
	check.Length(t, 5).Assert(s)
	check.String(t, `0, 0, 0, 0, 0`).Assert(s)
	checkTombCount(t, s, 0)

	s = New[int](5, 8)
	check.Length(t, 5).Assert(s)
	check.String(t, `0, 0, 0, 0, 0`).Assert(s)
	checkTombCount(t, s, 3)

	check.MatchError(t, `^invalid number of arguments \{count: 3, maximum: 2, usage: size and capacity\}$`).
		Panic(func() { New[int](1, 2, 3) })

	s = New[int](-1)
	check.Empty(t).Assert(s)
	check.String(t, ``).Assert(s)
	checkTombCount(t, s, 0)

	s = Fill[int](8, 5)
	check.Length(t, 5).Assert(s)
	check.String(t, `8, 8, 8, 8, 8`).Assert(s)
	checkTombCount(t, s, 0)

	s = Fill[int](8, -1)
	check.Empty(t).Assert(s)
	check.String(t, ``).Assert(s)
	checkTombCount(t, s, 0)

	s = Fill[int](8, 5, 8)
	check.Length(t, 5).Assert(s)
	check.String(t, `8, 8, 8, 8, 8`).Assert(s)
	checkTombCount(t, s, 3)

	check.MatchError(t, `^invalid number of arguments \{count: 2, maximum: 1, usage: capacity\}$`).
		Panic(func() { Fill[int](8, 5, 8, 13) })

	s = From[int](enumerator.Range(1, 3))
	check.Length(t, 3).Assert(s)
	check.String(t, `1, 2, 3`).Assert(s)
	checkTombCount(t, s, 7)

	s = From[int](nil, 3)
	check.Empty(t).Assert(s)
	check.String(t, ``).Assert(s)
	checkTombCount(t, s, 3)

	check.MatchError(t, `^invalid number of arguments \{count: 2, maximum: 1, usage: capacity\}$`).
		Panic(func() { From[int](nil, 7, 8) })
}

func Test_CapStack_UnstableIteration(t *testing.T) {
	s := With(0, 1, 2, 3, 4, 5)
	e := enumerator.Select(s.Enumerate(), func(value int) int {
		return 10 - value
	}).WhereNot(predicate.Eq(7))
	check.Equal(t, []int{10, 9, 8, 6, 5}).Assert(e.ToSlice())

	it := e.Iterate()
	check.True(t).Assert(it.Next())
	check.Equal(t, 10).Assert(it.Current())

	// Iteration still works because pushing doesn't modify existing nodes.
	s.Push(6)
	check.Equal(t, []int{4, 10, 9, 8, 6, 5}).Assert(e.ToSlice())
	check.True(t).Assert(it.Next())
	check.Equal(t, 9).Assert(it.Current())

	// Enumerations using the same closures still function with changes.
	// Iteration will break since popping may break the nodes that
	// some iteration may or may not be at.
	check.Equal(t, 6).Assert(s.Pop())
	check.Equal(t, []int{10, 9, 8, 6, 5}).Assert(e.ToSlice())
	check.MatchError(t, `Collection was modified; iteration may not continue`).
		Panic(func() { it.Next() })

	// Will panic because the stack is modified inside of foreach.
	// Stack will still have change for foreach because the iteration
	// will fail, not the modification.
	check.MatchError(t, `Collection was modified; iteration may not continue`).
		Panic(func() {
			e.Foreach(func(value int) {
				if value == 6 {
					s.Pop()
				}
			})
		})
	check.Equal(t, []int{9, 8, 6, 5}).Assert(e.ToSlice())

	it = e.Iterate()
	check.True(t).Assert(it.Next())
	check.Equal(t, 9).Assert(it.Current())

	s.PushFrom(enumerator.Enumerate(7, 8)) // continue iterators
	check.Equal(t, []int{3, 2, 9, 8, 6, 5}).Assert(e.ToSlice())
	check.True(t).Assert(it.Next())
	check.Equal(t, 8).Assert(it.Current())

	check.Equal(t, []int{7, 8}).Assert(s.Take(2)) // break iterators
	check.Equal(t, []int{9, 8, 6, 5}).Assert(e.ToSlice())
	check.MatchError(t, `Collection was modified; iteration may not continue`).
		Panic(func() { it.Next() })

	it = e.Iterate()
	check.True(t).Assert(it.Next())
	check.Equal(t, 9).Assert(it.Current())

	s.TrimTo(2) // break iterators
	check.Equal(t, []int{9, 8}).Assert(e.ToSlice())
	check.MatchError(t, `Collection was modified; iteration may not continue`).
		Panic(func() { it.Next() })

	it = e.Iterate()
	check.True(t).Assert(it.Next())
	check.Equal(t, 9).Assert(it.Current())

	s.Clear() // break iterators
	check.Empty(t).Assert(e.ToSlice())
	check.MatchError(t, `Collection was modified; iteration may not continue`).
		Panic(func() { it.Next() })
}

func Test_CapStack_OnChange(t *testing.T) {
	buf := &bytes.Buffer{}
	s := New[int]()
	lis := listener.New(func(args collections.ChangeArgs) {
		_, _ = buf.WriteString(args.Type().String())
	})
	defer lis.Cancel()
	check.True(t).Assert(lis.Subscribe(s.OnChange()))
	check.StringAndReset(t, ``).Assert(buf)

	s.Push()
	check.StringAndReset(t, ``).Assert(buf)
	s.Push(1, 2, 3)
	check.StringAndReset(t, `Added`).Assert(buf)
	check.Equal(t, 1).Assert(s.Pop())
	check.StringAndReset(t, `Removed`).Assert(buf)
	check.String(t, `2, 3`).Assert(s)

	s.PushFrom(nil)
	check.StringAndReset(t, ``).Assert(buf)
	s.PushFrom(enumerator.Enumerate[int]())
	check.StringAndReset(t, ``).Assert(buf)
	s.PushFrom(enumerator.Enumerate[int](4, 5, 6))
	check.StringAndReset(t, `Added`).Assert(buf)
	check.String(t, `4, 5, 6, 2, 3`).Assert(s)

	check.String(t, `[]`).Assert(s.Take(0))
	check.StringAndReset(t, ``).Assert(buf)
	check.String(t, `[4 5 6]`).Assert(s.Take(3))
	check.StringAndReset(t, `Removed`).Assert(buf)
	check.String(t, `2, 3`).Assert(s)

	s.Clear()
	check.StringAndReset(t, `Removed`).Assert(buf)
	s.Clear()
	check.StringAndReset(t, ``).Assert(buf)
	_, dequeue := s.TryPop()
	check.False(t).Assert(dequeue)
	check.StringAndReset(t, ``).Assert(buf)

	s.Push(1, 2, 3)
	check.StringAndReset(t, `Added`).Assert(buf)
	s.TrimTo(5)
	check.StringAndReset(t, ``).Assert(buf)
	s.TrimTo(3)
	check.StringAndReset(t, ``).Assert(buf)
	s.TrimTo(2)
	check.StringAndReset(t, `Removed`).Assert(buf)
	check.String(t, `1, 2`).Assert(s)
	s.TrimTo(-1)
	check.StringAndReset(t, `Removed`).Assert(buf)
	check.String(t, ``).Assert(s)
}
