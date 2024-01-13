package terror

import (
	"github.com/Snow-Gremlin/goToolbox/internal/liteUtils"
	"github.com/Snow-Gremlin/goToolbox/terrors"
)

// Iterator is a collection.Iterator.
type Iterator[T any] interface {
	Next() bool
	Current() T
}

// Walk returns an iterator function for walking through all the errors
// in the error tree in depth first order.
func Walk(err error) Iterator[error] {
	return &walkIteratorImp{
		head:    pushWalk(nil, err),
		current: nil,
	}
}

// Unwrap will return all the errors wrapped in the given error.
func Unwrap(err error) []error {
	if e, ok := err.(terrors.MonoWrap); ok {
		return []error{e.Unwrap()}
	}
	if e, ok := err.(terrors.MultiWrap); ok {
		return e.Unwrap()
	}
	return nil
}

type walkIteratorImp struct {
	head    *walkStackNode
	current error
}

func (it *walkIteratorImp) Next() bool {
	if it.head == nil {
		it.current = nil
		return false
	}

	it.head, it.current = popWalk(it.head)
	return true
}

func (it *walkIteratorImp) Current() error {
	return it.current
}

type walkStackNode struct {
	err  error
	next *walkStackNode
}

func pushWalk(next *walkStackNode, err error) *walkStackNode {
	if liteUtils.IsZero(err) {
		return next
	}
	return &walkStackNode{
		err:  err,
		next: next,
	}
}

func popWalk(head *walkStackNode) (*walkStackNode, error) {
	err := head.err
	head = head.next
	wrapped := Unwrap(err)
	for i := len(wrapped) - 1; i >= 0; i-- {
		head = pushWalk(head, wrapped[i])
	}
	return head, err
}
