package linkedList

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/collections/iterator"
	"github.com/Snow-Gremlin/goToolbox/collections/readonlyList"
	"github.com/Snow-Gremlin/goToolbox/internal/optional"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

func newImp[T any](s ...T) *linkedListImp[T] {
	count := len(s)
	list := &linkedListImp[T]{
		count:     0,
		head:      nil,
		tail:      nil,
		enumGuard: 0,
	}
	if count <= 0 {
		return list
	}

	chunk := make([]node[T], count)
	prev := &chunk[0]
	prev.value = s[0]
	for i := 1; i < count; i++ {
		n := &chunk[i]
		n.value = s[i]
		n.prev = prev
		prev.next = n
		prev = n
	}

	list.head = &chunk[0]
	list.tail = &chunk[count-1]
	list.count = count
	return list
}

func impFrom[T any](e collections.Enumerator[T]) *linkedListImp[T] {
	list := &linkedListImp[T]{
		count:     0,
		head:      nil,
		tail:      nil,
		enumGuard: 0,
	}
	if utils.IsNil(e) {
		return list
	}

	it := e.Iterate()
	if !it.Next() {
		return list
	}

	last := newNode(it.Current(), nil)
	list.head = last
	count := 1

	for it.Next() {
		n := newNode(it.Current(), last)
		last.next = n
		last = n
		count++
	}

	list.tail = last
	list.count = count
	return list
}

type linkedListImp[T any] struct {
	count     int
	head      *node[T]
	tail      *node[T]
	enumGuard uint
}

func (list *linkedListImp[T]) overwriteWith(temp *linkedListImp[T]) {
	list.tail = temp.tail
	list.head = temp.head
	list.count = temp.count
}

func (list *linkedListImp[T]) nodeAt(index int) *node[T] {
	if index2 := list.count - index - 1; index > index2 {
		return list.tail.backwards(index2)
	}
	return list.head.forward(index)
}

func (list *linkedListImp[T]) Enumerate() collections.Enumerator[T] {
	return enumerator.New(func() collections.Iterator[T] {
		n := list.head
		guardStash := list.enumGuard
		return iterator.New(func() (T, bool) {
			if n == nil {
				return utils.Zero[T](), false
			}
			if guardStash != list.enumGuard {
				panic(terror.UnstableIteration())
			}
			value := n.value
			n = n.next
			return value, true
		})
	})
}

func (list *linkedListImp[T]) Backwards() collections.Enumerator[T] {
	return enumerator.New(func() collections.Iterator[T] {
		n := list.tail
		guardStash := list.enumGuard
		return iterator.New(func() (T, bool) {
			if n == nil {
				return utils.Zero[T](), false
			}
			if guardStash != list.enumGuard {
				panic(terror.UnstableIteration())
			}
			value := n.value
			n = n.prev
			return value, true
		})
	})
}

func (list *linkedListImp[T]) Empty() bool {
	return list.head == nil
}

func (list *linkedListImp[T]) Count() int {
	return list.count
}

func (list *linkedListImp[T]) Contains(value T) bool {
	return list.IndexOf(value) >= 0
}

func (list *linkedListImp[T]) IndexOf(value T, after ...int) int {
	start := optional.After(after)
	for it, i := list.head, 0; it != nil; it, i = it.next, i+1 {
		if i > start && utils.Equal(it.value, value) {
			return i
		}
	}
	return -1
}

func (list *linkedListImp[T]) First() T {
	if n := list.head; n != nil {
		return n.value
	}
	panic(terror.EmptyCollection(`First`))
}

func (list *linkedListImp[T]) Last() T {
	if n := list.tail; n != nil {
		return n.value
	}
	panic(terror.EmptyCollection(`Last`))
}

func (list *linkedListImp[T]) Get(index int) T {
	if index < 0 || index >= list.count {
		panic(terror.OutOfBounds(index, list.count))
	}
	return list.nodeAt(index).value
}

func (list *linkedListImp[T]) TryGet(index int) (T, bool) {
	if index < 0 || index >= list.count {
		return utils.Zero[T](), false
	}
	return list.nodeAt(index).value, true
}

func (list *linkedListImp[T]) StartsWith(other collections.ReadonlyList[T]) bool {
	return list.Enumerate().StartsWith(other.Enumerate())
}

func (list *linkedListImp[T]) EndsWith(other collections.ReadonlyList[T]) bool {
	return list.Backwards().StartsWith(other.Backwards())
}

func (list *linkedListImp[T]) ToSlice() []T {
	return list.Enumerate().ToSlice()
}

func (list *linkedListImp[T]) CopyToSlice(s []T) {
	list.Enumerate().CopyToSlice(s)
}

func (list *linkedListImp[T]) String() string {
	return list.Enumerate().Join(`, `)
}

func (list *linkedListImp[T]) Equals(other any) bool {
	s, ok := other.(collections.Collection[T])
	return ok && list.count == s.Count() &&
		list.Enumerate().Equals(s.Enumerate())
}

func (list *linkedListImp[T]) prependOther(temp *linkedListImp[T]) {
	if temp.Empty() {
		return
	}
	if list.Empty() {
		list.overwriteWith(temp)
		return
	}
	temp.tail.next = list.head
	list.head.prev = temp.tail
	list.head = temp.head
	list.count += temp.count
}

func (list *linkedListImp[T]) Prepend(values ...T) {
	list.prependOther(newImp(values...))
}

func (list *linkedListImp[T]) PrependFrom(e collections.Enumerator[T]) {
	list.prependOther(impFrom(e))
}

