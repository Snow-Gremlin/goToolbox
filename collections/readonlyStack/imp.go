package readonlyStack

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/events"
)

type readonlyStackImp[T any] struct {
	s collections.ReadonlyStack[T]
}

func (r readonlyStackImp[T]) Enumerate() collections.Enumerator[T] {
	return r.s.Enumerate()
}

func (r readonlyStackImp[T]) Empty() bool {
	return r.s.Empty()
}

func (r readonlyStackImp[T]) Count() int {
	return r.s.Count()
}

func (r readonlyStackImp[T]) String() string {
	return r.s.String()
}

func (r readonlyStackImp[T]) Equals(other any) bool {
	return r.s.Equals(other)
}

func (r readonlyStackImp[T]) ToSlice() []T {
	return r.s.ToSlice()
}

func (r readonlyStackImp[T]) CopyToSlice(s []T) {
	r.s.CopyToSlice(s)
}

func (r readonlyStackImp[T]) ToList() collections.List[T] {
	return r.s.ToList()
}

func (r readonlyStackImp[T]) Peek() T {
	return r.s.Peek()
}

func (r readonlyStackImp[T]) TryPeek() (T, bool) {
	return r.s.TryPeek()
}

func (r readonlyStackImp[T]) OnChange() events.Event[collections.ChangeArgs] {
	return r.s.OnChange()
}
