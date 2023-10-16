package optional

import (
	"strings"
	"testing"

	"goToolbox/utils"
)

func Test_Optional_oneArg(t *testing.T) {
	checkEqual(t, -1, oneArg([]int{}, -1, `bork`))
	checkEqual(t, 42, oneArg([]int{}, 42, `bork`))
	checkEqual(t, 42, oneArg([]int{42}, -1, `bork`))
	checkEqual(t, -1, oneArg([]int{-1}, 42, `bork`))
	checkPanic(t, `invalid number of arguments {count: 2, maximum: 1, usage: bork}`,
		func() { oneArg([]int{1, 2}, -1, `bork`) })
}

func Test_Optional_twoArgs(t *testing.T) {
	a, b := twoArgs([]int{}, -1, -2, `bork`)
	checkEqual(t, -1, a)
	checkEqual(t, -2, b)
	a, b = twoArgs([]int{}, 42, 24, `bork`)
	checkEqual(t, 42, a)
	checkEqual(t, 24, b)
	a, b = twoArgs([]int{42}, -1, -2, `bork`)
	checkEqual(t, 42, a)
	checkEqual(t, -2, b)
	a, b = twoArgs([]int{-1}, 24, 42, `bork`)
	checkEqual(t, -1, a)
	checkEqual(t, 42, b)
	a, b = twoArgs([]int{13, 26}, -1, -2, `bork`)
	checkEqual(t, 13, a)
	checkEqual(t, 26, b)
	a, b = twoArgs([]int{-2, -1}, 24, 42, `bork`)
	checkEqual(t, -2, a)
	checkEqual(t, -1, b)
	checkPanic(t, `invalid number of arguments {count: 3, maximum: 2, usage: bork}`,
		func() { twoArgs([]int{1, 2, 3}, -1, -2, `bork`) })
}

func Test_Optional_SizeAndCapacity(t *testing.T) {
	size, capacity := SizeAndCapacity([]int{})
	checkEqual(t, 0, size)
	checkEqual(t, 0, capacity)
	size, capacity = SizeAndCapacity([]int{-4})
	checkEqual(t, 0, size)
	checkEqual(t, 0, capacity)
	size, capacity = SizeAndCapacity([]int{42})
	checkEqual(t, 42, size)
	checkEqual(t, 42, capacity)
	size, capacity = SizeAndCapacity([]int{-4, -3})
	checkEqual(t, 0, size)
	checkEqual(t, 0, capacity)
	size, capacity = SizeAndCapacity([]int{42, -3})
	checkEqual(t, 42, size)
	checkEqual(t, 42, capacity)
	size, capacity = SizeAndCapacity([]int{-4, 24})
	checkEqual(t, 0, size)
	checkEqual(t, 24, capacity)
	size, capacity = SizeAndCapacity([]int{42, 24})
	checkEqual(t, 42, size)
	checkEqual(t, 42, capacity)
	size, capacity = SizeAndCapacity([]int{13, 26})
	checkEqual(t, 13, size)
	checkEqual(t, 26, capacity)
	checkPanic(t, `invalid number of arguments {count: 4, maximum: 2, usage: size and capacity}`,
		func() { SizeAndCapacity([]int{1, 2, 3, 4}) })
}

func Test_Optional_Capacity(t *testing.T) {
	checkEqual(t, 0, Capacity([]int{}))
	checkEqual(t, 0, Capacity([]int{-4}))
	checkEqual(t, 42, Capacity([]int{42}))
	checkPanic(t, `invalid number of arguments {count: 2, maximum: 1, usage: capacity}`,
		func() { Capacity([]int{1, 2}) })
}

func Test_Optional_Size(t *testing.T) {
	checkEqual(t, 0, Size([]int{}))
	checkEqual(t, 0, Size([]int{-4}))
	checkEqual(t, 42, Size([]int{42}))
	checkPanic(t, `invalid number of arguments {count: 2, maximum: 1, usage: size}`,
		func() { Size([]int{1, 2}) })
}

func Test_Optional_After(t *testing.T) {
	checkEqual(t, -1, After([]int{}))
	checkEqual(t, -1, After([]int{-4}))
	checkEqual(t, 0, After([]int{0}))
	checkEqual(t, 42, After([]int{42}))
	checkPanic(t, `invalid number of arguments {count: 2, maximum: 1, usage: after index}`,
		func() { After([]int{1, 2}) })
}

func Test_Optional_Comparer(t *testing.T) {
	cmp1 := Comparer([]utils.Comparer[int]{})
	checkEqual(t, -1, cmp1(1, 2))
	checkEqual(t, 0, cmp1(3, 3))
	checkEqual(t, 1, cmp1(2, 1))
	cmp2 := Comparer([]utils.Comparer[string]{})
	checkEqual(t, -1, cmp2(`apple`, `cat`))
	checkEqual(t, 0, cmp2(`dog`, `dog`))
	checkEqual(t, 1, cmp2(`cat`, `apple`))
	cmp2 = Comparer([]utils.Comparer[string]{utils.Descender(strings.Compare)})
	checkEqual(t, 1, cmp2(`apple`, `cat`))
	checkEqual(t, 0, cmp2(`dog`, `dog`))
	checkEqual(t, -1, cmp2(`cat`, `apple`))
	cmp2 = Comparer([]utils.Comparer[string]{nil})
	checkEqual(t, -1, cmp2(`apple`, `cat`))
	checkEqual(t, 0, cmp2(`dog`, `dog`))
	checkEqual(t, 1, cmp2(`cat`, `apple`))
	checkPanic(t, `invalid number of arguments {count: 2, maximum: 1, usage: comparer}`,
		func() { Comparer([]utils.Comparer[string]{nil, nil}) })
	checkPanic(t, `must provide a comparer to compare this type {type: []int}`,
		func() { Comparer([]utils.Comparer[[]int]{}) })
}

func checkEqual(t *testing.T, expected, actual any) {
	if !utils.Equal(expected, actual) {
		t.Errorf("\n"+
			"Unexpected value:\n"+
			"\tExpected: %v (%T)\n"+
			"\tActual:   %v (%T)\n",
			expected, expected, actual, actual)
	}
}

func checkPanic(t *testing.T, expected string, handle func()) {
	actual := func() (msg string) {
		defer func() { msg = utils.String(recover()) }()
		handle()
		return ``
	}()
	checkEqual(t, expected, actual)
}
