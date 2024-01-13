package list

import (
	"slices"
	"strings"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/collections/iterator"
	"github.com/Snow-Gremlin/goToolbox/collections/readonlyList"
	"github.com/Snow-Gremlin/goToolbox/internal/optional"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

type listImp[T any] struct {
	s []T
}

func (list *listImp[T]) Enumerate() collections.Enumerator[T] {
	// Since we can use the length to keep the index valid
	// changes to the list don't have stop enumerators.
	// Changes may just cause the enumeration to be unstable.
	return enumerator.New(func() collections.Iterator[T] {
		index := -1
		return iterator.New(func() (T, bool) {
			if index < len(list.s)-1 {
				index++
				return list.s[index], true
			}
			return utils.Zero[T](), false
		})
	})
}

func (list *listImp[T]) Backwards() collections.Enumerator[T] {
	// See comment in Enumerate
	return enumerator.New(func() collections.Iterator[T] {
		index := len(list.s)
		return iterator.New(func() (T, bool) {
			if index = min(index, len(list.s)); index > 0 {
				index--
				return list.s[index], true
			}
			return utils.Zero[T](), false
		})
	})
}

func (list *listImp[T]) Empty() bool {
	return len(list.s) <= 0
}

func (list *listImp[T]) Count() int {
	return len(list.s)
}

func (list *listImp[T]) Contains(value T) bool {
	return list.IndexOf(value) >= 0
}

func (list *listImp[T]) IndexOf(value T, after ...int) int {
	for i, count := optional.After(after)+1, len(list.s); i < count; i++ {
		if utils.Equal(list.s[i], value) {
			return i
		}
	}
	return -1
}

func (list *listImp[T]) First() T {
	if len(list.s) <= 0 {
		panic(terror.EmptyCollection(`First`))
	}
	return list.s[0]
}

func (list *listImp[T]) Last() T {
	count := len(list.s)
	if count <= 0 {
		panic(terror.EmptyCollection(`Last`))
	}
	return list.s[count-1]
}

func (list *listImp[T]) Get(index int) T {
	if count := len(list.s); index < 0 || index >= count {
		panic(terror.OutOfBounds(index, count))
	}
	return list.s[index]
}

func (list *listImp[T]) TryGet(index int) (T, bool) {
	if index < 0 || index >= len(list.s) {
		return utils.Zero[T](), false
	}
	return list.s[index], true
}

func (list *listImp[T]) StartsWith(other collections.ReadonlyList[T]) bool {
	return list.Enumerate().StartsWith(other.Enumerate())
}

func (list *listImp[T]) EndsWith(other collections.ReadonlyList[T]) bool {
	return list.Backwards().StartsWith(other.Backwards())
}

func (list *listImp[T]) ToSlice() []T {
	return slices.Clone(list.s)
}

func (list *listImp[T]) CopyToSlice(s []T) {
	copy(s, list.s)
}

func (list *listImp[T]) String() string {
	return strings.Join(utils.Strings(list.s), `, `)
}

func (list *listImp[T]) Equals(other any) bool {
	s, ok := other.(collections.Collection[T])
	return ok && list.Count() == s.Count() &&
		list.Enumerate().Equals(s.Enumerate())
}

func (list *listImp[T]) Prepend(values ...T) {
	list.s = slices.Insert(list.s, 0, values...)
}

func (list *listImp[T]) PrependFrom(e collections.Enumerator[T]) {
	if !utils.IsNil(e) {
		list.Prepend(e.ToSlice()...)
	}
}

func (list *listImp[T]) Append(values ...T) {
	list.s = append(list.s, values...)
}

func (list *listImp[T]) AppendFrom(e collections.Enumerator[T]) {
	if !utils.IsNil(e) {
		list.Append(e.ToSlice()...)
	}
}

func (list *listImp[T]) TakeFirst() T {
	max := len(list.s) - 1
	if max < 0 {
		panic(terror.EmptyCollection(`TakeFirst`))
	}
	result := list.s[0]
	copy(list.s, list.s[1:])
	list.s[max] = utils.Zero[T]()
	list.s = list.s[:max]
	return result
}

func (list *listImp[T]) TakeFront(count int) collections.List[T] {
	fullCount := len(list.s)
	count = min(count, fullCount)
	if count <= 0 {
		return New[T]()
	}
	end := fullCount - count
	result := With(list.s[:count]...)
	copy(list.s, list.s[count:])
	utils.SetToZero(list.s, end, fullCount)
	list.s = list.s[:end]
	return result
}

func (list *listImp[T]) TakeLast() T {
	max := len(list.s) - 1
	if max < 0 {
		panic(terror.EmptyCollection(`TakeLast`))
	}
	result := list.s[max]
	list.s[max] = utils.Zero[T]()
	list.s = list.s[:max]
	return result
}

func (list *listImp[T]) TakeBack(count int) collections.List[T] {
	fullCount := len(list.s)
	count = min(count, fullCount)
	if count <= 0 {
		return New[T]()
	}
	end := fullCount - count
	result := With(list.s[end:]...)
	utils.SetToZero(list.s, end, fullCount)
	list.s = list.s[:end]
	return result
}

func (list *listImp[T]) Insert(index int, values ...T) {
	list.s = slices.Insert(list.s, index, values...)
}

func (list *listImp[T]) InsertFrom(index int, e collections.Enumerator[T]) {
	if !utils.IsNil(e) {
		list.Insert(index, e.ToSlice()...)
	}
}

func (list *listImp[T]) Remove(index, count int) {
	s := slices.Delete(list.s, index, index+count)
	utils.SetToZero(list.s, len(s)+1, len(list.s))
	list.s = s
}

func (list *listImp[T]) RemoveIf(handle collections.Predicate[T]) bool {
	s := slices.DeleteFunc(list.s, handle)
	oldCount, newCount := len(list.s), len(s)
	if oldCount == newCount {
		return false
	}
	utils.SetToZero(list.s, newCount+1, oldCount)
	list.s = s
	return true
}

func (list *listImp[T]) Set(index int, values ...T) {
	valCount := len(values)
	if valCount <= 0 {
		return
	}
	count := len(list.s)
	if index < 0 || index > count {
		panic(terror.OutOfBounds(index, count))
	}

	if index+valCount > count {
		list.s = append(list.s[:index], values...)
	} else {
		copy(list.s[index:], values)
	}
}

func (list *listImp[T]) SetFrom(index int, e collections.Enumerator[T]) {
	if !utils.IsNil(e) {
		list.Set(index, e.ToSlice()...)
	}
}

func (list *listImp[T]) Clear() {
	utils.SetToZero(list.s, 0, len(list.s))
	list.s = list.s[:0]
}

func (list *listImp[T]) Clone() collections.List[T] {
	return &listImp[T]{
		s: slices.Clone(list.s),
	}
}

func (list *listImp[T]) Readonly() collections.ReadonlyList[T] {
	return readonlyList.New(list)
}
