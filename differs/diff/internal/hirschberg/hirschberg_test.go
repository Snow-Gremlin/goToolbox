package hirschberg

import (
	"testing"

	"goToolbox/differs/data"
	"goToolbox/differs/diff/internal"
	"goToolbox/differs/diff/internal/collector"
	"goToolbox/differs/diff/internal/container"
	"goToolbox/utils"
)

func Test_Hirschberg_NoReduce(t *testing.T) {
	checkAll(t, New(nil, -1, false))
}

func Test_Hirschberg_UseReduce(t *testing.T) {
	checkAll(t, New(nil, -1, true))
}

func Test_Hirschberg_Hybrid(t *testing.T) {
	d := New(New(nil, 6, false), -1, false)
	check(t, d, `kitten kitten kitten`, `sitting sitting sitting`,
		`-1 +1 =3 -1 +1 =1 +1 =1 -1 +1 =3 -1 +1 =1 +1 =1 -1 +1 =3 -1 +1 =1 +1`)
	check(t, d, `saturday saturday saturday`, `sunday sunday sunday`,
		`=1 -2 =1 -1 +1 =5 -2 =1 -1 +1 =5 -2 =1 -1 +1 =3`)
	check(t, d, `satxrday satxrday satxrday`, `sunday sunday sunday`,
		`=1 -4 +2 =5 -4 +2 =5 -4 +2 =3`)
}

func checkAll(t *testing.T, d internal.Algorithm) {
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

// checks the levenshtein distance algorithm
func check(t *testing.T, d internal.Algorithm, a, b, exp string) {
	dat := data.Chars(a, b)
	col := collector.New(dat.ACount(), dat.BCount())
	cont := container.New(dat)
	d.Diff(cont, col)
	r := col.Finish()
	result := utils.String(r)
	if exp != result {
		t.Error("Hirschberg returned unexpected result:",
			"\n   Input A:  ", a,
			"\n   Input B:  ", b,
			"\n   Expected: ", exp,
			"\n   Result:   ", result)
	}
}
