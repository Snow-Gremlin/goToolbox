package sortedSet

import (
	"slices"
	"strings"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/changeArgs"
	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/collections/iterator"
	"github.com/Snow-Gremlin/goToolbox/collections/list"
	"github.com/Snow-Gremlin/goToolbox/collections/readonlySortedSet"
	"github.com/Snow-Gremlin/goToolbox/comp"
	"github.com/Snow-Gremlin/goToolbox/events"
	"github.com/Snow-Gremlin/goToolbox/events/event"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

type sortedSetImp[T any] struct {
	data     []T
	comparer comp.Comparer[T]
	event    events.Event[collections.ChangeArgs]
}

func (s *sortedSetImp[T]) find(value T) (int, bool) {
	return slices.BinarySearchFunc(s.data, value, s.comparer)
}

func (s *sortedSetImp[T]) grow(newLength int) {
	s.data = slices.Grow(s.data, newLength)
}

func (s *sortedSetImp[T]) addOne(value T, force bool) (T, bool) {
	index, found := s.find(value)
	if found {
		if force {
			s.data[index] = value
			return value, false
		}
		return s.data[index], false
	}
	s.data = slices.Insert(s.data, index, value)
	return value, true
}

func (s *sortedSetImp[T]) onAdded() {
	if s.event != nil {
		s.event.Invoke(changeArgs.NewAdded())
	}
}

func (s *sortedSetImp[T]) onRemoved() {
	if s.event != nil {
		s.event.Invoke(changeArgs.NewRemoved())
	}
}

func (s *sortedSetImp[T]) Enumerate() collections.Enumerator[T] {
	// Since we can use the length to keep the index valid
	// changes to the list don't have stop enumerators.
	// Changes may just cause the enumeration to be unstable.
	return enumerator.New(func() collections.Iterator[T] {
		index := -1
		return iterator.New(func() (T, bool) {
			if index < len(s.data)-1 {
				index++
				return s.data[index], true
			}
			return utils.Zero[T](), false
		})
	})
}

func (s *sortedSetImp[T]) Backwards() collections.Enumerator[T] {
	// See comment in Enumerate
	return enumerator.New(func() collections.Iterator[T] {
		index := len(s.data)
		return iterator.New(func() (T, bool) {
			if index = min(index, len(s.data)); index > 0 {
				index--
				return s.data[index], true
			}
			return utils.Zero[T](), false
		})
	})
}

func (s *sortedSetImp[T]) Empty() bool {
	return len(s.data) <= 0
}

func (s *sortedSetImp[T]) Count() int {
	return len(s.data)
}

func (s *sortedSetImp[T]) ToSlice() []T {
	return slices.Clip(s.data)
}

func (s *sortedSetImp[T]) CopyToSlice(s2 []T) {
	room := len(s2)
	for index, value := range s.data {
		if index >= room {
			break
		}
		s2[index] = value
	}
}

func (s *sortedSetImp[T]) ToList() collections.List[T] {
	return list.From(s.Enumerate())
}

func (s *sortedSetImp[T]) Contains(value T) bool {
	_, found := s.find(value)
	return found
}

func (s *sortedSetImp[T]) String() string {
	parts := utils.Strings(s.data)
	return strings.Join(parts, `, `)
}

func (s *sortedSetImp[T]) Equals(other any) bool {
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

func (s *sortedSetImp[T]) OnChange() events.Event[collections.ChangeArgs] {
	if s.event == nil {
		s.event = event.New[collections.ChangeArgs]()
	}
	return s.event
}

func (s *sortedSetImp[T]) Get(index int) T {
	if count := len(s.data); index < 0 || index >= count {
		panic(terror.OutOfBounds(index, count))
	}
	return s.data[index]
}

func (s *sortedSetImp[T]) TryGet(index int) (T, bool) {
	if index < 0 || index >= len(s.data) {
		return utils.Zero[T](), false
	}
	return s.data[index], true
}

func (s *sortedSetImp[T]) First() T {
	if len(s.data) <= 0 {
		panic(terror.EmptyCollection(`First`))
	}
	return s.data[0]
}

func (s *sortedSetImp[T]) Last() T {
	count := len(s.data)
	if count <= 0 {
		panic(terror.EmptyCollection(`Last`))
	}
	return s.data[count-1]
}

func (s *sortedSetImp[T]) IndexOf(value T) int {
	if index, found := s.find(value); found {
		return index
	}
	return -1
}

func (s *sortedSetImp[T]) add(values []T, force bool) bool {
	added := false
	s.grow(len(s.data) + len(values))
	for _, value := range values {
		_, oneAdded := s.addOne(value, force)
		added = oneAdded || added
	}
	if added {
		s.onAdded()
	}
	return added
}

func (s *sortedSetImp[T]) Add(values ...T) bool {
	return s.add(values, false)
}

func (s *sortedSetImp[T]) Overwrite(values ...T) bool {
	return s.add(values, true)
}

func (s *sortedSetImp[T]) addFrom(e collections.Enumerator[T], force bool) bool {
	if utils.IsNil(e) {
		return false
	}
	added := false
	it := e.Iterate()
	for it.Next() {
		_, oneAdded := s.addOne(it.Current(), force)
		added = oneAdded || added
	}
	if added {
		s.onAdded()
	}
	return added
}

func (s *sortedSetImp[T]) AddFrom(e collections.Enumerator[T]) bool {
	return s.addFrom(e, false)
}

func (s *sortedSetImp[T]) OverwriteFrom(e collections.Enumerator[T]) bool {
	return s.addFrom(e, true)
}

func (s *sortedSetImp[T]) TryAdd(value T) (T, bool) {
	value, added := s.addOne(value, false)
	if added {
		s.onAdded()
	}
	return value, added
}

func (s *sortedSetImp[T]) TakeFirst() T {
	maxIndex := len(s.data) - 1
	if maxIndex < 0 {
		panic(terror.EmptyCollection(`TakeFirst`))
	}
	result := s.data[0]
	copy(s.data, s.data[1:])
	s.data[maxIndex] = utils.Zero[T]()
	s.data = s.data[:maxIndex]
	s.onRemoved()
	return result
}

func (s *sortedSetImp[T]) TakeFront(count int) collections.List[T] {
	fullCount := len(s.data)
	count = min(count, fullCount)
	if count <= 0 {
		return list.New[T]()
	}
	end := fullCount - count
	result := list.With(s.data[:count]...)
	copy(s.data, s.data[count:])
	utils.SetToZero(s.data, end, fullCount)
	s.data = s.data[:end]
	s.onRemoved()
	return result
}

func (s *sortedSetImp[T]) TakeLast() T {
	maxIndex := len(s.data) - 1
	if maxIndex < 0 {
		panic(terror.EmptyCollection(`TakeLast`))
	}
	result := s.data[maxIndex]
	s.data[maxIndex] = utils.Zero[T]()
	s.data = s.data[:maxIndex]
	s.onRemoved()
	return result
}

func (s *sortedSetImp[T]) TakeBack(count int) collections.List[T] {
	fullCount := len(s.data)
	count = min(count, fullCount)
	if count <= 0 {
		return list.New[T]()
	}
	end := fullCount - count
	result := list.With(s.data[end:]...)
	utils.SetToZero(s.data, end, fullCount)
	s.data = s.data[:end]
	s.onRemoved()
	return result
}

func (s *sortedSetImp[T]) Remove(values ...T) bool {
	removed := false
	for _, value := range values {
		if index, found := s.find(value); found {
			s.data = slices.Delete(s.data, index, index+1)
			removed = true
		}
	}
	if removed {
		s.onRemoved()
	}
	return removed
}

func (s *sortedSetImp[T]) RemoveIf(predicate collections.Predicate[T]) bool {
	if utils.IsNil(predicate) {
		return false
	}
	oldLen := len(s.data)
	s.data = slices.DeleteFunc(s.data, predicate)
	if len(s.data) != oldLen {
		s.onRemoved()
		return true
	}
	return false
}

func (s *sortedSetImp[T]) RemoveRange(index, count int) {
	if count > 0 {
		s.data = slices.Delete(s.data, index, index+count)
		s.onRemoved()
	}
}

func (s *sortedSetImp[T]) Clear() {
	if len(s.data) > 0 {
		utils.SetToZero(s.data, 0, len(s.data)-1)
		s.data = s.data[:0]
		s.onRemoved()
	}
}

func (s *sortedSetImp[T]) Clone() collections.SortedSet[T] {
	return &sortedSetImp[T]{
		data:     slices.Clone(s.data),
		comparer: s.comparer,
		event:    nil,
	}
}

func (s *sortedSetImp[T]) Readonly() collections.ReadonlySortedSet[T] {
	return readonlySortedSet.New(s)
}
