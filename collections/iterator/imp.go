package iterator

import "github.com/Snow-Gremlin/goToolbox/utils"

type iteratorImp[T any] struct {
	fetcher Fetcher[T]
	current T
}

func (it *iteratorImp[T]) Next() bool {
	if it.fetcher == nil {
		return false
	}

	if next, has := it.fetcher(); has {
		it.current = next
		return true
	}

	it.current = utils.Zero[T]()
	it.fetcher = nil
	return false
}

func (it *iteratorImp[T]) Current() T {
	return it.current
}
