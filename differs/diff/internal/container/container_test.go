package container

import (
	"fmt"
	"testing"

	"goToolbox/differs/diff/internal"
	"goToolbox/differs/diff/internal/collector"
	"goToolbox/utils"
)

func Test_Equals(t *testing.T) {
	cont := newCont(`cat`, `kitten`)
	check(t, cont, `cat`, `kitten`, `0, 3, 0, 6, false`)
	intEqual(t, cont.ACount(), 3, `ACount`)
	intEqual(t, cont.BCount(), 6, `BCount`)
	boolEqual(t, cont.Equals(0, 0), false, `Equal(0, 0)`)
	boolEqual(t, cont.Equals(2, 2), true, `Equal(2, 2)`)
	boolEqual(t, cont.Equals(2, 3), true, `Equal(2, 3)`)
	boolEqual(t, cont.Equals(1, 2), false, `Equal(1, 2)`)
	intEqual(t, cont.SubstitutionCost(0, 0), internal.SubstitutionCost, `SubstitutionCost(0, 0)`)
	intEqual(t, cont.SubstitutionCost(2, 2), internal.EqualCost, `SubstitutionCost(2, 2)`)
}

func Test_Equals_Reversed(t *testing.T) {
	cont := reverse(newCont(`cat`, `kitten`))
	check(t, cont, `tac`, `nettik`, `0, 3, 0, 6, true`)
	intEqual(t, cont.ACount(), 3, `ACount`)
	intEqual(t, cont.BCount(), 6, `BCount`)
	boolEqual(t, cont.Equals(0, 0), false, `Equal(0, 0)`)
	boolEqual(t, cont.Equals(0, 2), true, `Equal(0, 2)`)
	boolEqual(t, cont.Equals(0, 3), true, `Equal(0, 3)`)
	boolEqual(t, cont.Equals(1, 2), false, `Equal(1, 2)`)
	intEqual(t, cont.SubstitutionCost(0, 0), internal.SubstitutionCost, `SubstitutionCost(0, 0)`)
	intEqual(t, cont.SubstitutionCost(0, 2), internal.EqualCost, `SubstitutionCost(2, 2)`)
}

func Test_Sub(t *testing.T) {
	cont := newCont(`abcdef`, `ghi`)
	check(t, cont, `abcdef`, `ghi`, `0, 6, 0, 3, false`)
	subCheck(t, cont, 0, 3, 0, 3, false, `abc`, `ghi`)
	subCheck(t, cont, 1, 4, 1, 3, false, `bcd`, `hi`)
	subCheck(t, cont, 0, 3, 0, 3, true, `cba`, `ihg`)
	subCheck(t, cont, 2, 5, 1, 3, true, `edc`, `ih`)

	sub := cont.Sub(2, 5, 1, 3, true)
	check(t, sub, `edc`, `ih`, `2, 3, 1, 2, true`)
}

func Test_Sub_Reversed(t *testing.T) {
	cont := reverse(newCont(`abcdef`, `ghi`))
	check(t, cont, `fedcba`, `ihg`, `0, 6, 0, 3, true`)
	subCheck(t, cont, 0, 3, 0, 3, false, `fed`, `ihg`)
	subCheck(t, cont, 1, 4, 1, 3, false, `edc`, `hg`)
	subCheck(t, cont, 0, 3, 0, 3, true, `def`, `ghi`)
	subCheck(t, cont, 2, 5, 1, 3, true, `bcd`, `gh`)
}

func Test_Reduce(t *testing.T) {
	reduceCheck(t, newCont(`abc`, `abc`), ``, ``, 3, 0)
	reduceCheck(t, newCont(`abc`, `def`), `abc`, `def`, 0, 0)
	reduceCheck(t, newCont(`abc`, `aef`), `bc`, `ef`, 1, 0)
	reduceCheck(t, newCont(`abc`, `dec`), `ab`, `de`, 0, 1)
	reduceCheck(t, newCont(`abc`, `ac`), `b`, ``, 1, 1)
	reduceCheck(t, newCont(`ac`, `abc`), ``, `b`, 1, 1)
	reduceCheck(t, newCont(`abcd`, `acd`), `b`, ``, 1, 2)
	reduceCheck(t, newCont(`abcd`, `abd`), `c`, ``, 2, 1)
	reduceCheck(t, newCont(`abc`, ``), `abc`, ``, 0, 0)
	reduceCheck(t, newCont(``, `abc`), ``, `abc`, 0, 0)
}

