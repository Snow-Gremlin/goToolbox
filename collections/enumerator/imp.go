package enumerator

import (
	"strconv"
	"strings"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/iterator"
	"github.com/Snow-Gremlin/goToolbox/collections/predicate"
	"github.com/Snow-Gremlin/goToolbox/internal/optional"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

type enumeratorImp[T any] struct {
	iterable collections.Iterable[T]
}

func (e enumeratorImp[T]) Iterate() collections.Iterator[T] {
	return e.iterable()
}

func (e enumeratorImp[T]) Where(p collections.Predicate[T]) collections.Enumerator[T] {
	return New(func() collections.Iterator[T] {
		return iterator.Where(e.Iterate(), p)
	})
}

func (e enumeratorImp[T]) WhereNot(p collections.Predicate[T]) collections.Enumerator[T] {
	return New(func() collections.Iterator[T] {
		return iterator.WhereNot(e.Iterate(), p)
	})
}

func (e enumeratorImp[T]) NotNil() collections.Enumerator[T] {
	return e.WhereNot(predicate.IsNil[T]())
}

func (e enumeratorImp[T]) NotZero() collections.Enumerator[T] {
	return e.WhereNot(predicate.IsZero[T]())
}

func (e enumeratorImp[T]) ToSlice() []T {
	return iterator.ToSlice(e.Iterate())
}

func (e enumeratorImp[T]) CopyToSlice(s []T) {
	iterator.CopyToSlice(e.Iterate(), s)
}

func (e enumeratorImp[T]) Foreach(m func(value T)) {
	iterator.Foreach(e.Iterate(), m)
}

func (e enumeratorImp[T]) DoUntilError(m func(value T) error) error {
	return iterator.DoUntilError(e.Iterate(), m)
}

func (e enumeratorImp[T]) Any(p collections.Predicate[T]) bool {
	return iterator.Any(e.Iterate(), p)
}

func (e enumeratorImp[T]) All(p collections.Predicate[T]) bool {
	return iterator.All(e.Iterate(), p)
}

func (e enumeratorImp[T]) StepsUntil(p collections.Predicate[T]) int {
	return iterator.StepsUntil(e.Iterate(), p)
}

func (e enumeratorImp[T]) Empty() bool {
	return iterator.Empty(e.Iterate())
}

func (e enumeratorImp[T]) Count() int {
	return iterator.Count(e.Iterate())
}

func (e enumeratorImp[T]) AtLeast(min int) bool {
	return iterator.AtLeast(e.Iterate(), min)
}

func (e enumeratorImp[T]) AtMost(max int) bool {
	return iterator.AtMost(e.Iterate(), max)
}

func (e enumeratorImp[T]) First() (T, bool) {
	return iterator.First(e.Iterate())
}

func (e enumeratorImp[T]) Last() (T, bool) {
	return iterator.Last(e.Iterate())
}

func (e enumeratorImp[T]) Single() (T, bool) {
	return iterator.Single(e.Iterate())
}

func (e enumeratorImp[T]) Skip(count int) collections.Enumerator[T] {
	return New(func() collections.Iterator[T] {
		return iterator.Skip(e.Iterate(), count)
	})
}

func (e enumeratorImp[T]) SkipWhile(p collections.Predicate[T]) collections.Enumerator[T] {
	return New(func() collections.Iterator[T] {
		return iterator.SkipWhile(e.Iterate(), p)
	})
}

func (e enumeratorImp[T]) Take(count int) collections.Enumerator[T] {
	return New(func() collections.Iterator[T] {
		return iterator.Take(e.Iterate(), count)
	})
}

func (e enumeratorImp[T]) TakeWhile(p collections.Predicate[T]) collections.Enumerator[T] {
	return New(func() collections.Iterator[T] {
		return iterator.TakeWhile(e.Iterate(), p)
	})
}

func (e enumeratorImp[T]) Replace(replacer collections.Selector[T, T]) collections.Enumerator[T] {
	return New(func() collections.Iterator[T] {
		return iterator.Replace(e.Iterate(), replacer)
	})
}

func (e enumeratorImp[T]) Reverse() collections.Enumerator[T] {
	return New(func() collections.Iterator[T] {
		return iterator.Reverse(e.Iterate())
	})
}

func (e enumeratorImp[T]) Strings() collections.Enumerator[string] {
	return Select[T, string](e, utils.String)
}

func (e enumeratorImp[T]) Quotes() collections.Enumerator[string] {
	return Select[T, string](e, func(value T) string {
		return strconv.Quote(utils.String(value))
	})
}

func (e enumeratorImp[T]) Trim() collections.Enumerator[string] {
	return Select(e.Strings(), strings.TrimSpace)
}

func (e enumeratorImp[T]) Join(sep string) string {
	return strings.Join(e.Strings().ToSlice(), sep)
}

func (e enumeratorImp[T]) Append(tails ...T) collections.Enumerator[T] {
	return New(func() collections.Iterator[T] {
		return iterator.Append(e.Iterate(), tails)
	})
}

func (e enumeratorImp[T]) Concat(tails ...collections.Enumerator[T]) collections.Enumerator[T] {
	return New(func() collections.Iterator[T] {
		its := make([]collections.Iterator[T], len(tails)+1)
		its[0] = e.Iterate()
		for i, t := range tails {
			its[i+1] = t.Iterate()
		}
		return iterator.Concat(its)
	})
}

func (e enumeratorImp[T]) SortInterweave(other collections.Enumerator[T], comparer ...utils.Comparer[T]) collections.Enumerator[T] {
	cmp := optional.Comparer(comparer)
	return New(func() collections.Iterator[T] {
		return iterator.SortInterweave(e.Iterate(), other.Iterate(), cmp)
	})
}

func (e enumeratorImp[T]) Sorted(comparer ...utils.Comparer[T]) bool {
	cmp := optional.Comparer(comparer)
	return iterator.Sorted(e.Iterate(), cmp)
}

func (e enumeratorImp[T]) Sort(comparer ...utils.Comparer[T]) collections.Enumerator[T] {
	cmp := optional.Comparer(comparer)
	return New(func() collections.Iterator[T] {
		return iterator.Sort(e.Iterate(), cmp)
	})
}

func (e enumeratorImp[T]) Merge(merger collections.Reducer[T, T]) T {
	return iterator.Merge(e.Iterate(), merger)
}

func (e enumeratorImp[T]) Max(comparer ...utils.Comparer[T]) T {
	cmp := optional.Comparer(comparer)
	return e.Merge(func(value, prior T) T {
		if cmp(value, prior) > 0 {
			return value
		}
		return prior
	})
}

func (e enumeratorImp[T]) Min(comparer ...utils.Comparer[T]) T {
	cmp := optional.Comparer(comparer)
	return e.Merge(func(value, prior T) T {
		if cmp(value, prior) < 0 {
			return value
		}
		return prior
	})
}

func (e enumeratorImp[T]) Buffered() collections.Enumerator[T] {
	count := 0
	var buffer []T
	loading := true
	var source collections.Iterator[T]
	return New(func() collections.Iterator[T] {
		if loading && source == nil {
			source = e.Iterate()
		}

		index := 0
		return iterator.New(func() (T, bool) {
			if index < count {
				value := buffer[index]
				index++
				return value, true
			}

			if loading {
				if source.Next() {
					value := source.Current()
					buffer = append(buffer, value)
					count = len(buffer)
					index++
					return value, true
				}

				source = nil
				loading = false
			}
			return utils.Zero[T](), false
		})
	})
}

func (e enumeratorImp[T]) StartsWith(other collections.Enumerator[T]) bool {
	return iterator.StartsWith(e.Iterate(), other.Iterate())
}

func (e enumeratorImp[T]) Equals(other any) bool {
	e2, ok := other.(collections.Enumerator[T])
	return ok && iterator.Equal(e.Iterate(), e2.Iterate())
}
