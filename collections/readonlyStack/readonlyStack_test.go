package readonlyStack

import (
	"fmt"
	"slices"
	"testing"

	"goToolbox/collections"
	"goToolbox/collections/enumerator"
	"goToolbox/collections/list"
	"goToolbox/testers/check"
	"goToolbox/utils"
)

type pseudoStackImp[T any] struct {
	s []T
}

func (s *pseudoStackImp[T]) Enumerate() collections.Enumerator[T] {
	return enumerator.Enumerate(s.s...)
}

func (s *pseudoStackImp[T]) Empty() bool {
	return len(s.s) <= 0
}

func (s *pseudoStackImp[T]) Count() int {
	return len(s.s)
}

func (s *pseudoStackImp[T]) String() string {
	return fmt.Sprint(s.s)
}

func (s *pseudoStackImp[T]) Equals(other any) bool {
	v, ok := other.(collections.Sliceable[T])
	return ok && utils.Equal(s.ToSlice(), v.ToSlice())
}

func (s *pseudoStackImp[T]) ToSlice() []T {
	return slices.Clone(s.s)
}

func (s *pseudoStackImp[T]) CopyToSlice(sc []T) {
	copy(sc, s.ToSlice())
}

func (s *pseudoStackImp[T]) ToList() collections.List[T] {
	return list.From(s.Enumerate())
}

func (s *pseudoStackImp[T]) Peek() T {
	return s.s[0]
}

func (s *pseudoStackImp[T]) TryPeek() (T, bool) {
	if s.Empty() {
		return utils.Zero[T](), false
	}
	return s.s[0], true
}

func Test_ReadonlyStack(t *testing.T) {
	s0 := &pseudoStackImp[int]{
		s: []int{1, 2, 3},
	}
	s1 := New(s0)
	s2 := &pseudoStackImp[int]{
		s: []int{1, 2, 3},
	}
	s3 := New(s2)
	check.False(t).Assert(s1.Empty())
	check.Length(t, 3).Assert(s1)
	check.Equal(t, []int{1, 2, 3}).Assert(s1.Enumerate().ToSlice())
	check.String(t, `[1 2 3]`).Assert(s1)
	check.Equal(t, s3).Assert(s1)

	sc := make([]int, 2)
	s1.CopyToSlice(sc)
	check.Equal(t, []int{1, 2}).Assert(sc)

	s0.s = append(s0.s, 34)
	check.Length(t, 4).Assert(s1)
	check.String(t, `[1 2 3 34]`).Assert(s1)
	check.Equal(t, []int{1, 2, 3, 34}).Assert(s1.ToSlice())
	check.Length(t, 4).Assert(s1.ToList())

	check.Equal(t, 1).Assert(s1.Peek())
	v, ok := s1.TryPeek()
	check.Equal(t, 1).Assert(v)
	check.True(t).Assert(ok)
	check.NotEqual(t, s3).Assert(s1)
}
