package readonlyVariantList

import (
	"goToolbox/collections"
	"goToolbox/collections/enumerator"
	"goToolbox/collections/iterator"
	"goToolbox/internal/optional"
	"goToolbox/terrors/terror"
	"goToolbox/utils"
)

type impReadonlyVariantList[T any] struct {
	countHandle func() int
	getHandle   func(int) T
}

func (list *impReadonlyVariantList[T]) liteGet(index int) T {
	if list == nil || list.getHandle == nil {
		return utils.Zero[T]()
	}
	return list.getHandle(index)
}

func (list *impReadonlyVariantList[T]) Enumerate() collections.Enumerator[T] {
	// Since we can use the length to keep the index valid
	// changes to the list don't have stop enumerators.
	// Changes may just cause the enumeration to be unstable.
	return enumerator.New(func() collections.Iterator[T] {
		index := -1
		return iterator.New(func() (T, bool) {
			if index++; index < list.Count() {
				return list.liteGet(index), true
			}
			return utils.Zero[T](), false
		})
	})
}

func (list *impReadonlyVariantList[T]) Backwards() collections.Enumerator[T] {
	// See comment in Enumerate
	return enumerator.New(func() collections.Iterator[T] {
		index := list.Count()
		return iterator.New(func() (T, bool) {
			if index = min(index, list.Count()) - 1; index >= 0 {
				return list.liteGet(index), true
			}
			return utils.Zero[T](), false
		})
	})
}

func (list *impReadonlyVariantList[T]) Empty() bool {
	return list.Count() <= 0
}

func (list *impReadonlyVariantList[T]) Count() int {
	if list == nil || list.countHandle == nil {
		return 0
	}
	return list.countHandle()
}

func (list *impReadonlyVariantList[T]) Contains(value T) bool {
	return list.IndexOf(value) >= 0
}

func (list *impReadonlyVariantList[T]) IndexOf(value T, after ...int) int {
	for i, count := optional.After(after)+1, list.Count(); i < count; i++ {
		if utils.Equal(list.liteGet(i), value) {
			return i
		}
	}
	return -1
}

func (list *impReadonlyVariantList[T]) First() T {
	if list.Count() <= 0 {
		panic(terror.EmptyCollection(`First`))
	}
	return list.liteGet(0)
}

func (list *impReadonlyVariantList[T]) Last() T {
	count := list.Count()
	if count <= 0 {
		panic(terror.EmptyCollection(`Last`))
	}
	return list.liteGet(count - 1)
}

func (list *impReadonlyVariantList[T]) Get(index int) T {
	if count := list.Count(); index < 0 || index >= count {
		panic(terror.OutOfBounds(index, count))
	}
	return list.liteGet(index)
}

func (list *impReadonlyVariantList[T]) TryGet(index int) (value T, ok bool) {
	if index < 0 || index >= list.Count() {
		return utils.Zero[T](), false
	}
	return list.liteGet(index), true
}

func (list *impReadonlyVariantList[T]) StartsWith(other collections.ReadonlyList[T]) bool {
	return list.Enumerate().StartsWith(other.Enumerate())
}

func (list *impReadonlyVariantList[T]) EndsWith(other collections.ReadonlyList[T]) bool {
	return list.Backwards().StartsWith(other.Backwards())
}

func (list *impReadonlyVariantList[T]) ToSlice() []T {
	return list.Enumerate().ToSlice()
}

func (list *impReadonlyVariantList[T]) CopyToSlice(sc []T) {
	list.Enumerate().CopyToSlice(sc)
}

func (list *impReadonlyVariantList[T]) String() string {
	return list.Enumerate().Strings().Join(`, `)
}

func (list *impReadonlyVariantList[T]) Equals(other any) bool {
	s, ok := other.(collections.Collection[T])
	return ok && list.Count() == s.Count() &&
		list.Enumerate().Equals(s.Enumerate())
}
