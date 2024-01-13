package linkedList

func newNode[T any](value T, prev *node[T]) *node[T] {
	return &node[T]{
		value: value,
		prev:  prev,
		next:  nil,
	}
}

type node[T any] struct {
	value T
	prev  *node[T]
	next  *node[T]
}

func (it *node[T]) forward(count int) *node[T] {
	for count > 0 && it != nil {
		count--
		it = it.next
	}
	return it
}

func (it *node[T]) backwards(count int) *node[T] {
	for count > 0 && it != nil {
		count--
		it = it.prev
	}
	return it
}
