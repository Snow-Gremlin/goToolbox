package sortedSet

import (
	"slices"
	"strings"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/changeArgs"
	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/collections/iterator"
	"github.com/Snow-Gremlin/goToolbox/collections/list"
	"github.com/Snow-Gremlin/goToolbox/collections/readonlySet"
	"github.com/Snow-Gremlin/goToolbox/comp"
	"github.com/Snow-Gremlin/goToolbox/events"
	"github.com/Snow-Gremlin/goToolbox/events/event"
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

func (s *sortedSetImp[T]) addOne(value T) bool {
	if index, found := s.find(value); !found {
		s.data = slices.Insert(s.data, index, value)
		return true
	}
	return false
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

func (s *sortedSetImp[T]) Add(values ...T) bool {
	added := false
	s.grow(len(s.data) + len(values))
	// TODO: Could improve by presorting values and zipping values together.
	for _, value := range values {
		added = s.addOne(value) || added
	}
	if added {
		s.onAdded()
	}
	return added
}

func (s *sortedSetImp[T]) AddFrom(e collections.Enumerator[T]) bool {
	if utils.IsNil(e) {
		return false
	}
	added := false
	it := e.Iterate()
	for it.Next() {
		added = s.addOne(it.Current()) || added
	}
	if added {
		s.onAdded()
	}
	return added
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

func (s *sortedSetImp[T]) Clear() {
	if len(s.data) > 0 {
		utils.SetToZero(s.data, 0, len(s.data)-1)
		s.data = s.data[:0]
		s.onRemoved()
	}
}

func (s *sortedSetImp[T]) Clone() collections.Set[T] {
	return &sortedSetImp[T]{
		data:     slices.Clone(s.data),
		comparer: s.comparer,
		event:    nil,
	}
}

func (s *sortedSetImp[T]) Readonly() collections.ReadonlySet[T] {
	return readonlySet.New(s)
}
