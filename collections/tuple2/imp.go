package tuple2

import (
	"goToolbox/collections"
	"goToolbox/terrors/terror"
	"goToolbox/utils"
)

type tuple2Imp[T1, T2 any] struct {
	value1 T1
	value2 T2
}

func (t tuple2Imp[T1, T2]) Count() int       { return 2 }
func (t tuple2Imp[T1, T2]) Value1() T1       { return t.value1 }
func (t tuple2Imp[T1, T2]) Value2() T2       { return t.value2 }
func (t tuple2Imp[T1, T2]) Values() (T1, T2) { return t.value1, t.value2 }
func (t tuple2Imp[T1, T2]) ToSlice() []any   { return []any{t.value1, t.value2} }

func (t tuple2Imp[T1, T2]) CopyToSlice(s []any) {
	room := len(s)
	if room >= 1 {
		s[0] = t.value1
		if room >= 2 {
			s[1] = t.value2
		}
	}
}

func (t tuple2Imp[T1, T2]) Get(index int) any {
	if v, ok := t.TryGet(index); ok {
		return v
	}
	panic(terror.OutOfBounds(index, t.Count()))
}

func (t tuple2Imp[T1, T2]) TryGet(index int) (any, bool) {
	switch index {
	case 0:
		return t.value1, true
	case 1:
		return t.value2, true
	default:
		return nil, false
	}
}

func (t tuple2Imp[T1, T2]) String() string {
	return `[` + utils.String(t.value1) +
		`, ` + utils.String(t.value2) + `]`
}

func (t tuple2Imp[T1, T2]) Equals(other any) bool {
	t2, ok := other.(collections.Tuple)
	return ok && t.Count() == t2.Count() &&
		utils.Equal[any](t.value1, t2.Get(0)) &&
		utils.Equal[any](t.value2, t2.Get(1))
}