func (list *linkedListImp[T]) appendOther(temp *linkedListImp[T]) {
	if temp.Empty() {
		return
	}
	if list.Empty() {
		list.overwriteWith(temp)
		return
	}
	temp.head.prev = list.tail
	list.tail.next = temp.head
	list.tail = temp.tail
	list.count += temp.count
}

func (list *linkedListImp[T]) Append(values ...T) {
	list.appendOther(newImp(values...))
}

func (list *linkedListImp[T]) AppendFrom(e collections.Enumerator[T]) {
	list.appendOther(impFrom(e))
}

func (list *linkedListImp[T]) TakeFirst() T {
	if list.head == nil {
		panic(terror.EmptyCollection(`TakeFirst`))
	}
	value := list.head.value
	list.head = list.head.next
	if list.head != nil {
		list.head.prev = nil
	} else {
		list.tail = nil
	}
	list.count--
	list.enumGuard++
	return value
}

func (list *linkedListImp[T]) TakeLast() T {
	if list.tail == nil {
		panic(terror.EmptyCollection(`TakeLast`))
	}
	value := list.tail.value
	list.tail = list.tail.prev
	if list.tail != nil {
		list.tail.next = nil
	} else {
		list.head = nil
	}
	list.count--
	list.enumGuard++
	return value
}

func (list *linkedListImp[T]) TakeFront(count int) collections.List[T] {
	count = min(count, list.count)
	if count <= 0 {
		return New[T]()
	}
	split := list.nodeAt(count - 1)
	result := &linkedListImp[T]{
		count:     count,
		head:      list.head,
		tail:      split,
		enumGuard: 0,
	}
	list.head = split.next
	split.next = nil
	if list.head == nil {
		list.tail = nil
	} else {
		list.head.prev = nil
	}
	list.count -= count
	list.enumGuard++
	return result
}

func (list *linkedListImp[T]) TakeBack(count int) collections.List[T] {
	count = min(count, list.count)
	if count <= 0 {
		return New[T]()
	}
	split := list.nodeAt(list.count - count)
	result := &linkedListImp[T]{
		count:     count,
		head:      split,
		tail:      list.tail,
		enumGuard: 0,
	}
	list.tail = split.prev
	split.prev = nil
	if list.tail == nil {
		list.head = nil
	} else {
		list.tail.next = nil
	}
	list.count -= count
	list.enumGuard++
	return result
}

func (list *linkedListImp[T]) insertOther(index int, temp *linkedListImp[T]) {
	if index < 0 || index > list.count {
		panic(terror.OutOfBounds(index, list.count))
	}
	if index == 0 {
		list.prependOther(temp)
		return
	}
	if index == list.count {
		list.appendOther(temp)
		return
	}

	split := list.nodeAt(index - 1)
	temp.tail.next = split.next
	temp.head.prev = split
	if split.next != nil {
		split.next.prev = temp.tail
	}
	split.next = temp.head
	list.count += temp.count
}

func (list *linkedListImp[T]) Insert(index int, values ...T) {
	list.insertOther(index, newImp(values...))
}

func (list *linkedListImp[T]) InsertFrom(index int, e collections.Enumerator[T]) {
	list.insertOther(index, impFrom(e))
}

func (list *linkedListImp[T]) Remove(index, count int) {
	if count <= 0 {
		return
	}
	stopIndex := index + count - 1
	if index < 0 || stopIndex >= list.count {
		panic(terror.OutOfBounds(stopIndex, list.count))
	}

	start := list.nodeAt(index)
	stop := start.forward(count)
	if start.prev != nil {
		start.prev.next = stop
	} else {
		list.head = stop
	}
	if stop == nil {
		list.tail = start.prev
	} else {
		stop.prev = start.prev
	}
	list.count -= count
	list.enumGuard++
}

func (list *linkedListImp[T]) RemoveIf(handle collections.Predicate[T]) bool {
	var prior *node[T]
	subCount := 0
	changed := false
	for it := list.head; it != nil; it = it.next {
		if handle(it.value) {
			subCount++
			continue
		}

		if subCount > 0 {
			it.prev = prior
			if prior != nil {
				prior.next = it
			} else {
				list.head = it
			}
			list.count -= subCount
			changed = true
		}
		subCount = 0
		prior = it
	}

	if subCount > 0 {
		list.tail = prior
		if prior != nil {
			prior.next = nil
		} else {
			list.head = nil
		}
		list.count -= subCount
		changed = true
	}

	if changed {
		list.enumGuard++
	}
	return changed
}

func (list *linkedListImp[T]) Set(index int, values ...T) {
	list.SetFrom(index, enumerator.Enumerate(values...))
}

func (list *linkedListImp[T]) SetFrom(index int, e collections.Enumerator[T]) {
	if utils.IsNil(e) {
		return
	}
	if index < 0 || index > list.count {
		panic(terror.OutOfBounds(index, list.count))
	}

	n := list.nodeAt(index)
	it := e.Iterate()
	for n != nil {
		if !it.Next() {
			return
		}
		n.value = it.Current()
		n = n.next
	}

	prev := list.tail
	count := 0
	for it.Next() {
		n = newNode(it.Current(), prev)
		prev.next = n
		prev = n
		count++
	}
	list.tail = n
	list.count += count
}

func (list *linkedListImp[T]) Clear() {
	list.head = nil
	list.tail = nil
	list.count = 0
	list.enumGuard++
}

func (list *linkedListImp[T]) Clone() collections.List[T] {
	return From(list.Enumerate())
}

func (list *linkedListImp[T]) Readonly() collections.ReadonlyList[T] {
	return readonlyList.New(list)
}
