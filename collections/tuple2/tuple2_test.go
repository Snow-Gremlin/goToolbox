package tuple2

import (
	"testing"

	"goToolbox/utils"
)

func Test_Tuple2(t *testing.T) {
	t2 := New[int, string](42, `Answer`)

	checkEqual(t, 42, t2.Value1())
	checkEqual(t, `Answer`, t2.Value2())

	checkEqual(t, 42, t2.Get(0))
	checkEqual(t, `Answer`, t2.Get(1))

	checkPanic(t, `index out of bounds {count: 2, index: -1}`, func() { t2.Get(-1) })
	checkPanic(t, `index out of bounds {count: 2, index: 2}`, func() { t2.Get(2) })

	checkEqual(t, 2, t2.Count())
	checkEqual(t, `[42, Answer]`, t2.String())
	checkEqual(t, []any{42, `Answer`}, t2.ToSlice())

	v1, v2 := t2.Values()
	checkEqual(t, 42, v1)
	checkEqual(t, `Answer`, v2)

	checkEqual(t, t2, New[int, string](42, `Answer`))
	checkNotEqual(t, t2, New[float64, string](42.0, `Answer`))
	checkNotEqual(t, t2, New[int, string](42, `Anther`))
	checkNotEqual(t, t2, []any{42, `Answer`})
	checkNotEqual(t, t2, 42)
	checkNotEqual(t, 42, t2)

	a := make([]any, 0)
	t2.CopyToSlice(a)
	checkEqual(t, []any{}, a)
	a = make([]any, 1)
	t2.CopyToSlice(a)
	checkEqual(t, []any{42}, a)
	a = make([]any, 2)
	t2.CopyToSlice(a)
	checkEqual(t, []any{42, `Answer`}, a)
	a = make([]any, 3)
	t2.CopyToSlice(a)
	checkEqual(t, []any{42, `Answer`, nil}, a)
}

func checkEqual(t *testing.T, exp, actual any) {
	t.Helper()
	if !utils.Equal(exp, actual) {
		t.Errorf("\n"+
			"Expected value didn't match the actual value:\n"+
			"Actual:   %v (%T)\n"+
			"Expected: %v (%T)", actual, actual, exp, exp)
	}
}

func checkNotEqual(t *testing.T, exp, actual any) {
	t.Helper()
	if utils.Equal(exp, actual) {
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
