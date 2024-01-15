package simpleSet

import "testing"

func Test_SimpleSet(t *testing.T) {
	m := With(1, 4, 8, 4)
	checkEqual(t, 3, m.Count(), `Count after With`)

	checkEqual(t, true, m.Has(1), `Has(1) is set`)
	checkEqual(t, false, m.Has(2), `Has(2) is not set`)
	checkEqual(t, true, m.Has(4), `Has(4) is set`)

	checkEqual(t, `1, 4, 8`, m.ToString(), `ToSlice via ToString`)
	checkEqual(t, `1, 4, 8`, m.Clone().ToString(), `ToSlice`)

	m.Remove(4)
	checkEqual(t, `1, 8`, m.ToString(), `after Remove(4) when set`)
	m.Remove(4)
	checkEqual(t, `1, 8`, m.ToString(), `after Remove(4) when not set`)
	checkEqual(t, false, m.Has(4), `Has(4) after removed`)

	checkEqual(t, true, m.RemoveTest(8), `RemoveTest(8) when set`)
	checkEqual(t, `1`, m.ToString(), `after RemoveTest(8) when set`)
	checkEqual(t, false, m.RemoveTest(8), `RemoveTest(8) when not set`)
	checkEqual(t, `1`, m.ToString(), `after RemoveTest(8) when not set`)
	checkEqual(t, false, m.Has(8), `Has(8) after removed`)

	m = New[int]()
	checkEqual(t, 0, m.Count(), `Count after New`)
	checkEqual(t, ``, m.ToString(), `ToString after New`)

	m.Set(5)
	checkEqual(t, `5`, m.ToString(), `after Set(5) when not set`)
	m.Set(5)
	checkEqual(t, `5`, m.ToString(), `after Set(5) when set`)

	checkEqual(t, true, m.SetTest(7), `SetTest(7) when not set`)
	checkEqual(t, `5, 7`, m.ToString(), `after SetTest(7) when not set`)
	checkEqual(t, false, m.SetTest(7), `SetTest(7) when set`)
	checkEqual(t, `5, 7`, m.ToString(), `after SetTest(7) when set`)

	m.Set(10)
	m.Set(11)
	m.Set(13)
	m.Set(15)
	m.Set(17)
	checkEqual(t, `10, 11, 13, 15, 17, 5, 7`, m.ToString(), `after preparing for RemoveIf`)

	checkEqual(t, false, m.RemoveIf(nil), `RemoveIf with nil predicate`)
	checkEqual(t, true, m.RemoveIf(func(v int) bool { return v%5 == 0 }), `RemoveIf multiples of 5 when set`)
	checkEqual(t, `11, 13, 17, 7`, m.ToString(), `after RemoveIf multiples of 5 when set`)
	checkEqual(t, false, m.RemoveIf(func(v int) bool { return v%5 == 0 }), `RemoveIf multiples of 5 when not set`)
	checkEqual(t, `11, 13, 17, 7`, m.ToString(), `after RemoveIf multiples of 5 when not set`)
}

func checkEqual(t *testing.T, exp, actual any, msg string) {
	if actual != exp {
		t.Errorf("\nUnexpected result in SimpleSet:\n"+
			"\tMessage:  %s\n"+
			"\tExpected: %v\n"+
			"\tActual:   %v\n", msg, exp, actual)
	}
}
