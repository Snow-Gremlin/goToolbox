package wagner

import (
	"fmt"
	"strings"
	"testing"

	"goToolbox/differs/data"
	"goToolbox/differs/diff/internal"
	"goToolbox/differs/diff/internal/collector"
	"goToolbox/differs/diff/internal/container"
	"goToolbox/utils"
)

func Test_Wagner(t *testing.T) {
	d := New(-1)
	check(t, d, `A`, `A`, `=1`)
	check(t, d, `A`, `B`, `-1 +1`)
	check(t, d, `A`, `AB`, `=1 +1`)
	check(t, d, `A`, `BA`, `+1 =1`)
	check(t, d, `AB`, `A`, `=1 -1`)
	check(t, d, `BA`, `A`, `-1 =1`)
	check(t, d, `kitten`, `sitting`, `-1 +1 =3 -1 +1 =1 +1`)
	check(t, d, `saturday`, `sunday`, `=1 -2 =1 -1 +1 =3`)
	check(t, d, `satxrday`, `sunday`, `=1 -4 +2 =3`)
	check(t, d, `ABC`, `ADB`, `=1 +1 =1 -1`)
}

func Test_NoResizeNeeded(t *testing.T) {
	d := New(25)
	boolEqual(t, noResizeNeeded(d, 5, 5), true, `5 x 5`)
	boolEqual(t, noResizeNeeded(d, 4, 4), true, `4 x 4`)
	boolEqual(t, noResizeNeeded(d, 2, 3), true, `2 x 3`)
	boolEqual(t, noResizeNeeded(d, 0, 0), true, `0 x 0`)
	boolEqual(t, noResizeNeeded(d, 3, 8), true, `3 x 8`)
	boolEqual(t, noResizeNeeded(d, 1, 25), true, `1 x 25`)
	boolEqual(t, noResizeNeeded(d, 5, 6), false, `5 x 6`)
	boolEqual(t, noResizeNeeded(d, 6, 6), false, `6 x 6`)
	boolEqual(t, noResizeNeeded(d, 3, 9), false, `3 x 9`)
	boolEqual(t, noResizeNeeded(d, 1, 26), false, `1 x 26`)
}

func noResizeNeeded(d internal.Algorithm, a, b int) bool {
	comp := data.Chars(strings.Repeat(`x`, a), strings.Repeat(`y`, b))
	return d.NoResizeNeeded(container.New(comp))
}

func boolEqual(t *testing.T, value, exp bool, msg string) {
	if value != exp {
		t.Error(fmt.Sprint("Unexpected boolean value:",
			"\n   Message:  ", msg,
			"\n   Value:    ", value,
			"\n   Expected: ", exp))
	}
}

// checks the levenshtein distance algorithm
func check(t *testing.T, d internal.Algorithm, a, b, exp string) {
	dat := data.Chars(a, b)
	col := collector.New(dat.ACount(), dat.BCount())
	cont := container.New(dat)
	d.Diff(cont, col)
	r := col.Finish()
	result := utils.String(r)
	if exp != result {
		t.Error("Wagner returned unexpected result:",
			"\n   Input A:  ", a,
			"\n   Input B:  ", b,
			"\n   Expected: ", exp,
			"\n   Result:   ", result)
	}
}
