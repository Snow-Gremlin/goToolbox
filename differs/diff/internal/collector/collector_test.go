package collector

import (
	"fmt"
	"testing"

	"goToolbox/differs"
	"goToolbox/utils"
)

func Test_Basics(t *testing.T) {
	col := New(19, 20)

	col.InsertAdded(1)
	col.InsertRemoved(1)
	col.InsertAdded(2)
	col.InsertRemoved(2)
	col.InsertEqual(3)

	col.InsertAdded(4)
	col.InsertEqual(2)
	col.InsertEqual(2)

	col.InsertRemoved(5)
	col.InsertEqual(2)
	col.InsertEqual(3)

	col.InsertRemoved(-6)
	col.InsertEqual(-6)
	col.InsertAdded(-6)

	r := col.Finish()
	intEqual(t, r.ACount(), 19, `Collection A Count`)
	intEqual(t, r.BCount(), 20, `Collection B Count`)
	intEqual(t, r.Count(), 7, `Collection Count`)
	intEqual(t, r.Total(), 27, `Collection Total`)
	readEqual(t, r, `=5 -5 =4 +4 =3 -3 +3`)
}

func Test_Error(t *testing.T) {
	col := New(6, 6)

	col.InsertAdded(1)
	col.InsertRemoved(1)
	col.InsertEqual(3)
	col.InsertRemoved(2)
	col.InsertAdded(2)
	col.InsertSubstitute(3)

	r := col.Finish()
	intEqual(t, r.ACount(), 6, `Collection A Count`)
	intEqual(t, r.BCount(), 6, `Collection B Count`)
	intEqual(t, r.Count(), 5, `Collection Count`)
	intEqual(t, r.Total(), 15, `Collection Total`)
	readEqual(t, r, `-5 +5 =3 -1 +1`)
}

func intEqual(t *testing.T, value, exp int, msg string) {
	if value != exp {
		t.Error(fmt.Sprint("Unexpected integer value:",
			"\n   Message:  ", msg,
			"\n   Value:    ", value,
			"\n   Expected: ", exp))
	}
}

func readEqual(t *testing.T, r differs.Result, exp string) {
	if result := utils.String(r); result != exp {
		t.Error(fmt.Sprint("Unexpected collection read:",
			"\n   Value:    ", result,
			"\n   Expected: ", exp))
	}
}
