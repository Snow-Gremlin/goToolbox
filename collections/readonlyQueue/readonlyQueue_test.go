package readonlyQueue

import (
	"fmt"
	"slices"
	"testing"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/collections/list"
	"github.com/Snow-Gremlin/goToolbox/events"
	"github.com/Snow-Gremlin/goToolbox/testers/check"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

type pseudoQueueImp[T any] struct {
	q []T
	e events.Event[collections.ChangeArgs]
}

func (q *pseudoQueueImp[T]) Enumerate() collections.Enumerator[T] {
	return enumerator.Enumerate(q.q...)
}

func (q *pseudoQueueImp[T]) Empty() bool {
	return len(q.q) <= 0
}

func (q *pseudoQueueImp[T]) Count() int {
	return len(q.q)
}

func (q *pseudoQueueImp[T]) String() string {
	return fmt.Sprint(q.q)
}

func (q *pseudoQueueImp[T]) Equals(other any) bool {
	v, ok := other.(collections.Sliceable[T])
	return ok && utils.Equal(q.ToSlice(), v.ToSlice())
}

func (q *pseudoQueueImp[T]) ToSlice() []T {
	return slices.Clone(q.q)
}

func (q *pseudoQueueImp[T]) CopyToSlice(sc []T) {
	copy(sc, q.ToSlice())
}

func (q *pseudoQueueImp[T]) ToList() collections.List[T] {
	return list.From(q.Enumerate())
}

func (q *pseudoQueueImp[T]) Peek() T {
	return q.q[0]
}

func (q *pseudoQueueImp[T]) TryPeek() (T, bool) {
	if q.Empty() {
		return utils.Zero[T](), false
	}
	return q.q[0], true
}

func (q *pseudoQueueImp[T]) OnChange() events.Event[collections.ChangeArgs] {
	return q.e
}

func Test_ReadonlyQueue(t *testing.T) {
	q0 := &pseudoQueueImp[int]{
		q: []int{1, 2, 3},
		e: nil,
	}
	q1 := New(q0)
	q2 := &pseudoQueueImp[int]{
		q: []int{1, 2, 3},
		e: nil,
	}
	q3 := New(q2)
	check.False(t).Assert(q1.Empty())
	check.Length(t, 3).Assert(q1)
	check.Equal(t, []int{1, 2, 3}).Assert(q1.Enumerate().ToSlice())
	check.String(t, `[1 2 3]`).Assert(q1)
	check.Equal(t, q3).Assert(q1)

	q0.q = append(q0.q, 34)
	check.Length(t, 4).Assert(q1)
	check.String(t, `[1 2 3 34]`).Assert(q1)
	check.Equal(t, []int{1, 2, 3, 34}).Assert(q1.ToSlice())
	check.String(t, `1, 2, 3, 34`).Assert(q1.ToList())

	p := make([]int, 5)
	q1.CopyToSlice(p)
	check.Equal(t, []int{1, 2, 3, 34, 0}).Assert(p)

	check.Equal(t, 1).Assert(q1.Peek())
	v, ok := q1.TryPeek()
	check.Equal(t, 1).Assert(v)
	check.True(t).Assert(ok)
	check.NotEqual(t, q3).Assert(q1)

	check.Same(t, q0.OnChange()).Assert(q1.OnChange())
	check.Same(t, q2.OnChange()).Assert(q3.OnChange())
}
