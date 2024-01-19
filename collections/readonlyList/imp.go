package readonlyList

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/events"
)

type readonlyListImp[T any] struct {
	list collections.ReadonlyList[T]
}

func (r readonlyListImp[T]) Enumerate() collections.Enumerator[T] {
	return r.list.Enumerate()
}

func (r readonlyListImp[T]) Backwards() collections.Enumerator[T] {
	return r.list.Backwards()
}

func (r readonlyListImp[T]) Empty() bool {
	return r.list.Empty()
}

func (r readonlyListImp[T]) Count() int {
	return r.list.Count()
}

func (r readonlyListImp[T]) Contains(value T) bool {
	return r.list.Contains(value)
}

func (r readonlyListImp[T]) IndexOf(value T, after ...int) int {
	return r.list.IndexOf(value, after...)
}

func (r readonlyListImp[T]) First() T {
	return r.list.First()
}

func (r readonlyListImp[T]) Last() T {
	return r.list.Last()
}

func (r readonlyListImp[T]) Get(index int) T {
	return r.list.Get(index)
}

func (r readonlyListImp[T]) TryGet(index int) (T, bool) {
	return r.list.TryGet(index)
}

func (r readonlyListImp[T]) StartsWith(other collections.ReadonlyList[T]) bool {
	return r.list.StartsWith(other)
}

func (r readonlyListImp[T]) EndsWith(other collections.ReadonlyList[T]) bool {
	return r.list.EndsWith(other)
}

func (r readonlyListImp[T]) ToSlice() []T {
	return r.list.ToSlice()
}

func (r readonlyListImp[T]) CopyToSlice(sc []T) {
	r.list.CopyToSlice(sc)
}

func (r readonlyListImp[T]) String() string {
	return r.list.String()
}

func (r readonlyListImp[T]) Equals(other any) bool {
	return r.list.Equals(other)
}

func (r readonlyListImp[T]) OnChange() events.Event[collections.ChangeArgs] {
	return r.list.OnChange()
}
