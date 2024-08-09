package readonlySortedSet

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/events"
)

type readonlySortedSetImp[T any] struct {
	s collections.ReadonlySortedSet[T]
}

func (r readonlySortedSetImp[T]) Enumerate() collections.Enumerator[T] {
	return r.s.Enumerate()
}

func (r readonlySortedSetImp[T]) Empty() bool {
	return r.s.Empty()
}

func (r readonlySortedSetImp[T]) Count() int {
	return r.s.Count()
}

func (r readonlySortedSetImp[T]) Get(index int) T {
	return r.s.Get(index)
}

func (r readonlySortedSetImp[T]) TryGet(index int) (T, bool) {
	return r.s.TryGet(index)
}

func (r readonlySortedSetImp[T]) First() T {
	return r.s.First()
}

func (r readonlySortedSetImp[T]) Last() T {
	return r.s.Last()
}

func (r readonlySortedSetImp[T]) Backwards() collections.Enumerator[T] {
	return r.s.Backwards()
}

func (r readonlySortedSetImp[T]) IndexOf(value T) int {
	return r.s.IndexOf(value)
}

func (r readonlySortedSetImp[T]) ToSlice() []T {
	return r.s.ToSlice()
}

func (r readonlySortedSetImp[T]) CopyToSlice(sc []T) {
	r.s.CopyToSlice(sc)
}

func (r readonlySortedSetImp[T]) ToList() collections.List[T] {
	return r.s.ToList()
}

func (r readonlySortedSetImp[T]) Contains(key T) bool {
	return r.s.Contains(key)
}

func (r readonlySortedSetImp[T]) String() string {
	return r.s.String()
}

func (r readonlySortedSetImp[T]) Equals(other any) bool {
	return r.s.Equals(other)
}

func (r readonlySortedSetImp[T]) OnChange() events.Event[collections.ChangeArgs] {
	return r.s.OnChange()
}
