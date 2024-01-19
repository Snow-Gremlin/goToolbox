package readonlyList

import "github.com/Snow-Gremlin/goToolbox/collections"

// New wraps another list in a readonly shell.
func New[T any](list collections.ReadonlyList[T]) collections.ReadonlyList[T] {
	return readonlyListImp[T]{list: list}
}
