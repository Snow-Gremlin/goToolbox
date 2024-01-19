package capQueue

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

func validate[T any](t *testing.T, queue collections.Queue[T], nc nodeCollection) {
	q := queue.(*capQueueImp[T])
	for n := range nc {
		nc[n] = false
	}

	touched := map[*node[T]]int{}
	if q.head == nil {
		check.Nil(t).Assert(q.tail)
		check.Zero(t).Assert(q.count)

	} else {
		n := q.head
		touched[n] = 0
		nc[n] = true
		count := 1
		for n.next != nil {
			n = n.next
			count++

			old, exists := touched[n]
			check.False(t).
				With(`old index`, old).
				Name(`loop check in queue`).
				With(`index`, count).
				Require(exists)

			touched[n] = count
			nc[n] = true
		}

		check.Equal(t, count).Name(`Count`).Assert(q.count)
		check.Same(t, n).Name(`Tail`).Assert(q.tail)
	}

	g := q.graveyard
	count := 0
	for g != nil {
		old, exists := touched[g]
		c := check.False(t).
			Name(`queue graveyard check`)
		if old >= 0 {
			c = c.With(`old index in queue`, old)
		} else {
			c = c.With(`old index in graveyard`, -old-1)
		}
		c.Require(exists)

		touched[g] = -1 - count
		nc[g] = true
		check.Zero(t).Assert(g.value)

		g = g.next
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

func checkTombCount[T any](t *testing.T, queue collections.Queue[T], expCount int) {
	check.Equal(t, expCount).Assert(queue.(*capQueueImp[T]).tombs())
}

func Test_CapQueue(t *testing.T) {
	q := With(1, 2, 3)
	qnc := nodeCollection{}
	validate(t, q, qnc)
	check.False(t).Assert(q.Empty())
	check.Length(t, 3).Assert(q)
	check.Equal(t, []int{1, 2, 3}).Assert(q.Enumerate().ToSlice())
	check.Equal(t, list.With(1, 2, 3)).Assert(q.ToList())
	check.Equal(t, []int{1, 2, 3}).Assert(q.Readonly().ToSlice())
	check.String(t, `1, 2, 3`).Assert(q)

	s := make([]int, 2)
	q.CopyToSlice(s)
	check.Equal(t, []int{1, 2}).Assert(s)

	s = make([]int, 5)
	q.CopyToSlice(s)
	check.Equal(t, []int{1, 2, 3, 0, 0}).Assert(s)

	q.Enqueue() // no effect
	validate(t, q, qnc)
	q.Enqueue(4, 5)
	validate(t, q, qnc)
	check.Length(t, 5).Assert(q)
	check.String(t, `1, 2, 3, 4, 5`).Assert(q)
	check.Equal(t, []int{1, 2, 3, 4, 5}).Assert(q.ToSlice())
	checkTombCount(t, q, 5)

	check.Equal(t, 1).Assert(q.Peek())
	v, ok := q.TryPeek()
	check.Equal(t, 1).Assert(v)
	check.True(t).Assert(ok)

	q2 := q.Clone()
	q2nc := nodeCollection{}
	check.Equal(t, q).Assert(q2)
	check.String(t, `1, 2, 3, 4, 5`).Assert(q2)
	validate(t, q2, q2nc)

	check.Equal(t, 1).Assert(q.Dequeue())
	validate(t, q, qnc)
	v, ok = q.TryDequeue()
	check.Equal(t, 2).Assert(v)
	check.True(t).Assert(ok)
	validate(t, q, qnc)
	check.String(t, `3, 4, 5`).Assert(q)
	checkTombCount(t, q, 7)

	check.Equal(t, []int{}).Assert(q.Take(0))
	check.Equal(t, []int{3}).Assert(q.Take(1))
	validate(t, q, qnc)
	check.Equal(t, []int{4, 5}).Assert(q.Take(3))
	validate(t, q, qnc)
	check.String(t, ``).Assert(q)
	check.True(t).Assert(q.Empty())
	check.NotEqual(t, q).Assert(q2)
	checkTombCount(t, q, 10)

	check.MatchError(t, `^collection contains no values \{action: Dequeue\}$`).
		Panic(func() { q.Dequeue() })
	check.MatchError(t, `^collection contains no values \{action: Peek\}$`).
		Panic(func() { q.Peek() })
	v, ok = q.TryDequeue()
	check.Zero(t).Assert(v)
	check.False(t).Assert(ok)
	v, ok = q.TryPeek()
	check.Zero(t).Assert(v)
	check.False(t).Assert(ok)
	checkTombCount(t, q, 10)

	checkTombCount(t, q2, 5)
	q2.Clear()
	check.String(t, ``).Assert(q2)
	validate(t, q2, q2nc)
	checkTombCount(t, q2, 10)

	q2.(*capQueueImp[int]).growCap(4) // no effect
	checkTombCount(t, q2, 10)

	q2.Clear() // no effect
	check.String(t, ``).Assert(q2)
	validate(t, q2, q2nc)
	q2.Enqueue(42)
	check.String(t, `42`).Assert(q2)
	checkTombCount(t, q2, 9)
	validate(t, q2, q2nc)
	check.Equal(t, 42).Assert(q2.Dequeue())
	check.String(t, ``).Assert(q2)
	validate(t, q2, q2nc)
	checkTombCount(t, q2, 10)

	q2.EnqueueFrom(enumerator.Range(4, 3))
	check.String(t, `4, 5, 6`).Assert(q2)
	validate(t, q2, q2nc)
	checkTombCount(t, q2, 7)

	q2.EnqueueFrom(enumerator.Range(9, 3))
	check.String(t, `4, 5, 6, 9, 10, 11`).Assert(q2)
	validate(t, q2, q2nc)
	checkTombCount(t, q2, 4)

	q2.(*capQueueImp[int]).growCap(10)
	checkTombCount(t, q2, 4)
	q2.Clip()
	q2nc = nodeCollection{} // reset for because of Clip
	checkTombCount(t, q2, 0)
	q2.Clip() // no effect
	check.String(t, `4, 5, 6, 9, 10, 11`).Assert(q2)
	validate(t, q2, q2nc)
	checkTombCount(t, q2, 0)
}

func Test_CapQueue_New(t *testing.T) {
	s := New[int]()
	check.Empty(t).Assert(s)
	check.String(t, ``).Assert(s)
	checkTombCount(t, s, 0)

	s = New[int](5)
	check.Length(t, 5).Assert(s)
	check.String(t, `0, 0, 0, 0, 0`).Assert(s)
	checkTombCount(t, s, 0)

	s = New[int](5, 9)
	check.Length(t, 5).Assert(s)
	check.String(t, `0, 0, 0, 0, 0`).Assert(s)
	checkTombCount(t, s, 4)

	s = New[int](5, 4)
	check.Length(t, 5).Assert(s)
	check.String(t, `0, 0, 0, 0, 0`).Assert(s)
	checkTombCount(t, s, 0)

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

	s = Fill[int](8, 5, 9)
	check.Length(t, 5).Assert(s)
	check.String(t, `8, 8, 8, 8, 8`).Assert(s)
	checkTombCount(t, s, 4)

	s = Fill[int](8, 5, -1)
	check.Length(t, 5).Assert(s)
	check.String(t, `8, 8, 8, 8, 8`).Assert(s)
	checkTombCount(t, s, 0)

	check.MatchError(t, `^invalid number of arguments \{count: 2, maximum: 1, usage: capacity\}$`).
		Panic(func() { Fill(1, 2, 3, 4) })

	s = From[int](nil)
	check.Empty(t).Assert(s)
	check.String(t, ``).Assert(s)
	checkTombCount(t, s, 0)

	s = From[int](nil, 5)
	check.Empty(t).Assert(s)
	check.String(t, ``).Assert(s)
	checkTombCount(t, s, 5)

	s = From[int](nil, -1)
	check.Empty(t).Assert(s)
	check.String(t, ``).Assert(s)
	checkTombCount(t, s, 0)

	check.MatchError(t, `^invalid number of arguments \{count: 2, maximum: 1, usage: capacity\}$`).
		Panic(func() { From[int](nil, 1, 2) })
}

func Test_CapQueue_UnstableIteration(t *testing.T) {
	q := With(0, 1, 2, 3, 4, 5)
	e := enumerator.Select(q.Enumerate(), func(value int) int {
		return 10 - value
	}).WhereNot(predicate.Eq(7))
	check.Equal(t, []int{10, 9, 8, 6, 5}).Assert(e.ToSlice())

	it := e.Iterate()
	check.True(t).Assert(it.Next())
	check.Equal(t, 10).Assert(it.Current())

	// Iteration still works because Enqueue doesn't modify existing nodes.
	q.Enqueue(6)
	check.Equal(t, []int{10, 9, 8, 6, 5, 4}).Assert(e.ToSlice())
	check.True(t).Assert(it.Next())
	check.Equal(t, 9).Assert(it.Current())

	// Enumerations using the same closures still function with changes.
	// Iteration will break since dequeue-ing may break the nodes that
	// some iteration may or may not be at.
	check.Equal(t, 0).Assert(q.Dequeue())
	check.Equal(t, []int{9, 8, 6, 5, 4}).Assert(e.ToSlice())
	check.MatchError(t, `Collection was modified; iteration may not continue`).
		Panic(func() { it.Next() })

	// Will panic because the queue is modified inside of foreach.
	// Queue will still have change for foreach because the iteration
	// will fail, not the modification.
	check.MatchError(t, `Collection was modified; iteration may not continue`).
		Panic(func() {
			e.Foreach(func(value int) {
				if value == 6 {
					q.Dequeue()
				}
			})
		})
	check.Equal(t, []int{8, 6, 5, 4}).Assert(e.ToSlice())

	it = e.Iterate()
	check.True(t).Assert(it.Next())
	check.Equal(t, 8).Assert(it.Current())

	q.EnqueueFrom(enumerator.Enumerate(7, 8)) // continue iterators
	check.Equal(t, []int{8, 6, 5, 4, 3, 2}).Assert(e.ToSlice())
	check.True(t).Assert(it.Next())
	check.Equal(t, 6).Assert(it.Current())

	check.Equal(t, []int{2, 3}).Assert(q.Take(2)) // break iterators
	check.Equal(t, []int{6, 5, 4, 3, 2}).Assert(e.ToSlice())
	check.MatchError(t, `Collection was modified; iteration may not continue`).
		Panic(func() { it.Next() })

	it = e.Iterate()
	check.True(t).Assert(it.Next())
	check.Equal(t, 6).Assert(it.Current())

	q.Clear() // break iterators
	check.Empty(t).Assert(e.ToSlice())
	check.MatchError(t, `Collection was modified; iteration may not continue`).
		Panic(func() { it.Next() })
}

func Test_CapQueue_OnChange(t *testing.T) {
	buf := &bytes.Buffer{}
	q := New[int]()
	lis := listener.New(func(args collections.ChangeArgs) {
		_, _ = buf.WriteString(args.Type().String())
	})
	defer lis.Cancel()
	check.True(t).Assert(lis.Subscribe(q.OnChange()))
	check.StringAndReset(t, ``).Assert(buf)

	q.Enqueue()
	check.StringAndReset(t, ``).Assert(buf)
	q.Enqueue(1, 2, 3)
	check.StringAndReset(t, `Added`).Assert(buf)
	check.Equal(t, 1).Assert(q.Dequeue())
	check.StringAndReset(t, `Removed`).Assert(buf)
	check.String(t, `2, 3`).Assert(q)

	q.EnqueueFrom(nil)
	check.StringAndReset(t, ``).Assert(buf)
	q.EnqueueFrom(enumerator.Enumerate[int]())
	check.StringAndReset(t, ``).Assert(buf)
	q.EnqueueFrom(enumerator.Enumerate[int](4, 5, 6))
	check.StringAndReset(t, `Added`).Assert(buf)
	check.String(t, `2, 3, 4, 5, 6`).Assert(q)

	check.String(t, `[]`).Assert(q.Take(0))
	check.StringAndReset(t, ``).Assert(buf)
	check.String(t, `[2 3 4]`).Assert(q.Take(3))
	check.StringAndReset(t, `Removed`).Assert(buf)
	check.String(t, `5, 6`).Assert(q)

	q.Clear()
	check.StringAndReset(t, `Removed`).Assert(buf)
	q.Clear()
	check.StringAndReset(t, ``).Assert(buf)
	_, dequeue := q.TryDequeue()
	check.False(t).Assert(dequeue)
	check.StringAndReset(t, ``).Assert(buf)
}
