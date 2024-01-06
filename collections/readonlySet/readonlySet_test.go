package readonlySet

import (
	"maps"
	"slices"
	"testing"

	"goToolbox/collections"
	"goToolbox/collections/enumerator"
	"goToolbox/collections/list"
	"goToolbox/internal/simpleSet"
	"goToolbox/testers/check"
	"goToolbox/utils"
)

type pseudoSetImp struct {
	m simpleSet.Set[int]
}

func newPseudoImp(values ...int) *pseudoSetImp {
	return &pseudoSetImp{simpleSet.With(values...)}
}

func (s *pseudoSetImp) Enumerate() collections.Enumerator[int] {
	return enumerator.Enumerate(s.m.ToSlice()...)
}

func (s *pseudoSetImp) Empty() bool {
	return s.m.Count() <= 0
}

func (s *pseudoSetImp) Count() int {
	return s.m.Count()
}

func (s *pseudoSetImp) ToSlice() []int {
	return utils.SortedKeys(s.m)
}

func (s *pseudoSetImp) CopyToSlice(sc []int) {
	copy(sc, s.ToSlice())
}

func (s *pseudoSetImp) ToList() collections.List[int] {
	return list.From(s.Enumerate())
}

func (s *pseudoSetImp) Contains(key int) bool {
	return s.m.Has(key)
}

func (s *pseudoSetImp) String() string {
	return s.Enumerate().Strings().Sort().Join(`, `)
}

func (s *pseudoSetImp) Equals(other any) bool {
	s2, ok := other.(collections.ReadonlySet[int])
	if !ok || s.Count() != s2.Count() {
		return false
	}

	for _, v := range s2.ToSlice() {
		if !s2.Contains(v) {
			return false
		}
	}

	return true
}

func Test_ReadonlySet(t *testing.T) {
	s0 := newPseudoImp(1, 2, 3)
	s1 := New(s0)
	check.Length(t, 3).Assert(s1)
	check.String(t, `1, 2, 3`).Assert(s1)
	check.False(t).Assert(s1.Enumerate().Empty())
	check.False(t).Assert(s1.Empty())

	p := s1.ToSlice()
	slices.Sort(p)
	check.Equal(t, []int{1, 2, 3}).Assert(p)
	check.Length(t, 3).Assert(s1.ToList())

	p = make([]int, 5)
	s1.CopyToSlice(p)
	slices.Sort(p)
	check.Equal(t, []int{0, 0, 1, 2, 3}).Assert(p)

	check.True(t).Assert(s1.Contains(1))
	check.False(t).Assert(s1.Contains(4))

	s2 := &pseudoSetImp{
		m: maps.Clone(s0.m),
	}
	s3 := New(s2)
	check.Equal(t, s3).Assert(s1)
	check.String(t, `1, 2, 3`).Assert(s3)

	s2.m.Set(5)
	check.String(t, `1, 2, 3, 5`).Assert(s3)
	check.NotEqual(t, s3).Assert(s1)
}
