package stack

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/collections/iterator"
	"github.com/Snow-Gremlin/goToolbox/collections/list"
	"github.com/Snow-Gremlin/goToolbox/collections/readonlyStack"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

type (
	node[T any] struct {
		value T
		prev  *node[T]
	}

	stackImp[T any] struct {
		count     int
		head      *node[T]
		enumGuard uint
	}
)

func newImp[T any]() *stackImp[T] {
	return &stackImp[T]{
		count:     0,
		head:      nil,
		enumGuard: 0,
	}
}

func (s *stackImp[T]) Enumerate() collections.Enumerator[T] {
	return enumerator.New(func() collections.Iterator[T] {
		n := s.head
		guardStash := s.enumGuard
		return iterator.New(func() (T, bool) {
			if n == nil {
				return utils.Zero[T](), false
			}
			if guardStash != s.enumGuard {
				// Only removing nodes disrupts the iterations.
				// However, adding nodes doesn't cause a problem.
				panic(terror.UnstableIteration())
			}
			value := n.value
			n = n.prev
			return value, true
		})
	})
}

func (s *stackImp[T]) Empty() bool {
	return s.head == nil
}

func (s *stackImp[T]) Count() int {
	return s.count
}

func (s *stackImp[T]) String() string {
	return s.Enumerate().Join(`, `)
}

func (s *stackImp[T]) ToSlice() []T {
	return s.Enumerate().ToSlice()
}

func (s *stackImp[T]) CopyToSlice(s2 []T) {
	s.Enumerate().CopyToSlice(s2)
}

func (q *stackImp[T]) ToList() collections.List[T] {
	return list.From(q.Enumerate())
}

func (s *stackImp[T]) Peek() T {
	if s.head != nil {
		return s.head.value
	}
	panic(terror.EmptyCollection(`Peek`))
}

func (s *stackImp[T]) TryPeek() (T, bool) {
	if s.head != nil {
		return s.head.value, true
	}
	return utils.Zero[T](), false
}

func (s *stackImp[T]) pushOne(value T) {
	s.head = &node[T]{
		value: value,
		prev:  s.head,
	}
	s.count++
}

func (s *stackImp[T]) popOne() T {
	v := s.head.value
	s.head = s.head.prev
	s.count--
	s.enumGuard++
	return v
}

func (s *stackImp[T]) Push(values ...T) {
	for i := len(values) - 1; i >= 0; i-- {
		s.pushOne(values[i])
	}
}

func (s *stackImp[T]) PushFrom(e collections.Enumerator[T]) {
	if utils.IsNil(e) {
		return
	}

	it := e.Iterate()
	if !it.Next() {
		return
	}

	newHead := &node[T]{
		value: it.Current(),
		prev:  nil,
	}
	prev := newHead
	count := 1
	for it.Next() {
		n := &node[T]{
			value: it.Current(),
			prev:  nil,
		}
		prev.prev = n
		prev = n
		count++
	}

	prev.prev = s.head
	s.head = newHead
	s.count += count
}

func (s *stackImp[T]) Take(count int) []T {
	count = min(count, s.count)
	result := make([]T, count)
	for i := 0; i < count; i++ {
		result[i] = s.popOne()
	}
	return result
}

func (s *stackImp[T]) Pop() T {
	if s.head == nil {
		panic(terror.EmptyCollection(`Pop`))
	}
	return s.popOne()
}

func (s *stackImp[T]) TryPop() (T, bool) {
	if s.head == nil {
		return utils.Zero[T](), false
	}
	return s.popOne(), true
}

func (s *stackImp[T]) TrimTo(count int) {
	if count <= 0 {
		s.head = nil
		s.count = 0
		return
	}

	prev := s.head
	for i := 1; i < count; i++ {
		prev = prev.prev
		if prev == nil {
			return
		}
	}
	prev.prev = nil
	s.count = count
	s.enumGuard++
}

func (s *stackImp[T]) Clear() {
	s.head = nil
	s.count = 0
	s.enumGuard++
}

func (s *stackImp[T]) Clip() {
	// no effect
}

func (s *stackImp[T]) Equals(other any) bool {
	s2, ok := other.(collections.Collection[T])
	return ok && s.count == s2.Count() &&
		s.Enumerate().Equals(s2.Enumerate())
}

func (s *stackImp[T]) Clone() collections.Stack[T] {
	return From(s.Enumerate())
}

func (s *stackImp[T]) Readonly() collections.ReadonlyStack[T] {
	return readonlyStack.New(s)
}
