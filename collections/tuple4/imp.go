package tuple4

import (
	"goToolbox/collections"
	"goToolbox/terrors/terror"
	"goToolbox/utils"
)

type tuple4Imp[T1, T2, T3, T4 any] struct {
	value1 T1
	value2 T2
	value3 T3
	value4 T4
}

func (t tuple4Imp[T1, T2, T3, T4]) Count() int { return 4 }
func (t tuple4Imp[T1, T2, T3, T4]) Value1() T1 { return t.value1 }
func (t tuple4Imp[T1, T2, T3, T4]) Value2() T2 { return t.value2 }
func (t tuple4Imp[T1, T2, T3, T4]) Value3() T3 { return t.value3 }
func (t tuple4Imp[T1, T2, T3, T4]) Value4() T4 { return t.value4 }

func (t tuple4Imp[T1, T2, T3, T4]) Values() (T1, T2, T3, T4) {
	return t.value1, t.value2, t.value3, t.value4
}

func (t tuple4Imp[T1, T2, T3, T4]) ToSlice() []any {
	return []any{t.value1, t.value2, t.value3, t.value4}
}

func (t tuple4Imp[T1, T2, T3, T4]) CopyToSlice(s []any) {
	room := len(s)
	if room >= 1 {
		s[0] = t.value1
		if room >= 2 {
			s[1] = t.value2
			if room >= 3 {
				s[2] = t.value3
				if room >= 4 {
					s[3] = t.value4
				}
			}
		}
	}
}

func (t tuple4Imp[T1, T2, T3, T4]) Get(index int) any {
	if v, ok := t.TryGet(index); ok {
		return v
	}
	panic(terror.OutOfBounds(index, t.Count()))
}

func (t tuple4Imp[T1, T2, T3, T4]) TryGet(index int) (any, bool) {
	switch index {
	case 0:
		return t.value1, true
	case 1:
		return t.value2, true
	case 2:
		return t.value3, true
	case 3:
		return t.value4, true
	default:
		return nil, false
	}
}

func (t tuple4Imp[T1, T2, T3, T4]) String() string {
	return `[` + utils.String(t.value1) +
		`, ` + utils.String(t.value2) +
		`, ` + utils.String(t.value3) +
		`, ` + utils.String(t.value4) + `]`
}

func (t tuple4Imp[T1, T2, T3, T4]) Equals(other any) bool {
	t2, ok := other.(collections.Tuple)
	return ok && t.Count() == t2.Count() &&
		utils.Equal[any](t.value1, t2.Get(0)) &&
		utils.Equal[any](t.value2, t2.Get(1)) &&
		utils.Equal[any](t.value3, t2.Get(2)) &&
		utils.Equal[any](t.value4, t2.Get(3))
}
