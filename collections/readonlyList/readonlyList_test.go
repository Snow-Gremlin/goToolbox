package readonlyList

import (
	"fmt"
	"slices"
	"testing"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/internal/optional"
	"github.com/Snow-Gremlin/goToolbox/testers/check"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

type pseudoList[T any] struct {
	list []T
}

func pseudoWith[T any](values ...T) *pseudoList[T] {
	return &pseudoList[T]{list: values}
}

func (p *pseudoList[T]) Enumerate() collections.Enumerator[T] {
	return enumerator.Enumerate(p.list...)
}

func (p *pseudoList[T]) Backwards() collections.Enumerator[T] {
	s2 := slices.Clone(p.list)
	slices.Reverse(s2)
	return enumerator.Enumerate(s2...)
}

func (p *pseudoList[T]) Empty() bool {
	return len(p.list) <= 0
}

func (p *pseudoList[T]) Count() int {
	return len(p.list)
}

func (p *pseudoList[T]) Contains(value T) bool {
	return p.IndexOf(value) >= 0
}

func (p *pseudoList[T]) IndexOf(value T, after ...int) int {
	for i, count := optional.After(after)+1, len(p.list); i < count; i++ {
		if utils.Equal(p.list[i], value) {
			return i
		}
	}
	return -1
}

func (p *pseudoList[T]) First() T {
	return p.list[0]
}

func (p *pseudoList[T]) Last() T {
	return p.list[len(p.list)-1]
}

func (p *pseudoList[T]) Get(index int) T {
	return p.list[index]
}

func (p *pseudoList[T]) TryGet(index int) (T, bool) {
	if index < 0 || index >= len(p.list) {
		return utils.Zero[T](), false
	}
	return p.list[index], true
}

func (list *pseudoList[T]) StartsWith(other collections.ReadonlyList[T]) bool {
	return list.Enumerate().StartsWith(other.Enumerate())
}

func (list *pseudoList[T]) EndsWith(other collections.ReadonlyList[T]) bool {
	return list.Backwards().StartsWith(other.Backwards())
}

func (p *pseudoList[T]) ToSlice() []T {
	return slices.Clone(p.list)
}

func (p *pseudoList[T]) CopyToSlice(sc []T) {
	copy(sc, p.ToSlice())
}

func (p *pseudoList[T]) String() string {
	return fmt.Sprint(p.list)
}

func (p *pseudoList[T]) Equals(other any) bool {
	s, ok := other.(collections.ReadonlyList[T])
	if !ok || len(p.list) != s.Count() {
		return false
	}

	for i, v := range p.list {
		if !utils.Equal(v, s.Get(i)) {
			return false
		}
	}

	return true
}

func Test_ReadonlyList(t *testing.T) {
	s0 := pseudoWith(1, 2, 3, 4, 5)
	s := New(s0)
	check.String(t, `[1 2 3 4 5]`).Assert(s)
	check.Length(t, 5).Assert(s)
	check.False(t).Assert(s.Empty())
	check.Equal(t, []int{1, 2, 3, 4, 5}).Assert(s.ToSlice())
	check.Equal(t, []int{1, 2, 3, 4, 5}).Assert(s.Enumerate().ToSlice())
	check.Equal(t, []int{5, 4, 3, 2, 1}).Assert(s.Backwards().ToSlice())

	p := make([]int, 3)
	s.CopyToSlice(p)
	check.Equal(t, []int{1, 2, 3}).Assert(p)

	p = make([]int, 8)
	s.CopyToSlice(p)
	check.Equal(t, []int{1, 2, 3, 4, 5, 0, 0, 0}).Assert(p)

	check.Equal(t, 3).Assert(s.IndexOf(4))
	check.Equal(t, -1).Assert(s.IndexOf(6))
	check.True(t).Assert(s.Contains(3))
	check.False(t).Assert(s.Contains(-1))
	check.Equal(t, 1).Assert(s.First())
	check.Equal(t, 5).Assert(s.Last())
	check.Equal(t, 1).Assert(s.Get(0))
	check.Equal(t, 2).Assert(s.Get(1))
	check.Equal(t, 3).Assert(s.Get(2))
	check.Equal(t, 4).Assert(s.Get(3))
	check.Equal(t, 5).Assert(s.Get(4))

	v, ok := s.TryGet(2)
	check.True(t).Assert(ok)
	check.Equal(t, 3).Assert(v)

	v, ok = s.TryGet(-1)
	check.False(t).Assert(ok)
	check.Zero(t).Assert(v)

	check.True(t).Withf(`[%s].StartsWith([1, 2, 3])`, s.String()).Assert(s.StartsWith(pseudoWith(1, 2, 3)))
	check.False(t).Withf(`[%s].StartsWith([1, 2, 4])`, s.String()).Assert(s.StartsWith(pseudoWith(1, 2, 4)))
	check.True(t).Withf(`[%s].EndsWith([3, 4, 5])`, s.String()).Assert(s.EndsWith(pseudoWith(3, 4, 5)))
	check.False(t).Withf(`[%s].EndsWith([3, 2, 5])`, s.String()).Assert(s.EndsWith(pseudoWith(3, 2, 5)))

	s2 := New(&pseudoList[int]{
		list: slices.Clone(s0.list),
	})
	check.Equal(t, s).Assert(s2)
	check.Equal(t, s2).Assert(s)
	s0.list[0] = 42
	check.NotEqual(t, s).Assert(s2)
	check.NotEqual(t, s2).Assert(s)
}