func Test_Reduce_Reversed(t *testing.T) {
	reduceCheck(t, reverse(newCont(`abc`, `abc`)), ``, ``, 0, 3)
	reduceCheck(t, reverse(newCont(`abc`, `def`)), `cba`, `fed`, 0, 0)
	reduceCheck(t, reverse(newCont(`abc`, `aef`)), `cb`, `fe`, 0, 1)
	reduceCheck(t, reverse(newCont(`abc`, `dec`)), `ba`, `ed`, 1, 0)
	reduceCheck(t, reverse(newCont(`abc`, `ac`)), `b`, ``, 1, 1)
	reduceCheck(t, reverse(newCont(`ac`, `abc`)), ``, `b`, 1, 1)
	reduceCheck(t, reverse(newCont(`abcd`, `acd`)), `b`, ``, 2, 1)
	reduceCheck(t, reverse(newCont(`abcd`, `abd`)), `c`, ``, 1, 2)
	reduceCheck(t, reverse(newCont(`abc`, ``)), `cba`, ``, 0, 0)
	reduceCheck(t, reverse(newCont(``, `abc`)), ``, `cba`, 0, 0)
}

func Test_EndCase(t *testing.T) {
	endCaseCheck(t, newCont(`abc`, `abc`), false, ``)
	endCaseCheck(t, newCont(``, ``), true, ``)

	endCaseCheck(t, newCont(`a`, ``), true, `-1`)
	endCaseCheck(t, newCont(`ab`, ``), true, `-2`)
	endCaseCheck(t, newCont(`abc`, ``), true, `-3`)

	endCaseCheck(t, newCont(``, `a`), true, `+1`)
	endCaseCheck(t, newCont(``, `ab`), true, `+2`)
	endCaseCheck(t, newCont(``, `abc`), true, `+3`)

	endCaseCheck(t, newCont(`abc`, `a`), true, `=1 -2`)
	endCaseCheck(t, newCont(`abc`, `b`), true, `-1 =1 -1`)
	endCaseCheck(t, newCont(`abc`, `c`), true, `-2 =1`)
	endCaseCheck(t, newCont(`abc`, `d`), true, `-3 +1`)

	endCaseCheck(t, newCont(`a`, `abc`), true, `=1 +2`)
	endCaseCheck(t, newCont(`b`, `abc`), true, `+1 =1 +1`)
	endCaseCheck(t, newCont(`c`, `abc`), true, `+2 =1`)
	endCaseCheck(t, newCont(`d`, `abc`), true, `-1 +3`)
}

func Test_EndCase_Reverse(t *testing.T) {
	endCaseCheck(t, reverse(newCont(`abc`, `abc`)), false, ``)
	endCaseCheck(t, reverse(newCont(``, ``)), true, ``)

	endCaseCheck(t, reverse(newCont(`a`, ``)), true, `-1`)
	endCaseCheck(t, reverse(newCont(`ab`, ``)), true, `-2`)
	endCaseCheck(t, reverse(newCont(`abc`, ``)), true, `-3`)

	endCaseCheck(t, reverse(newCont(``, `a`)), true, `+1`)
	endCaseCheck(t, reverse(newCont(``, `ab`)), true, `+2`)
	endCaseCheck(t, reverse(newCont(``, `abc`)), true, `+3`)

	endCaseCheck(t, reverse(newCont(`abc`, `a`)), true, `-2 =1`)
	endCaseCheck(t, reverse(newCont(`abc`, `b`)), true, `-1 =1 -1`)
	endCaseCheck(t, reverse(newCont(`abc`, `c`)), true, `=1 -2`)
	endCaseCheck(t, reverse(newCont(`abc`, `d`)), true, `-3 +1`)

	endCaseCheck(t, reverse(newCont(`a`, `abc`)), true, `+2 =1`)
	endCaseCheck(t, reverse(newCont(`b`, `abc`)), true, `+1 =1 +1`)
	endCaseCheck(t, reverse(newCont(`c`, `abc`)), true, `=1 +2`)
	endCaseCheck(t, reverse(newCont(`d`, `abc`)), true, `-1 +3`)
}

func boolEqual(t *testing.T, value, exp bool, msg string) {
	if value != exp {
		t.Error(fmt.Sprint("Unexpected boolean value:",
			"\n   Message:  ", msg,
			"\n   Value:    ", value,
			"\n   Expected: ", exp))
	}
}

func intEqual(t *testing.T, value, exp int, msg string) {
	if value != exp {
		t.Error(fmt.Sprint("Unexpected integer value:",
			"\n   Message:  ", msg,
			"\n   Value:    ", value,
			"\n   Expected: ", exp))
	}
}

