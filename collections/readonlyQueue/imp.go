package readonlyQueue

import "github.com/Snow-Gremlin/goToolbox/collections"

type readonlyQueueImp[T any] struct {
	q collections.ReadonlyQueue[T]
}

func (r *readonlyQueueImp[T]) Enumerate() collections.Enumerator[T] {
	return r.q.Enumerate()
}

func (r *readonlyQueueImp[T]) Empty() bool {
	return r.q.Empty()
}

func (r *readonlyQueueImp[T]) Count() int {
	return r.q.Count()
}

func (r *readonlyQueueImp[T]) String() string {
	return r.q.String()
}

func (r *readonlyQueueImp[T]) Equals(other any) bool {
	return r.q.Equals(other)
}

func (r *readonlyQueueImp[T]) ToSlice() []T {
	return r.q.ToSlice()
}

func (r *readonlyQueueImp[T]) CopyToSlice(sc []T) {
	r.q.CopyToSlice(sc)
}

func (r *readonlyQueueImp[T]) ToList() collections.List[T] {
	return r.q.ToList()
}

func (r *readonlyQueueImp[T]) Peek() T {
	return r.q.Peek()
}

func (r *readonlyQueueImp[T]) TryPeek() (T, bool) {
	return r.q.TryPeek()
}
