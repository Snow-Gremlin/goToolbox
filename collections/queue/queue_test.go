package queue

import (
	"testing"

	"goToolbox/collections"
	"goToolbox/collections/enumerator"
	"goToolbox/collections/list"
	"goToolbox/testers/check"
)

func validate[T any](t *testing.T, queue collections.Queue[T]) {
	q := queue.(*queueImp[T])
	if q.head == nil {
		check.Nil(t).Assert(q.tail)
		check.Zero(t).Assert(q.count)
		return
	}

	n := q.head
	touched := map[*node[T]]bool{}
	touched[n] = true
	count := 1
	for n.next != nil {
		n = n.next
		_, exists := touched[n]
		check.False(t).
			Name(`loop check in queue`).
			With(`count`, count).
			Require(exists)
		touched[n] = true
		count++
	}

	check.Equal(t, count).Name(`Count`).Assert(q.count)
	check.Same(t, n).Name(`Tail`).Assert(q.tail)
}

func Test_Queue(t *testing.T) {
	q := With(1, 2, 3)
	validate(t, q)
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
	q.Enqueue(4, 5)
	validate(t, q)
	check.Length(t, 5).Assert(q)
	check.String(t, `1, 2, 3, 4, 5`).Assert(q)
	check.Equal(t, []int{1, 2, 3, 4, 5}).Assert(q.ToSlice())

	check.Equal(t, 1).Assert(q.Peek())
	v, ok := q.TryPeek()
	check.Equal(t, 1).Assert(v)
	check.True(t).Assert(ok)

	q2 := q.Clone()
	check.Equal(t, q).Assert(q2)
	check.String(t, `1, 2, 3, 4, 5`).Assert(q2)
	validate(t, q2)

	check.Equal(t, 1).Assert(q.Dequeue())
	validate(t, q)
	v, ok = q.TryDequeue()
	check.Equal(t, 2).Assert(v)
	check.True(t).Assert(ok)
	validate(t, q)
	check.String(t, `3, 4, 5`).Assert(q)

	check.Equal(t, []int{}).Assert(q.Take(0))
	check.Equal(t, []int{3}).Assert(q.Take(1))
	validate(t, q)
	check.Equal(t, []int{4, 5}).Assert(q.Take(3))
	validate(t, q)
	check.String(t, ``).Assert(q)
	check.True(t).Assert(q.Empty())
	check.NotEqual(t, q).Assert(q2)

	check.MatchError(t, `^collection contains no values \{action: Dequeue\}$`).Panic(func() { q.Dequeue() })
	check.MatchError(t, `^collection contains no values \{action: Peek\}$`).Panic(func() { q.Peek() })
	v, ok = q.TryDequeue()
	check.Zero(t).Assert(v)
	check.False(t).Assert(ok)
	v, ok = q.TryPeek()
	check.Zero(t).Assert(v)
	check.False(t).Assert(ok)

	q2.Clear()
	check.String(t, ``).Assert(q2)
	validate(t, q2)
	q2.Clear() // no effect
	check.String(t, ``).Assert(q2)
	validate(t, q2)
	q2.Enqueue(42)
	check.String(t, `42`).Assert(q2)
	validate(t, q2)
	check.Equal(t, 42).Assert(q2.Dequeue())
	check.String(t, ``).Assert(q2)
	validate(t, q2)

	q2.EnqueueFrom(enumerator.Range(4, 3))
	check.String(t, `4, 5, 6`).Assert(q2)
	validate(t, q2)

	q2.EnqueueFrom(enumerator.Range(9, 3))
	check.String(t, `4, 5, 6, 9, 10, 11`).Assert(q2)
	validate(t, q2)

	q2.Clip() // no effect
	check.String(t, `4, 5, 6, 9, 10, 11`).Assert(q2)
	validate(t, q2)
}

func Test_Queue_New(t *testing.T) {
	q := New[int]()
	check.Empty(t).Assert(q)
	check.String(t, ``).Assert(q)

	q = New[int](5)
	check.Length(t, 5).Assert(q)
	check.String(t, `0, 0, 0, 0, 0`).Assert(q)

	check.MatchError(t, `^invalid number of arguments \{count: 2, maximum: 1, usage: size\}$`).
		Panic(func() { New[int](1, 2) })

	q = New[int](-1)
	check.Empty(t).Assert(q)
	check.String(t, ``).Assert(q)

	q = Fill[int](8, 5)
	check.Length(t, 5).Assert(q)
	check.String(t, `8, 8, 8, 8, 8`).Assert(q)

	q = Fill[int](8, -1)
	check.Empty(t).Assert(q)
	check.String(t, ``).Assert(q)

	q = From[int](nil)
	check.Empty(t).Assert(q)
	check.String(t, ``).Assert(q)
}

func Test_Queue_UnstableIteration(t *testing.T) {
	q := With(2, 3, 4, 5, 6)
	it := q.Enumerate().Iterate()
	check.True(t).Assert(it.Next())
	check.Equal(t, 2).Assert(it.Current())

	q.Enqueue(8)
	check.True(t).Assert(it.Next())
	check.Equal(t, 3).Assert(it.Current())

	check.Equal(t, 2).Assert(q.Dequeue())
	check.MatchError(t, `Collection was modified; iteration may not continue`).
		Panic(func() { it.Next() })
}
