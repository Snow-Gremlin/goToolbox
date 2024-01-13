package readonlySet

import "github.com/Snow-Gremlin/goToolbox/collections"

type readonlySetImp[T comparable] struct {
	s collections.ReadonlySet[T]
}

func (r *readonlySetImp[T]) Enumerate() collections.Enumerator[T] {
	return r.s.Enumerate()
}

func (r *readonlySetImp[T]) Empty() bool {
	return r.s.Empty()
}

func (r *readonlySetImp[T]) Count() int {
	return r.s.Count()
}

func (r *readonlySetImp[T]) ToSlice() []T {
	return r.s.ToSlice()
}

func (r *readonlySetImp[T]) CopyToSlice(sc []T) {
	r.s.CopyToSlice(sc)
}

func (r *readonlySetImp[T]) ToList() collections.List[T] {
	return r.s.ToList()
}

func (r *readonlySetImp[T]) Contains(key T) bool {
	return r.s.Contains(key)
}

func (r *readonlySetImp[T]) String() string {
	return r.s.String()
}

func (r *readonlySetImp[T]) Equals(other any) bool {
	return r.s.Equals(other)
}
