package data

import (
	"testing"

	"goToolbox/differs"
)

func Test_Data_New(t *testing.T) {
	d := New(6, 4, func(aIndex, bIndex int) bool {
		return aIndex*2 == bIndex*3
	})
	checkSize(t, d, 6, 4)
	checkData(t, d, 0, 0, true)
	checkData(t, d, 3, 2, true)
	checkData(t, d, 6, 4, true)
	checkData(t, d, 1, 1, false)
	checkData(t, d, 2, 3, false)
	checkData(t, d, 4, 6, false)
}

func Test_Data_Strings(t *testing.T) {
	d := Strings(
		[]string{`cat`, `dog`, `apple`, `love`},
		[]string{`dog`, `cat`, `hate`})
	checkSize(t, d, 4, 3)
	checkData(t, d, 0, 0, false)
	checkData(t, d, 1, 1, false)
	checkData(t, d, 0, 1, true)
	checkData(t, d, 1, 0, true)
	checkData(t, d, 3, 2, false)
}

func Test_Data_Ints(t *testing.T) {
	d := Ints(
		[]int{3, 6, 1, 2, 4, 2},
		[]int{5, 4, 6, 2, 1})
	checkSize(t, d, 6, 5)
	checkData(t, d, 0, 0, false)
	checkData(t, d, 1, 1, false)
	checkData(t, d, 3, 3, true)
	checkData(t, d, 5, 3, true)
	checkData(t, d, 3, 2, false)
}

func Test_Data_Bytes(t *testing.T) {
	d := Bytes(
		[]byte{3, 6, 1, 2, 4, 2},
		[]byte{5, 4, 6, 2, 1})
	checkSize(t, d, 6, 5)
	checkData(t, d, 0, 0, false)
	checkData(t, d, 1, 1, false)
	checkData(t, d, 3, 3, true)
	checkData(t, d, 5, 3, true)
	checkData(t, d, 3, 2, false)
}

func Test_Data_Runes(t *testing.T) {
	d := Runes(
		[]rune{'H', 'e', 'l', 'l', 'o'},
		[]rune{'W', 'o', 'r', 'l', 'd'})
	checkSize(t, d, 5, 5)
	checkData(t, d, 0, 0, false)
	checkData(t, d, 1, 1, false)
	checkData(t, d, 3, 3, true)
	checkData(t, d, 4, 1, true)
	checkData(t, d, 1, 4, false)
}

func Test_Data_Chars(t *testing.T) {
	d := Chars(`Hello`, `World`)
	checkSize(t, d, 5, 5)
	checkData(t, d, 0, 0, false)
	checkData(t, d, 1, 1, false)
	checkData(t, d, 3, 3, true)
	checkData(t, d, 4, 1, true)
	checkData(t, d, 1, 4, false)
}

func Test_Data_RuneSlice(t *testing.T) {
	d := RuneSlice(
		[][]rune{[]rune(`cat`), []rune(`dog`), []rune(`apple`), []rune(`love`)},
		[][]rune{[]rune(`dog`), []rune(`cat`), []rune(`hate`)})
	checkSize(t, d, 4, 3)
	checkData(t, d, 0, 0, false)
	checkData(t, d, 1, 1, false)
	checkData(t, d, 0, 1, true)
	checkData(t, d, 1, 0, true)
	checkData(t, d, 3, 2, false)
}

func Test_Data_Any(t *testing.T) {
	d := Any([]any{3, `cat`, 1.2, pe(`dog`)},
		[]any{nil, 4, 3, `cat`, 1.2})
	checkSize(t, d, 4, 5)
	checkData(t, d, 0, 0, false)
	checkData(t, d, 1, 1, false)
	checkData(t, d, 1, 3, true)
	checkData(t, d, 2, 4, true)
	checkData(t, d, 3, 0, false)
}

func Test_Data_Regex(t *testing.T) {
	d := Regex([]string{`^1\d+$`, `^mi\w{2}$`},
		[]any{1, 22, 12, 123, `milk`, `mit`, `miak`, `mitch`})
	checkSize(t, d, 2, 8)
	checkData(t, d, 0, 0, false)
	checkData(t, d, 0, 1, false)
	checkData(t, d, 0, 2, true)
	checkData(t, d, 0, 3, true)
	checkData(t, d, 0, 4, false)
	checkData(t, d, 0, 5, false)
	checkData(t, d, 0, 6, false)
	checkData(t, d, 0, 7, false)
	checkData(t, d, 1, 0, false)
	checkData(t, d, 1, 1, false)
	checkData(t, d, 1, 2, false)
	checkData(t, d, 1, 3, false)
	checkData(t, d, 1, 4, true)
	checkData(t, d, 1, 5, false)
	checkData(t, d, 1, 6, true)
	checkData(t, d, 1, 7, false)
}

type pseudoEquatable struct {
	value string
}

func (p *pseudoEquatable) Equals(other any) bool {
	if p2, ok := other.(*pseudoEquatable); ok {
		if p == nil || p2 == nil {
			return p == p2
		}
		return p.value == p2.value
	}
	return false
}

func pe(value string) *pseudoEquatable {
	return &pseudoEquatable{value: value}
}

func checkSize(t *testing.T, d differs.Data, expA, expB int) {
	aCount, bCount := d.ACount(), d.BCount()
	if aCount != expA || bCount != expB {
		t.Errorf("\nUnexpected value:\n"+
			"\tA Count: actual = %d, expected = %d\n"+
			"\tB Count: actual = %d, expected = %d\n", aCount, expA, bCount, expB)
	}
}

func checkData(t *testing.T, d differs.Data, aIndex, bIndex int, exp bool) {
	if d.Equals(aIndex, bIndex) != exp {
		t.Errorf("\nUnexpected Equals Result:\n"+
			"\tExpected: %t\n"+
			"\tA Index:  %d\n"+
			"\tB Index:  %d\n", exp, aIndex, bIndex)
	}
}
