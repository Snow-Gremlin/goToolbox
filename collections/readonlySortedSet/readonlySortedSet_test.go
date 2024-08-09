package readonlySortedSet

import (
	"slices"
	"sort"
	"testing"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/collections/list"
	"github.com/Snow-Gremlin/goToolbox/comp"
	"github.com/Snow-Gremlin/goToolbox/events"
	"github.com/Snow-Gremlin/goToolbox/events/event"
	"github.com/Snow-Gremlin/goToolbox/testers/check"
)

type pseudoSortedSetImp struct {
	data []int
	e    events.Event[collections.ChangeArgs]
}

func newPseudoImp(values ...int) *pseudoSortedSetImp {
	sort.Ints(values)
	return &pseudoSortedSetImp{
		data: values,
		e:    event.New[collections.ChangeArgs](),
	}
}

func (s *pseudoSortedSetImp) Enumerate() collections.Enumerator[int] {
	return enumerator.Enumerate(s.data...)
}

func (s *pseudoSortedSetImp) Empty() bool {
	return len(s.data) <= 0
}

func (s *pseudoSortedSetImp) Count() int {
	return len(s.data)
}

func (s *pseudoSortedSetImp) Get(index int) int {
	return s.data[index]
}

func (s *pseudoSortedSetImp) TryGet(index int) (int, bool) {
	if index < 0 || index >= len(s.data) {
		return 0, false
	}
	return s.data[index], true
}

func (s *pseudoSortedSetImp) First() int {
	return s.data[0]
}

func (s *pseudoSortedSetImp) Last() int {
	return s.data[len(s.data)-1]
}

func (s *pseudoSortedSetImp) Backwards() collections.Enumerator[int] {
	return enumerator.Enumerate(s.data...).Reverse()
}

func (s *pseudoSortedSetImp) ToSlice() []int {
	return slices.Clone(s.data)
}

func (s *pseudoSortedSetImp) CopyToSlice(sc []int) {
	copy(sc, s.ToSlice())
}

func (s *pseudoSortedSetImp) ToList() collections.List[int] {
	return list.From(s.Enumerate())
}

func (s *pseudoSortedSetImp) Contains(value int) bool {
	return s.IndexOf(value) >= 0
}

func (s *pseudoSortedSetImp) IndexOf(value int) int {
	for i, count := 0, len(s.data); i < count; i++ {
		if comp.Equal(s.data[i], value) {
			return i
		}
	}
	return -1
}

func (s *pseudoSortedSetImp) String() string {
	return s.Enumerate().Strings().Sort().Join(`, `)
}

func (s *pseudoSortedSetImp) Equals(other any) bool {
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

func (s *pseudoSortedSetImp) OnChange() events.Event[collections.ChangeArgs] {
	return s.e
}

func Test_ReadonlySet(t *testing.T) {
	s0 := newPseudoImp(1, 2, 3)
	s1 := New(s0)
	check.Length(t, 3).Assert(s1)
	check.String(t, `1, 2, 3`).Assert(s1)
	check.False(t).Assert(s1.Enumerate().Empty())
	check.False(t).Assert(s1.Backwards().Empty())
	check.False(t).Assert(s1.Empty())
	check.Equal(t, 1).Assert(s1.First())
	check.Equal(t, 3).Assert(s1.Last())
	check.Equal(t, 1).Assert(s1.IndexOf(2))
	check.Equal(t, -1).Assert(s1.IndexOf(5))
	check.Equal(t, 3).Assert(s1.Get(2))
	v, ok := s1.TryGet(3)
	check.Equal(t, 0).Assert(v)
	check.False(t).Assert(ok)

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
	check.Same(t, s0.OnChange()).Assert(s1.OnChange())

	s2 := &pseudoSortedSetImp{
		data: slices.Clone(s0.data),
		e:    event.New[collections.ChangeArgs](),
	}
	s3 := New(s2)
	check.Equal(t, s3).Assert(s1)
	check.String(t, `1, 2, 3`).Assert(s3)

	s2.data = append(s2.data, 5)
	check.String(t, `1, 2, 3, 5`).Assert(s3)
	check.NotEqual(t, s3).Assert(s1)
	check.Same(t, s2.OnChange()).Assert(s3.OnChange())
}
