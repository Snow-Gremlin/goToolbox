package result

import (
	"fmt"
	"testing"

	"goToolbox/collections/stack"
	"goToolbox/collections/tuple2"
	"goToolbox/differs/step"
)

func Test_Result(t *testing.T) {
	s := stack.With(
		tuple2.New(step.Added, 3),
		tuple2.New(step.Removed, 2))
	r := New(s.Enumerate(), s.Count(), 3, 2, 5, 3, 2)

	if r.Count() != 2 {
		t.Errorf(`the count was %d but it should have been 2`, r.Count())
	}

	if r.ACount() != 3 {
		t.Errorf(`the A count was %d but it should have been 3`, r.ACount())
	}

	if r.BCount() != 2 {
		t.Errorf(`the B count was %d but it should have been 2`, r.BCount())
	}

	if r.Total() != 5 {
		t.Errorf(`the total was %d but it should have been 5`, r.Total())
	}

	if r.AddedCount() != 3 {
		t.Errorf(`the added count was %d but it should have been 3`, r.AddedCount())
	}

	if r.RemovedCount() != 2 {
		t.Errorf(`the removed count was %d but it should have been 2`, r.RemovedCount())
	}

	if !r.HasDiff() {
		t.Error(`the diff should have been true`)
	}

	if result := r.Enumerate().Join(`, `); result != `[+, 3], [-, 2]` {
		t.Errorf(`the enumerate was %q but it should have been "[+, 3], [-, 2]"`, result)
	}

	if str := fmt.Sprint(r); str != `+3 -2` {
		t.Errorf(`the string was not expected, it was %q`, str)
	}
}
