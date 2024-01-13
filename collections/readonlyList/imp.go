package readonlyList

import "github.com/Snow-Gremlin/goToolbox/collections"

type readonlyListImp[T any] struct {
	list collections.ReadonlyList[T]
}

func (list *readonlyListImp[T]) Enumerate() collections.Enumerator[T] {
	return list.list.Enumerate()
}

func (list *readonlyListImp[T]) Backwards() collections.Enumerator[T] {
	return list.list.Backwards()
}

func (list *readonlyListImp[T]) Empty() bool {
	return list.list.Empty()
}

func (list *readonlyListImp[T]) Count() int {
	return list.list.Count()
}

func (list *readonlyListImp[T]) Contains(value T) bool {
	return list.list.Contains(value)
}

func (list *readonlyListImp[T]) IndexOf(value T, after ...int) int {
	return list.list.IndexOf(value, after...)
}

func (list *readonlyListImp[T]) First() T {
	return list.list.First()
}

func (list *readonlyListImp[T]) Last() T {
	return list.list.Last()
}

func (list *readonlyListImp[T]) Get(index int) T {
	return list.list.Get(index)
}

func (list *readonlyListImp[T]) TryGet(index int) (T, bool) {
	return list.list.TryGet(index)
}

func (list *readonlyListImp[T]) StartsWith(other collections.ReadonlyList[T]) bool {
	return list.list.StartsWith(other)
}

func (list *readonlyListImp[T]) EndsWith(other collections.ReadonlyList[T]) bool {
	return list.list.EndsWith(other)
}

func (list *readonlyListImp[T]) ToSlice() []T {
	return list.list.ToSlice()
}

func (list *readonlyListImp[T]) CopyToSlice(sc []T) {
	list.list.CopyToSlice(sc)
}

func (list *readonlyListImp[T]) String() string {
	return list.list.String()
}

func (list *readonlyListImp[T]) Equals(other any) bool {
	return list.list.Equals(other)
}
