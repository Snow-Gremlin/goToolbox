package terror

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"goToolbox/internal/liteUtils"
	"goToolbox/terrors"
	"goToolbox/terrors/stacked"
)

func pokeTheBear(e terrors.TError) {
	e.ResetStack(0)
}

func Test_TError(t *testing.T) {
	e1 := New(`Hello`)
	checkEqual(t, `Hello`, e1.Message())
	checkEqual(t, `Hello`, e1.Error())
	checkEqual(t, `Hello`, e1.String())
	checkMatch(t, `terror\.Test_TError\(`, e1.Stack())
	checkEqual(t, 0, len(e1.Context()))
	checkEqual(t, 0, len(e1.Unwrap()))
	checkMatch(t, `^Hello$`, e1)

	e2 := e1.Clone()
	e3 := e2.With(`cat`, 3).
		With(`dog`, `woof`).
		With(`cat`, 9)
	checkEqual(t, e3, e2)
	checkEqual(t, 0, len(e1.Context()))
	checkEqual(t, 2, len(e2.Context()))
	checkEqual(t, map[string]any{
		`cat`: 9,
		`dog`: `woof`,
	}, e2.Context())
	checkMatch(t, `^Hello \{cat: 9, dog: woof\}$`, e2)

	e4 := New(``)
	checkMatch(t, `^unknown error$`, e4)

	e5 := fmt.Errorf(`World - %w`, e4)
	e6 := e1.Clone()
	e7 := e6.WithError(e5).WithError(nil)
	checkEqual(t, e6, e7)
	checkEqual(t, 0, len(e1.Unwrap()))
	checkEqual(t, []error{e5}, e6.Unwrap())
	checkMatch(t, `^Hello: World - unknown error: unknown error$`, e6)

	e8 := fmt.Errorf(`Blue`)
	e6.WithError(e8)
	checkEqual(t, []error{e5, e8}, e6.Unwrap())
	checkMatch(t, `^Hello: \[World - unknown error: unknown error, Blue\]$`, e6)

	checkMatch(t, `terror\.Test_TError\(`, e6.Stack())
	checkMatch(t, `terror\.Test_TError\(`, e4.Stack())
	checkEqual(t, e4.Stack(), stacked.DeepestStacked(e6).Stack())
	checkEqual(t, false, Walk(nil).Next())

	pokeTheBear(e6)
	checkMatch(t, `terror\.pokeTheBear\(`, e6.Stack())

	e4.WithError(e6)
	checkMatch(t, `^Hello: \[World - unknown error: unknown error, Blue\]$`, e6)
}

func Test_TError_Equals(t *testing.T) {
	e1 := New(`hello`)
	e2 := e1.Clone()
	e3 := errors.New(`hello`)
	e4 := e1.Clone().With(`cat`, 3)
	e5 := e1.Clone().WithError(New(` world`))
	e6 := e1.Clone().WithError(errors.New(` world`))
	e7 := e1.Clone().WithError(New(` world`))
	e8 := New(`hello`)

	checkEqual(t, e1, e1)
	checkEqual(t, e2, e1)
	checkNotEqual(t, e3, e1)
	checkNotEqual(t, e4, e1)
	checkNotEqual(t, e5, e1)
	checkNotEqual(t, e6, e1)
	checkNotEqual(t, e6, e5)
	checkEqual(t, e7, e5)
	checkEqual(t, e8, e1)
	checkNotEqual(t, e8.Stack(), e1.Stack())
	checkNotEqual(t, `hello`, e1)
}

func Test_TError_PredefinedErrors(t *testing.T) {
	checkMatch(t, `^index out of bounds \{count: 8, index: 12\}$`, OutOfBounds(12, 8))
	checkMatch(t, `^collection contains no values \{action: Slap\}$`, EmptyCollection(`Slap`))
	checkMatch(t, `^invalid number of arguments \{count: 12, maximum: 4, usage: Scrap\}$`, InvalidArgCount(4, 12, `Scrap`))
	checkMatch(t, `^argument may not be nil \{name: Snap\}$`, NilArg(`Snap`))
	checkMatch(t, `^Collection was modified; iteration may not continue$`, UnstableIteration())

	checkEqual(t, nil, RecoveredPanic(nil))
	checkMatch(t, `^recovered panic \{recovered: Snick\}$`, RecoveredPanic(`Snick`))
	checkMatch(t, `^Scramble$`, RecoveredPanic(New(`Scramble`)))
	checkMatch(t, `^recovered panic \{recovered: 123\}$`, RecoveredPanic(123))
	checkMatch(t, `^recovered panic: Snoop$`, RecoveredPanic(errors.New(`Snoop`)))
}

func checkEqual(t *testing.T, exp, actual any) {
	t.Helper()
	if !liteUtils.Equal(exp, actual) {
		t.Errorf("\n"+
			"Expected value didn't match the actual value:\n"+
			"Actual:   %v (%T)\n"+
			"Expected: %v (%T)", actual, actual, exp, exp)
	}
}

func checkNotEqual(t *testing.T, exp, actual any) {
	t.Helper()
	if liteUtils.Equal(exp, actual) {
		t.Errorf("\n"+
			"Expected value shouldn't have matched the actual value:\n"+
			"Actual:   %v (%T)\n"+
			"Expected: %v (%T)", actual, actual, exp, exp)
	}
}

func checkMatch(t *testing.T, pattern string, actual any) {
	t.Helper()
	got := liteUtils.String(actual)
	match, err := regexp.MatchString(pattern, got)
	if err != nil {
		t.Fatalf(`Regex error: %v`, err)
	}
	if !match {
		t.Errorf("\n"+
			"Unexpected result:\n"+
			"Actual:  %v (%T)\n"+
			"Pattern: %v", got, actual, pattern)
	}
}
