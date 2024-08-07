package tuple4

import (
	"testing"

	"github.com/Snow-Gremlin/goToolbox/comp"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

func Test_Tuple4(t *testing.T) {
	t4 := New[int, string, bool, float64](42, `Answer`, true, 3.14)

	checkEqual(t, 42, t4.Value1())
	checkEqual(t, `Answer`, t4.Value2())
	checkEqual(t, true, t4.Value3())
	checkEqual(t, 3.14, t4.Value4())

	checkEqual(t, 42, t4.Get(0))
	checkEqual(t, `Answer`, t4.Get(1))
	checkEqual(t, true, t4.Get(2))
	checkEqual(t, 3.14, t4.Get(3))

	checkPanic(t, `index out of bounds {count: 4, index: -1}`, func() { t4.Get(-1) })
	checkPanic(t, `index out of bounds {count: 4, index: 4}`, func() { t4.Get(4) })

	checkEqual(t, 4, t4.Count())
	checkEqual(t, `[42, Answer, true, 3.14]`, t4.String())
	checkEqual(t, []any{42, `Answer`, true, 3.14}, t4.ToSlice())

	v1, v2, v3, v4 := t4.Values()
	checkEqual(t, 42, v1)
	checkEqual(t, `Answer`, v2)
	checkEqual(t, true, v3)
	checkEqual(t, 3.14, v4)

	checkEqual(t, t4, New[int, string, bool, float64](42, `Answer`, true, 3.14))
	checkNotEqual(t, t4, New[float64, string, bool, float64](42.0, `Answer`, true, 3.14))
	checkNotEqual(t, t4, New[int, string, bool, float64](42, `Anther`, true, 3.14))
	checkNotEqual(t, t4, []any{42, `Answer`, true, 3.14})
	checkNotEqual(t, t4, 42)
	checkNotEqual(t, 42, t4)

	a := make([]any, 0)
	t4.CopyToSlice(a)
	checkEqual(t, []any{}, a)
	a = make([]any, 1)
	t4.CopyToSlice(a)
	checkEqual(t, []any{42}, a)
	a = make([]any, 2)
	t4.CopyToSlice(a)
	checkEqual(t, []any{42, `Answer`}, a)
	a = make([]any, 3)
	t4.CopyToSlice(a)
	checkEqual(t, []any{42, `Answer`, true}, a)
	a = make([]any, 4)
	t4.CopyToSlice(a)
	checkEqual(t, []any{42, `Answer`, true, 3.14}, a)
	a = make([]any, 5)
	t4.CopyToSlice(a)
	checkEqual(t, []any{42, `Answer`, true, 3.14, nil}, a)
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
