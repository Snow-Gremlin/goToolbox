package set

import (
	"slices"
	"strings"

	"goToolbox/collections"
	"goToolbox/collections/enumerator"
	"goToolbox/collections/iterator"
	"goToolbox/collections/list"
	"goToolbox/collections/readonlySet"
	"goToolbox/internal/simpleSet"
	"goToolbox/utils"
)

type setImp[T comparable] struct {
	m simpleSet.Set[T]
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

func (q *setImp[T]) ToList() collections.List[T] {
	return list.From(q.Enumerate())
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

func (s *setImp[T]) Add(values ...T) bool {
	added := false
	for _, key := range values {
		added = s.m.SetTest(key) || added
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
	return added
}

func (s *setImp[T]) Remove(values ...T) bool {
	removed := false
	for _, key := range values {
		removed = s.m.RemoveTest(key) || removed
	}
	return removed
}

func (s *setImp[T]) RemoveIf(predicate collections.Predicate[T]) bool {
	return s.m.RemoveIf(predicate)
}

func (s *setImp[T]) Clear() {
	s.m = simpleSet.New[T]()
}

func (s *setImp[T]) Clone() collections.Set[T] {
	return &setImp[T]{
		m: s.m.Clone(),
	}
}

func (s *setImp[T]) Readonly() collections.ReadonlySet[T] {
	return readonlySet.New(s)
}
