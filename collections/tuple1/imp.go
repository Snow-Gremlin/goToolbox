package tuple1

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/comp"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

type tuple1Imp[T1 any] struct {
	value1 T1
}

func (t tuple1Imp[T1]) Count() int     { return 1 }
func (t tuple1Imp[T1]) Value1() T1     { return t.value1 }
func (t tuple1Imp[T1]) ToSlice() []any { return []any{t.value1} }

func (t tuple1Imp[T1]) CopyToSlice(s []any) {
	if len(s) >= 1 {
		s[0] = t.value1
	}
}

func (t tuple1Imp[T1]) Get(index int) any {
	if v, ok := t.TryGet(index); ok {
		return v
	}
	panic(terror.OutOfBounds(index, t.Count()))
}

func (t tuple1Imp[T1]) TryGet(index int) (any, bool) {
	switch index {
	case 0:
		return t.value1, true
	default:
		return nil, false
	}
}

func (t tuple1Imp[T1]) String() string { return `[` + utils.String(t.value1) + `]` }

func (t tuple1Imp[T1]) Equals(other any) bool {
	t2, ok := other.(collections.Tuple)
	return ok && t.Count() == t2.Count() &&
		comp.Equal[any](t.value1, t2.Get(0))
}
