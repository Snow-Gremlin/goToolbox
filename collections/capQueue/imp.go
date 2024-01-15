package capQueue

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

// growthRate is the number of nodes to add when
// a new node is needed and the graveyard is empty.
const growthRate = 10

type (
	node[T any] struct {
		value T
		next  *node[T]
	}

	capQueueImp[T any] struct {
		count     int
		head      *node[T]
		tail      *node[T]
		graveyard *node[T]
		enumGuard uint
		event     events.Event[collections.ChangeArgs]
	}
)

func newImp[T any]() *capQueueImp[T] {
	return &capQueueImp[T]{
		count:     0,
		head:      nil,
		tail:      nil,
		graveyard: nil,
		enumGuard: 0,
		event:     nil,
	}
}

func (q *capQueueImp[T]) newNode(value T) *node[T] {
	if q.graveyard == nil {
		q.addTombs(growthRate)
	}

	n := q.graveyard
	q.graveyard = n.next
	n.next = nil
	n.value = value
	return n
}

func (q *capQueueImp[T]) tombs() int {
	count := 0
	n := q.graveyard
	for n != nil {
		count++
		n = n.next
	}
	return count
}

func (q *capQueueImp[T]) growCap(capacity int) {
	capacity -= q.count
	if capacity <= 0 {
		return
	}

	capacity -= q.tombs()
	if capacity <= 0 {
		return
	}

	q.addTombs(capacity)
}

func (q *capQueueImp[T]) addTombs(count int) {
	c := make([]node[T], count)
	prev := q.graveyard
	for i := count - 1; i >= 0; i-- {
		n := &c[i]
		n.next = prev
		prev = n
	}
	q.graveyard = &c[0]
}

func (q *capQueueImp[T]) onEnqueued() {
	if q.event != nil {
		q.event.Invoke(changeArgs.NewAdded())
	}
}

func (q *capQueueImp[T]) onDequeued() {
	if q.event != nil {
		q.event.Invoke(changeArgs.NewRemoved())
	}
}

func (q *capQueueImp[T]) Enumerate() collections.Enumerator[T] {
	return enumerator.New(func() collections.Iterator[T] {
		n := q.head
		guardStash := q.enumGuard
		return iterator.New(func() (T, bool) {
			if n == nil {
				return utils.Zero[T](), false
			}
			if guardStash != q.enumGuard {
				// Only removing nodes disrupts the iterations because those
				// nodes are moved to the graveyard. Continuing would then be
				// moving through the graveyard which we don't want to do.
				// However, adding nodes doesn't cause a problem.
				panic(terror.UnstableIteration())
			}
			value := n.value
			n = n.next
			return value, true
		})
	})
}

func (q *capQueueImp[T]) Empty() bool {
	return q.count <= 0
}

func (q *capQueueImp[T]) Count() int {
	return q.count
}

func (q *capQueueImp[T]) String() string {
	return q.Enumerate().Join(`, `)
}

func (q *capQueueImp[T]) ToSlice() []T {
	return q.Enumerate().ToSlice()
}

func (q *capQueueImp[T]) CopyToSlice(s []T) {
	q.Enumerate().CopyToSlice(s)
}

func (q *capQueueImp[T]) ToList() collections.List[T] {
	return list.From(q.Enumerate())
}

func (q *capQueueImp[T]) Peek() T {
	if q.head != nil {
		return q.head.value
	}
	panic(terror.EmptyCollection(`Peek`))
}

func (q *capQueueImp[T]) TryPeek() (T, bool) {
	if q.head != nil {
		return q.head.value, true
	}
	return utils.Zero[T](), false
}

func (q *capQueueImp[T]) OnChange() events.Event[collections.ChangeArgs] {
	if q.event == nil {
		q.event = event.New[collections.ChangeArgs]()
	}
	return q.event
}

func (q *capQueueImp[T]) Enqueue(values ...T) {
	count := len(values)
	if count <= 0 {
		return
	}

	prev := q.newNode(values[0])
	if q.tail == nil {
		q.head = prev
	} else {
		q.tail.next = prev
	}

	for i := 1; i < count; i++ {
		n := q.newNode(values[i])
		prev.next = n
		prev = n
	}

	q.tail = prev
	q.count += count
	q.onEnqueued()
}

func (q *capQueueImp[T]) EnqueueFrom(e collections.Enumerator[T]) {
	if utils.IsNil(e) {
		return
	}

	it := e.Iterate()
	if !it.Next() {
		return
	}

	count := 1
	first := q.newNode(it.Current())
	prev := first

	for it.Next() {
		n := q.newNode(it.Current())
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

func (q *capQueueImp[T]) Take(count int) []T {
	count = min(count, q.count)
	if count <= 0 {
		return []T{}
	}
	result := make([]T, count)
	n := q.head
	z := utils.Zero[T]()
	p := n
	for i := 0; i < count; i++ {
		result[i] = n.value
		n.value = z
		p = n
		n = n.next
	}
	p.next = q.graveyard
	q.graveyard = q.head

	q.head = n
	if q.head == nil {
		q.tail = nil
	}
	q.count -= count
	q.enumGuard++
	q.onDequeued()
	return result
}

func (q *capQueueImp[T]) Dequeue() T {
	if v, ok := q.TryDequeue(); ok {
		return v
	}
	panic(terror.EmptyCollection(`Dequeue`))
}

func (q *capQueueImp[T]) TryDequeue() (T, bool) {
	n := q.head
	if n == nil {
		return utils.Zero[T](), false
	}

	v := n.value
	q.head = n.next
	if q.head == nil {
		q.tail = nil
	}
	q.count--

	n.value = utils.Zero[T]()
	n.next = q.graveyard
	q.graveyard = n
	q.enumGuard++
	q.onDequeued()
	return v, true
}

func (q *capQueueImp[T]) Clear() {
	if q.head == nil {
		return
	}

	z := utils.Zero[T]()
	for n := q.head; n != nil; n = n.next {
		n.value = z
	}
	q.tail.next = q.graveyard
	q.graveyard = q.head

	q.head = nil
	q.tail = nil
	q.count = 0
	q.enumGuard++
	q.onDequeued()
}

func (q *capQueueImp[T]) Clip() {
	q.graveyard = nil
}

func (q *capQueueImp[T]) Equals(other any) bool {
	s, ok := other.(collections.Collection[T])
	return ok && q.count == s.Count() &&
		q.Enumerate().Equals(s.Enumerate())
}

func (q *capQueueImp[T]) Clone() collections.Queue[T] {
	return From(q.Enumerate())
}

func (q *capQueueImp[T]) Readonly() collections.ReadonlyQueue[T] {
	return readonlyQueue.New(q)
}
