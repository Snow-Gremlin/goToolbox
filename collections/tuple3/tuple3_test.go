package tuple3

import (
	"testing"

	"github.com/Snow-Gremlin/goToolbox/comp"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

func Test_Tuple3(t *testing.T) {
	t3 := New[int, string, bool](42, `Answer`, true)

	checkEqual(t, 42, t3.Value1())
	checkEqual(t, `Answer`, t3.Value2())
	checkEqual(t, true, t3.Value3())

	checkEqual(t, 42, t3.Get(0))
	checkEqual(t, `Answer`, t3.Get(1))
	checkEqual(t, true, t3.Get(2))

	checkPanic(t, `index out of bounds {count: 3, index: -1}`, func() { t3.Get(-1) })
	checkPanic(t, `index out of bounds {count: 3, index: 3}`, func() { t3.Get(3) })

	checkEqual(t, 3, t3.Count())
	checkEqual(t, `[42, Answer, true]`, t3.String())
	checkEqual(t, []any{42, `Answer`, true}, t3.ToSlice())

	v1, v2, v3 := t3.Values()
	checkEqual(t, 42, v1)
	checkEqual(t, `Answer`, v2)
	checkEqual(t, true, v3)

	checkEqual(t, t3, New[int, string, bool](42, `Answer`, true))
	checkNotEqual(t, t3, New[float64, string, bool](42.0, `Answer`, true))
	checkNotEqual(t, t3, New[int, string, bool](42, `Anther`, true))
	checkNotEqual(t, t3, []any{42, `Answer`, true})
	checkNotEqual(t, t3, 42)
	checkNotEqual(t, 42, t3)

	a := make([]any, 0)
	t3.CopyToSlice(a)
	checkEqual(t, []any{}, a)
	a = make([]any, 1)
	t3.CopyToSlice(a)
	checkEqual(t, []any{42}, a)
	a = make([]any, 2)
	t3.CopyToSlice(a)
	checkEqual(t, []any{42, `Answer`}, a)
	a = make([]any, 3)
	t3.CopyToSlice(a)
	checkEqual(t, []any{42, `Answer`, true}, a)
	a = make([]any, 4)
	t3.CopyToSlice(a)
	checkEqual(t, []any{42, `Answer`, true, nil}, a)
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
