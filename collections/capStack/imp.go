package capStack

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/changeArgs"
	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/collections/iterator"
	"github.com/Snow-Gremlin/goToolbox/collections/list"
	"github.com/Snow-Gremlin/goToolbox/collections/readonlyStack"
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
		prev  *node[T]
	}

	capStackImp[T any] struct {
		count     int
		head      *node[T]
		graveyard *node[T]
		enumGuard uint
		event     events.Event[collections.ChangeArgs]
	}
)

func newImp[T any]() *capStackImp[T] {
	return &capStackImp[T]{
		count:     0,
		head:      nil,
		graveyard: nil,
		enumGuard: 0,
		event:     nil,
	}
}

func (s *capStackImp[T]) newNode(value T) *node[T] {
	if s.graveyard == nil {
		s.addTombs(growthRate)
	}

	n := s.graveyard
	s.graveyard = n.prev
	n.prev = nil
	n.value = value
	return n
}

func (s *capStackImp[T]) tombs() int {
	count := 0
	n := s.graveyard
	for n != nil {
		count++
		n = n.prev
	}
	return count
}

func (s *capStackImp[T]) growCap(capacity int) {
	capacity -= s.count
	if capacity <= 0 {
		return
	}

	capacity -= s.tombs()
	if capacity <= 0 {
		return
	}

	s.addTombs(capacity)
}

func (s *capStackImp[T]) addTombs(count int) {
	c := make([]node[T], count)
	prev := s.graveyard
	for i := count - 1; i >= 0; i-- {
		n := &c[i]
		n.prev = prev
		prev = n
	}
	s.graveyard = &c[0]
}

func (s *capStackImp[T]) entombFrom(n *node[T]) {
	z := utils.Zero[T]()
	if n != nil {
		g := n
		g.value = z
		for g.prev != nil {
			g = g.prev
			g.value = z
		}
		g.prev = s.graveyard
	}
	s.graveyard = n
	s.enumGuard++
}

func (s *capStackImp[T]) pushOne(value T) {
	n := s.newNode(value)
	n.prev = s.head
	s.head = n
	s.count++
}

func (s *capStackImp[T]) popOne() T {
	n := s.head
	v := n.value
	s.head = n.prev
	s.count--

	n.prev = s.graveyard
	n.value = utils.Zero[T]()
	s.graveyard = n
	s.enumGuard++

	return v
}

func (s *capStackImp[T]) onPushed() {
	if s.event != nil {
		s.event.Invoke(changeArgs.NewAdded())
	}
}

func (s *capStackImp[T]) onPopped() {
	if s.event != nil {
		s.event.Invoke(changeArgs.NewRemoved())
	}
}

func (s *capStackImp[T]) Enumerate() collections.Enumerator[T] {
	return enumerator.New(func() collections.Iterator[T] {
		n := s.head
		guardStash := s.enumGuard
		return iterator.New(func() (T, bool) {
			if n == nil {
				return utils.Zero[T](), false
			}
			if guardStash != s.enumGuard {
				// Only removing nodes disrupts the iterations because those
				// nodes are moved to the graveyard. Continuing would then be
				// moving through the graveyard which we don't want to do.
				// However, adding nodes doesn't cause a problem.
				panic(terror.UnstableIteration())
			}
			value := n.value
			n = n.prev
			return value, true
		})
	})
}

func (s *capStackImp[T]) Empty() bool {
	return s.head == nil
}

func (s *capStackImp[T]) Count() int {
	return s.count
}

func (s *capStackImp[T]) String() string {
	return s.Enumerate().Join(`, `)
}

func (s *capStackImp[T]) ToSlice() []T {
	return s.Enumerate().ToSlice()
}

func (s *capStackImp[T]) CopyToSlice(s2 []T) {
	s.Enumerate().CopyToSlice(s2)
}

func (q *capStackImp[T]) ToList() collections.List[T] {
	return list.From(q.Enumerate())
}

func (s *capStackImp[T]) Peek() T {
	if s.head != nil {
		return s.head.value
	}
	panic(terror.EmptyCollection(`Peek`))
}

func (s *capStackImp[T]) TryPeek() (T, bool) {
	if s.head != nil {
		return s.head.value, true
	}
	return utils.Zero[T](), false
}

func (s *capStackImp[T]) OnChange() events.Event[collections.ChangeArgs] {
	if s.event == nil {
		s.event = event.New[collections.ChangeArgs]()
	}
	return s.event
}

func (s *capStackImp[T]) Push(values ...T) {
	if length := len(values); length > 0 {
		for i := length - 1; i >= 0; i-- {
			s.pushOne(values[i])
		}
		s.onPushed()
	}
}

func (s *capStackImp[T]) PushFrom(e collections.Enumerator[T]) {
	if utils.IsNil(e) {
		return
	}

	it := e.Iterate()
	if !it.Next() {
		return
	}

	newHead := s.newNode(it.Current())
	prev := newHead
	count := 1
	for it.Next() {
		n := s.newNode(it.Current())
		prev.prev = n
		prev = n
		count++
	}

	prev.prev = s.head
	s.head = newHead
	s.count += count
	s.onPushed()
}

func (s *capStackImp[T]) Take(count int) []T {
	count = min(count, s.count)
	if count <= 0 {
		return []T{}
	}
	result := make([]T, count)
	for i := 0; i < count; i++ {
		result[i] = s.popOne()
	}
	s.onPopped()
	return result
}

func (s *capStackImp[T]) Pop() T {
	if v, ok := s.TryPop(); ok {
		return v
	}
	panic(terror.EmptyCollection(`Pop`))
}

func (s *capStackImp[T]) TryPop() (T, bool) {
	if s.head == nil {
		return utils.Zero[T](), false
	}
	v := s.popOne()
	s.onPopped()
	return v, true
}

func (s *capStackImp[T]) TrimTo(count int) {
	if count <= 0 {
		s.Clear()
		return
	}

	prev := s.head
	for i := 1; i < count; i++ {
		if prev == nil {
			return
		}
		prev = prev.prev
	}
	if prev.prev == nil {
		return
	}

	s.entombFrom(prev.prev)
	prev.prev = nil
	s.count = count
	s.onPopped()
}

func (s *capStackImp[T]) Clear() {
	if s.head == nil {
		return
	}

	s.entombFrom(s.head)

	s.head = nil
	s.count = 0
	s.onPopped()
}

func (s *capStackImp[T]) Clip() {
	s.graveyard = nil
}

func (s *capStackImp[T]) Equals(other any) bool {
	s2, ok := other.(collections.Collection[T])
	return ok && s.count == s2.Count() &&
		s.Enumerate().Equals(s2.Enumerate())
}

func (s *capStackImp[T]) Clone() collections.Stack[T] {
	return From(s.Enumerate())
}

func (s *capStackImp[T]) Readonly() collections.ReadonlyStack[T] {
	return readonlyStack.New(s)
}
