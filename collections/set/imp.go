package set

import (
	"maps"
	"slices"
	"strings"

	"goToolbox/collections"
	"goToolbox/collections/enumerator"
	"goToolbox/collections/iterator"
	"goToolbox/collections/list"
	"goToolbox/collections/readonlySet"
	"goToolbox/utils"
)

type setImp[T comparable] struct {
	m map[T]struct{}
}

func (s *setImp[T]) Enumerate() collections.Enumerator[T] {
	// Since Go randomizes the order of values, to keep a consistent
	// iteration, all the keys must be collected once before iteration.
	// The keys will still be in random order but consistent.
	// Because the keys are collected ahead of iterations, changes to
	// the set may just cause the enumeration to be unstable
	// but doesn't require it to be stopped.
	return enumerator.New(func() collections.Iterator[T] {
		keys := utils.Keys(s.m)
		index, count := -1, len(keys)-1
		return iterator.New(func() (T, bool) {
			if index < count {
				index++
				return keys[index], true
			}
			return utils.Zero[T](), false
		})
	})
}

func (s *setImp[T]) Empty() bool {
	return len(s.m) <= 0
}

func (s *setImp[T]) Count() int {
	return len(s.m)
}

func (s *setImp[T]) ToSlice() []T {
	return utils.Keys(s.m)
}

func (s *setImp[T]) CopyToSlice(s2 []T) {
	index, room := 0, len(s2)
	for key := range s.m {
		if index >= room {
			break
		}
		s2[index] = key
		index++
	}
}

func (q *setImp[T]) ToList() collections.List[T] {
	return list.From(q.Enumerate())
}

func (s *setImp[T]) Contains(key T) bool {
	_, ok := s.m[key]
	return ok
}

func (s *setImp[T]) String() string {
	parts := utils.Strings(utils.Keys(s.m))
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
		if _, ok := s.m[key]; !ok {
			s.m[key] = struct{}{}
			added = true
		}
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
		key := it.Current()
		if _, ok := s.m[key]; !ok {
			s.m[key] = struct{}{}
			added = true
		}
	}
	return added
}

func (s *setImp[T]) Remove(values ...T) bool {
	removed := false
	for _, key := range values {
		if _, ok := s.m[key]; ok {
			delete(s.m, key)
			removed = true
		}

	}
	return removed
}

func (s *setImp[T]) RemoveIf(handle collections.Predicate[T]) bool {
	priorCount := s.Count()
	maps.DeleteFunc(s.m, func(key T, _ struct{}) bool { return handle(key) })
	return s.Count() != priorCount
}

func (s *setImp[T]) Clear() {
	s.m = make(map[T]struct{})
}

func (s *setImp[T]) Clone() collections.Set[T] {
	return &setImp[T]{
		m: maps.Clone(s.m),
	}
}

func (s *setImp[T]) Readonly() collections.ReadonlySet[T] {
	return readonlySet.New(s)
}
