package tuple1

import (
	"testing"

	"github.com/Snow-Gremlin/goToolbox/comp"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

func Test_Tuple1(t *testing.T) {
	t1 := New[int](42)

	checkEqual(t, 42, t1.Value1())
	checkEqual(t, 42, t1.Get(0))

	checkPanic(t, `index out of bounds {count: 1, index: -1}`, func() { t1.Get(-1) })
	checkPanic(t, `index out of bounds {count: 1, index: 1}`, func() { t1.Get(1) })

	checkEqual(t, 1, t1.Count())
	checkEqual(t, `[42]`, t1.String())
	checkEqual(t, []any{42}, t1.ToSlice())

	checkEqual(t, t1, New[int](42))
	checkNotEqual(t, t1, New[float64](42.0))
	checkNotEqual(t, t1, New[int](43))
	checkNotEqual(t, t1, []any{42})
	checkNotEqual(t, t1, 42)
	checkNotEqual(t, 42, t1)

	a := make([]any, 0)
	t1.CopyToSlice(a)
	checkEqual(t, []any{}, a)
	a = make([]any, 1)
	t1.CopyToSlice(a)
	checkEqual(t, []any{42}, a)
	a = make([]any, 2)
	t1.CopyToSlice(a)
	checkEqual(t, []any{42, nil}, a)
}

func checkEqual(t *testing.T, exp, actual any) {
	t.Helper()
	if !comp.Equal(exp, actual) {
		t.Errorf("\n"+
			"Expected value didn't match the actual value:\n"+
			"Actual:   %v (%T)\n"+
			"Expected: %v (%T)", actual, actual, exp, exp)
	}
}

func checkNotEqual(t *testing.T, exp, actual any) {
	t.Helper()
	if comp.Equal(exp, actual) {
		t.Errorf("\n"+
			"Expected value shouldn't have matched the actual value:\n"+
			"Actual:   %v (%T)\n"+
			"Expected: %v (%T)", actual, actual, exp, exp)
	}
}

func checkPanic(t *testing.T, exp string, handle func()) {
	t.Helper()
	actual := func() (r string) {
		defer func() { r = utils.String(recover()) }()
		handle()
		t.Error(`expected a panic`)
		return ``
	}()
	checkEqual(t, exp, actual)
}
