package set

import (
	"slices"
	"strings"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/changeArgs"
	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/collections/iterator"
	"github.com/Snow-Gremlin/goToolbox/collections/list"
	"github.com/Snow-Gremlin/goToolbox/collections/readonlySet"
	"github.com/Snow-Gremlin/goToolbox/events"
	"github.com/Snow-Gremlin/goToolbox/events/event"
	"github.com/Snow-Gremlin/goToolbox/internal/simpleSet"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

type setImp[T comparable] struct {
	m     simpleSet.Set[T]
	event events.Event[collections.ChangeArgs]
}

func (s *setImp[T]) onAdded() {
	if s.event != nil {
		s.event.Invoke(changeArgs.NewAdded())
	}
}

func (s *setImp[T]) onRemoved() {
	if s.event != nil {
		s.event.Invoke(changeArgs.NewRemoved())
	}
}

func (s *setImp[T]) Enumerate() collections.Enumerator[T] {
	// Since Go randomizes the order of values, to keep a consistent
	// iteration, all the keys must be collected once before iteration.
	// The keys will still be in random order but consistent.
	// Because the keys are collected ahead of iterations, changes to
	// the set may just cause the enumeration to be unstable
	// but doesn't require it to be stopped.
	return enumerator.New(func() collections.Iterator[T] {
		values := s.m.ToSlice()
		index, count := -1, len(values)-1
		return iterator.New(func() (T, bool) {
			if index < count {
				index++
				return values[index], true
			}
			return utils.Zero[T](), false
		})
	})
}

func (s *setImp[T]) Empty() bool {
	return s.m.Count() <= 0
}

func (s *setImp[T]) Count() int {
	return s.m.Count()
}

func (s *setImp[T]) ToSlice() []T {
	return s.m.ToSlice()
}

func (s *setImp[T]) CopyToSlice(s2 []T) {
	index, room := 0, len(s2)
	for value := range s.m {
		if index >= room {
			break
		}
		s2[index] = value
		index++
	}
}

func (s *setImp[T]) ToList() collections.List[T] {
	return list.From(s.Enumerate())
}

func (s *setImp[T]) Contains(value T) bool {
	return s.m.Has(value)
}

func (s *setImp[T]) String() string {
	parts := utils.Strings(s.m.ToSlice())
	slices.Sort(parts)
	return strings.Join(parts, `, `)
}

func (s *setImp[T]) Equals(other any) bool {
	s2, ok := other.(collections.Collection[T])
	if !ok || s.Count() != s2.Count() {
		return false
	}

	it := s2.Enumerate().Iterate()
	for it.Next() {
		if !s.Contains(it.Current()) {
			return false
		}
	}
	return true
}

func (s *setImp[T]) OnChange() events.Event[collections.ChangeArgs] {
	if s.event == nil {
		s.event = event.New[collections.ChangeArgs]()
	}
	return s.event
}

func (s *setImp[T]) Add(values ...T) bool {
	added := false
	for _, key := range values {
		added = s.m.SetTest(key) || added
	}
	if added {
		s.onAdded()
	}
	return added
}

func (s *setImp[T]) AddFrom(e collections.Enumerator[T]) bool {
	if utils.IsNil(e) {
		return false
	}
	added := false
	it := e.Iterate()
	for it.Next() {
		added = s.m.SetTest(it.Current()) || added
	}
	if added {
		s.onAdded()
	}
	return added
}

func (s *setImp[T]) TakeAny() T {
	for value := range s.m {
		delete(s.m, value)
		s.onRemoved()
		return value
	}
	panic(terror.EmptyCollection(`TakeAny`))
}

func (s *setImp[T]) TakeMany(count int) []T {
	count = min(count, s.Count())
	if count <= 0 {
		return []T{}
	}
	index := 0
	results := make([]T, count)
	for value := range s.m {
		if index >= count {
			break
		}
		results[index] = value
		index++
		delete(s.m, value)
	}
	s.onRemoved()
	return results
}

func (s *setImp[T]) Remove(values ...T) bool {
	removed := false
	for _, key := range values {
		removed = s.m.RemoveTest(key) || removed
	}
	if removed {
		s.onRemoved()
	}
	return removed
}

func (s *setImp[T]) RemoveIf(predicate collections.Predicate[T]) bool {
	if s.m.RemoveIf(predicate) {
		s.onRemoved()
		return true
	}
	return false
}

func (s *setImp[T]) Clear() {
	if len(s.m) > 0 {
		s.m = simpleSet.New[T]()
		s.onRemoved()
	}
}

func (s *setImp[T]) Clone() collections.Set[T] {
	return &setImp[T]{
		m:     s.m.Clone(),
		event: nil,
	}
}

func (s *setImp[T]) Readonly() collections.ReadonlySet[T] {
	return readonlySet.New(s)
}