type pseudoData struct {
	a, b string
}

func (pd *pseudoData) ACount() int {
	return len(pd.a)
}

func (pd *pseudoData) BCount() int {
	return len(pd.b)
}

func (pd *pseudoData) Equals(aIndex, bIndex int) bool {
	return pd.a[aIndex] == pd.b[bIndex]
}

func newCont(a, b string) internal.Container {
	return New(&pseudoData{a: a, b: b})
}

func reverse(c internal.Container) internal.Container {
	i := c.(*containerImp)
	return c.Sub(0, i.aCount, 0, i.bCount, !i.reverse)
}

func (cont *containerImp) AAdjust(aIndex int) int {
	if cont.reverse {
		return cont.aCount - 1 - aIndex + cont.aOffset
	}
	return aIndex + cont.aOffset
}

func (cont *containerImp) BAdjust(bIndex int) int {
	if cont.reverse {
		return cont.bCount - 1 - bIndex + cont.bOffset
	}
	return bIndex + cont.bOffset
}

func (cont *containerImp) AParts() string {
	parts := make([]byte, cont.aCount)
	pd := cont.data.(*pseudoData)
	for i := 0; i < cont.aCount; i++ {
		parts[i] = pd.a[cont.AAdjust(i)]
	}
	return string(parts)
}

func (cont *containerImp) BParts() string {
	parts := make([]byte, cont.bCount)
	pd := cont.data.(*pseudoData)
	for j := 0; j < cont.bCount; j++ {
		parts[j] = pd.b[cont.BAdjust(j)]
	}
	return string(parts)
}

func (cont *containerImp) String() string {
	return fmt.Sprintf(`%d, %d, %d, %d, %t`,
		cont.aOffset, cont.aCount, cont.bOffset, cont.bCount, cont.reverse)
}

func check(t *testing.T, cont internal.Container, expA, expB, expStr string) {
	t.Helper()
	i := cont.(*containerImp)
	resultA := i.AParts()
	resultB := i.BParts()
	resultStr := i.String()
	if (resultA != expA) || (resultB != expB) || (resultStr != expStr) {
		t.Error(fmt.Sprint(
			"Unexpected resulting container:",
			"\n   Container:  ", cont,
			"\n   Result A:   ", resultA, " => ", expA,
			"\n   Result B:   ", resultB, " => ", expB,
			"\n   Result Str: ", resultStr, " => ", expStr))
	}
}

func subCheck(t *testing.T, cont internal.Container, aLow, aHigh, bLow, bHigh int, reverse bool, expA, expB string) {
	t.Helper()
	sub := cont.Sub(aLow, aHigh, bLow, bHigh, reverse)
	i := sub.(*containerImp)
	resultA := i.AParts()
	resultB := i.BParts()
	if (resultA != expA) || (resultB != expB) {
		t.Error(fmt.Sprint(
			"Unexpected results from Reduce:",
			"\n   Original: ", cont,
			"\n   Sub:      ", sub,
			"\n   A Parts:  ", resultA, " => ", expA,
			"\n   B Parts:  ", resultB, " => ", expB))
	}
}

func reduceCheck(t *testing.T, cont internal.Container, expA, expB string, expBefore, expAfter int) {
	t.Helper()
	sub, before, after := cont.Reduce()
	i := sub.(*containerImp)
	resultA := i.AParts()
	resultB := i.BParts()
	if (before != expBefore) || (after != expAfter) || (resultA != expA) || (resultB != expB) {
		t.Error(fmt.Sprint(
			"Unexpected results from Reduce:",
			"\n   Original: ", cont,
			"\n   Reduces:  ", sub,
			"\n   A Parts:  ", resultA, " => ", expA,
			"\n   B Parts:  ", resultB, " => ", expB,
			"\n   Before:   ", before, " => ", expBefore,
			"\n   After:    ", after, " => ", expAfter))
	}
}

func endCaseCheck(t *testing.T, cont internal.Container, expBool bool, expCol string) {
	t.Helper()
	col := collector.New(-1, -1)
	resultBool := cont.EndCase(col)
	r := col.Finish()
	resultCol := utils.String(r)
	if (resultBool != expBool) || (resultCol != expCol) {
		t.Error(fmt.Sprint("Unexpected EndCase results:",
			"\n   Container:  ", cont,
			"\n   Result:     ", resultBool, " => ", expBool,
			"\n   Collection: ", resultCol, " => ", expCol))
	}
}
