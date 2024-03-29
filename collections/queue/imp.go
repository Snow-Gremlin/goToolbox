package queue

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/changeArgs"
	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/collections/iterator"
	"github.com/Snow-Gremlin/goToolbox/collections/list"
	"github.com/Snow-Gremlin/goToolbox/collections/readonlyQueue"
	"github.com/Snow-Gremlin/goToolbox/events"
	"github.com/Snow-Gremlin/goToolbox/events/event"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

type (
	node[T any] struct {
		value T
		next  *node[T]
	}

	queueImp[T any] struct {
		count     int
		head      *node[T]
		tail      *node[T]
		enumGuard uint
		event     events.Event[collections.ChangeArgs]
	}
)

func newNode[T any](value T) *node[T] {
	return &node[T]{
		value: value,
		next:  nil,
	}
}

func newImp[T any]() *queueImp[T] {
	return &queueImp[T]{
		count:     0,
		head:      nil,
		tail:      nil,
		enumGuard: 0,
		event:     nil,
	}
}

func (q *queueImp[T]) onEnqueued() {
	if q.event != nil {
		q.event.Invoke(changeArgs.NewAdded())
	}
}

func (q *queueImp[T]) onDequeued() {
	if q.event != nil {
		q.event.Invoke(changeArgs.NewRemoved())
	}
}

func (q *queueImp[T]) Enumerate() collections.Enumerator[T] {
	return enumerator.New(func() collections.Iterator[T] {
		n := q.head
		guardStash := q.enumGuard
		return iterator.New(func() (T, bool) {
			if n == nil {
				return utils.Zero[T](), false
			}
			if guardStash != q.enumGuard {
				// Only removing nodes disrupts the iterations.
				// However, adding nodes doesn't cause a problem.
				panic(terror.UnstableIteration())
			}
			value := n.value
			n = n.next
			return value, true
		})
	})
}

func (q *queueImp[T]) Empty() bool {
	return q.count <= 0
}

func (q *queueImp[T]) Count() int {
	return q.count
}

func (q *queueImp[T]) String() string {
	return q.Enumerate().Join(`, `)
}

func (q *queueImp[T]) ToSlice() []T {
	return q.Enumerate().ToSlice()
}

func (q *queueImp[T]) CopyToSlice(s []T) {
	q.Enumerate().CopyToSlice(s)
}

func (q *queueImp[T]) ToList() collections.List[T] {
	return list.From(q.Enumerate())
}

func (q *queueImp[T]) Peek() T {
	if q.head != nil {
		return q.head.value
	}
	panic(terror.EmptyCollection(`Peek`))
}

func (q *queueImp[T]) TryPeek() (T, bool) {
	if q.head != nil {
		return q.head.value, true
	}
	return utils.Zero[T](), false
}

func (q *queueImp[T]) OnChange() events.Event[collections.ChangeArgs] {
	if q.event == nil {
		q.event = event.New[collections.ChangeArgs]()
	}
	return q.event
}

func (q *queueImp[T]) Enqueue(values ...T) {
	count := len(values)
	if count <= 0 {
		return
	}

	prev := newNode(values[0])
	if q.tail == nil {
		q.head = prev
	} else {
		q.tail.next = prev
	}

	for i := 1; i < count; i++ {
		n := newNode(values[i])
		prev.next = n
		prev = n
	}

	q.tail = prev
	q.count += count
	q.onEnqueued()
}

func (q *queueImp[T]) EnqueueFrom(e collections.Enumerator[T]) {
	if utils.IsNil(e) {
		return
	}

	it := e.Iterate()
	if !it.Next() {
		return
	}

	count := 1
	first := newNode(it.Current())
	prev := first

	for it.Next() {
		n := newNode(it.Current())
		prev.next = n
		prev = n
		count++
	}

	if q.tail != nil {
		q.tail.next = first
	} else {
		q.head = first
	}
	q.tail = prev
	q.count += count
	q.onEnqueued()
}

func (q *queueImp[T]) Take(count int) []T {
	count = min(count, q.count)
	if count <= 0 {
		return []T{}
	}
	result := make([]T, count)
	n := q.head
	for i := 0; i < count; i++ {
		result[i] = n.value
		n = n.next
	}
	q.head = n
	if q.head == nil {
		q.tail = nil
	}
	q.count -= count
	q.enumGuard++
	q.onDequeued()
	return result
}

func (q *queueImp[T]) Dequeue() T {
	if v, ok := q.TryDequeue(); ok {
		return v
	}
	panic(terror.EmptyCollection(`Dequeue`))
}

func (q *queueImp[T]) TryDequeue() (T, bool) {
	if q.head == nil {
		return utils.Zero[T](), false
	}
	v := q.head.value
	q.head = q.head.next
	if q.head == nil {
		q.tail = nil
	}
	q.count--
	q.enumGuard++
	q.onDequeued()
	return v, true
}

func (q *queueImp[T]) Clear() {
	if q.count > 0 {
		q.head = nil
		q.tail = nil
		q.count = 0
		q.enumGuard++
		q.onDequeued()
	}
}

func (q *queueImp[T]) Clip() {
	// no effect
}

func (q *queueImp[T]) Equals(other any) bool {
	s, ok := other.(collections.Collection[T])
	return ok && q.count == s.Count() &&
		q.Enumerate().Equals(s.Enumerate())
}

func (q *queueImp[T]) Clone() collections.Queue[T] {
	return From(q.Enumerate())
}

func (q *queueImp[T]) Readonly() collections.ReadonlyQueue[T] {
	return readonlyQueue.New(q)
}
