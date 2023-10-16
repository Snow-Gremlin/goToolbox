package utils

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"
)

func checkIsZero[T any](t *testing.T, exp bool, value T) {
	if exp != IsZero(value) {
		t.Errorf("\n"+
			"Unexpected value from IsZero:\n"+
			"\tType:     %T\n"+
			"\tExpected: %t\n", value, exp)
	}
}

func checkAreEqual[T comparable](t *testing.T, actual, expected T) {
	if actual != expected {
		t.Errorf("\n"+
			"Actual value did not match expected value:\n"+
			"\tActual:   %v\n"+
			"\tExpected: %v\n", actual, expected)
	}
}

func Test_Utils_RemoveZero(t *testing.T) {
	s1 := []int{0, 12, 0, 3, 0}
	s2 := RemoveZeros(s1)

	// Check s1 did not change
	checkAreEqual(t, len(s1), 5)
	checkAreEqual(t, cap(s1), 5)
	checkIsZero(t, true, s1[0])
	checkIsZero(t, false, s1[1])
	checkIsZero(t, true, s1[2])
	checkIsZero(t, false, s1[3])
	checkIsZero(t, true, s1[4])

	// Check s2 has no zeros
	checkAreEqual(t, len(s2), 2)
	checkAreEqual(t, cap(s2), 2)
	checkIsZero(t, false, s2[0])
	checkIsZero(t, false, s2[1])
}

func Test_Utils_SetToZero(t *testing.T) {
	s := []int{1, 2, 3, 4}
	SetToZero(s, 1, 4)

	checkAreEqual(t, len(s), 4)
	checkIsZero(t, false, s[0])
	checkIsZero(t, true, s[1])
	checkIsZero(t, true, s[2])
	checkIsZero(t, true, s[3])
}

func checkIsNil[T any](t *testing.T, expNil, expOk bool, value T) {
	if isNil, ok := TryIsNil(value); expNil != isNil || expOk != ok {
		t.Errorf("\n"+
			"Unexpected value from TryIsNil:\n"+
			"\tType:     %T\n"+
			"\tActual:   %t, %t\n"+
			"\tExpected: %t, %t\n", value, isNil, ok, expNil, expOk)
	}

	if isNil := IsNil(value); expNil != isNil {
		t.Errorf("\n"+
			"Unexpected value from IsNil:\n"+
			"\tType:     %T\n"+
			"\tActual:   %t\n"+
			"\tExpected: %t\n", value, isNil, expNil)
	}
}

func Test_Utils_IsNil(t *testing.T) {
	v1 := 12
	checkIsNil(t, false, false, v1)
	checkIsNil(t, false, true, &v1)
	checkIsNil(t, true, true, (*int)(nil))

	var v2 error
	checkIsNil(t, true, true, v2)
	v3 := (*strconv.NumError)(nil)
	checkIsNil(t, true, true, v3)
	v2 = v3
	checkIsNil(t, true, true, v2)
	v3 = &strconv.NumError{Func: `Oops`}
	checkIsNil(t, false, true, v3)
	v2 = v3
	checkIsNil(t, false, true, v2)

	var v4 []int
	checkIsNil(t, true, true, v4)
	v4 = []int{}
	checkIsNil(t, false, true, v4)

	var v5 map[string]int
	checkIsNil(t, true, true, v5)
	v5 = map[string]int{}
	checkIsNil(t, false, true, v5)
}

func checkLength[T any](t *testing.T, expLen int, expOk bool, value T) {
	len, ok := Length(value)
	if expLen != len || expOk != ok {
		t.Errorf("\n"+
			"Unexpected value from Length:\n"+
			"\tType:     %T\n"+
			"\tActual:   %d, %t\n"+
			"\tExpected: %d, %t\n", value, len, ok, expLen, expOk)
	}
}

type (
	lenObj    struct{ len int }
	lengthObj struct{ length int }
	countObj  struct{ count int }
)

func (n lenObj) Len() int       { return n.len }
func (n lengthObj) Length() int { return n.length }
func (n countObj) Count() int   { return n.count }

func Test_Utils_Length(t *testing.T) {
	v1 := 12
	checkLength(t, 0, false, v1)
	checkLength(t, 0, false, &v1)
	checkLength(t, 0, false, (*int)(nil))
	checkLength[any](t, 0, false, nil)

	var v2 []int
	checkLength(t, 0, true, v2)
	v2 = []int{1, 2, 3, 4}
	checkLength(t, 4, true, v2)

	var v3 map[string]int
	checkLength(t, 0, true, v3)
	v3 = map[string]int{`A`: 1, `B`: 2, `C`: 3, `D`: 4}
	checkLength(t, 4, true, v3)

	var v4 string
	checkLength(t, 0, true, v4)
	v4 = `Pudding`
	checkLength(t, 7, true, v4)

	obj1 := lenObj{len: 14}
	checkLength(t, 14, true, obj1)

	obj2 := lengthObj{length: 27}
	checkLength(t, 27, true, obj2)

	obj3 := countObj{count: 336}
	checkLength(t, 336, true, obj3)
	checkLength[interface{ Count() int }](t, 0, false, nil)
}

func checkSortedKeys[TKey comparable, TValue any, TMap ~map[TKey]TValue](t *testing.T, value TMap, exp string, cmp ...Comparer[TKey]) {
	keys := SortedKeys(value, cmp...)
	if result := strings.Join(Strings(keys), `, `); exp != result {
		t.Errorf("\n"+
			"Unexpected value from SortedKeys:\n"+
			"\tKey Type: %T\n"+
			"\tActual:   %s\n"+
			"\tExpected: %s\n", Zero[TKey](), result, exp)
	}
}

func catchPanic(handle func()) (msg string) {
	defer func() { msg = String(recover()) }()
	handle()
	return ``
}

func Test_Utils_SortedKeys(t *testing.T) {
	checkSortedKeys(t,
		map[int]float64{3: 4.3, 1: 2.16, 6: 333.333, 5: 12.34},
		`1, 3, 5, 6`)

	checkSortedKeys(t,
		map[float64]int{4.3: 3, 2.16: 1, 333.333: 6, 12.34: 5},
		`2.16, 4.3, 12.34, 333.333`)

	checkSortedKeys(t,
		map[rune]float64{'k': 4.3, 'a': 2.16, 'q': 333.333, 'H': 12.34},
		`72, 97, 107, 113`)

	checkSortedKeys(t,
		map[string]float64{"cat": 4.3, "pig": 2.16, "Dog": 333.333, "apple": 12.34},
		`Dog, apple, cat, pig`)

	checkSortedKeys(t,
		map[int]float64{3: 4.3, 1: 2.16, 6: 333.333, 5: 12.34},
		`6, 5, 3, 1`, Descender(OrderedComparer[int]()))

	checkSortedKeys(t,
		map[float64]int{4.3: 3, 2.16: 1, 333.333: 6, 12.34: 5},
		`333.333, 12.34, 4.3, 2.16`, Descender(OrderedComparer[float64]()))

	checkSortedKeys(t,
		map[rune]float64{'k': 4.3, 'a': 2.16, 'q': 333.333, 'H': 12.34},
		`113, 107, 97, 72`, Descender(OrderedComparer[rune]()))

	checkSortedKeys(t,
		map[string]float64{"cat": 4.3, "pig": 2.16, "Dog": 333.333, "apple": 12.34},
		`pig, cat, apple, Dog`, Descender(OrderedComparer[string]()))

	errStr := catchPanic(func() {
		SortedKeys(map[string]float64{"apple": 12.34}, OrderedComparer[string](), OrderedComparer[string]())
	})
	checkEqual(t, `invalid number of arguments {count: 2, maximum: 1, usage: comparer}`, errStr, true)

	errStr = catchPanic(func() {
		SortedKeys(map[complex128]float64{}, nil)
	})
	checkEqual(t, `must provide a comparer to compare this type {type: complex128}`, errStr, true)
}

func checkComparer[T any](t *testing.T, cmp Comparer[T], x, y T, exp int) {
	if actual := cmp(x, y); exp != actual {
		t.Errorf("\n"+
			"Unexpected value from Comparer:\n"+
			"\tKey Type:    %T\n"+
			"\tLeft Value:  %v\n"+
			"\tRight Value: %v\n"+
			"\tActual:      %d\n"+
			"\tExpected:    %d\n", x, x, y, actual, exp)
	}
}

func Test_Utils_OrderedComparer(t *testing.T) {
	c := OrderedComparer[string]()
	checkComparer(t, c, `banana`, `cat`, -1)
	checkComparer(t, c, `cat`, `banana`, 1)
	checkComparer(t, c, `banana`, `banana`, 0)
	checkComparer(t, c, `cat`, `cat`, 0)
}

type pseudoComparable struct {
	name string
}

func (c *pseudoComparable) CompareTo(other *pseudoComparable) int {
	if c == nil {
		if other == nil {
			return 0
		}
		return -1
	}
	if other == nil {
		return 1
	}
	return strings.Compare(c.name, other.name)
}

func Test_Utils_Comparable(t *testing.T) {
	c := ComparableComparer[*pseudoComparable]()
	pc0 := (*pseudoComparable)(nil)
	pc1 := &pseudoComparable{name: `banana`}
	pc2 := &pseudoComparable{name: `cat`}

	checkComparer(t, c, pc0, pc0, 0)
	checkComparer(t, c, pc0, pc1, -1)
	checkComparer(t, c, pc0, pc2, -1)

	checkComparer(t, c, pc1, pc0, 1)
	checkComparer(t, c, pc1, pc1, 0)
	checkComparer(t, c, pc1, pc2, -1)

	checkComparer(t, c, pc2, pc0, 1)
	checkComparer(t, c, pc2, pc1, 1)
	checkComparer(t, c, pc2, pc2, 0)
}

func Test_Utils_ComparerFromLess(t *testing.T) {
	cmp := ComparerForLess(func(x, y string) bool {
		return len(x) < len(y)
	})

	values := []string{`cat`, `dogs`, `doggo`, `apple`, `ox`}
	exp := `ox, cat, dogs, doggo, apple`
	c := slices.Clone(values)
	slices.SortFunc(c, cmp)
	if result := strings.Join(c, `, `); result != exp {
		t.Errorf("\n"+
			"Unexpected value from ComparerFromLess sort:\n"+
			"\tActual:   %s\n"+
			"\tExpected: %s\n", result, exp)
	}
}

func Test_Utils_EpsilonComparer(t *testing.T) {
	cmp := EpsilonComparer(0.01)
	checkComparer(t, cmp, 0.0, 0.0, 0)
	checkComparer(t, cmp, 1.0, 1.0, 0)
	checkComparer(t, cmp, -1.0, -1.0, 0)

	checkComparer(t, cmp, 0.0, 0.005, 0)
	checkComparer(t, cmp, 0.0, 0.01, 0)
	checkComparer(t, cmp, 0.0, 0.02, -1)
	checkComparer(t, cmp, 0.0, 1.0, -1)
	checkComparer(t, cmp, 0.0, -0.005, 0)
	checkComparer(t, cmp, 0.0, -0.01, 0)
	checkComparer(t, cmp, 0.0, -0.02, 1)
	checkComparer(t, cmp, 0.0, -1.0, 1)

	checkComparer(t, cmp, 0.005, 0.0, 0)
	checkComparer(t, cmp, 0.01, 0.0, 0)
	checkComparer(t, cmp, 0.02, 0.0, 1)
	checkComparer(t, cmp, 1.0, 0.0, 1)
	checkComparer(t, cmp, -0.005, 0.0, 0)
	checkComparer(t, cmp, -0.01, 0.0, 0)
	checkComparer(t, cmp, -0.02, 0.0, -1)
	checkComparer(t, cmp, -1.0, 0.0, -1)

	cmp = EpsilonComparer(-1.0) // defaults to ordered comparer, epsilon = 0
	checkComparer(t, cmp, 0.0, 0.0, 0)
	checkComparer(t, cmp, 1.0, 1.0, 0)
	checkComparer(t, cmp, -1.0, -1.0, 0)
}

func Test_Utils_DefaultComparer(t *testing.T) {
	checkComparer(t, DefaultComparer[int](), 1, 3, -1)
	checkComparer(t, DefaultComparer[int8](), 1, 3, -1)
	checkComparer(t, DefaultComparer[int16](), 1, 3, -1)
	checkComparer(t, DefaultComparer[int32](), 1, 3, -1)
	checkComparer(t, DefaultComparer[int64](), 1, 3, -1)

	checkComparer(t, DefaultComparer[uint](), 1, 3, -1)
	checkComparer(t, DefaultComparer[uint8](), 1, 3, -1)
	checkComparer(t, DefaultComparer[uint16](), 1, 3, -1)
	checkComparer(t, DefaultComparer[uint32](), 1, 3, -1)
	checkComparer(t, DefaultComparer[uint64](), 1, 3, -1)

	checkComparer(t, DefaultComparer[float32](), 1.0, 3.0, -1)
	checkComparer(t, DefaultComparer[float64](), 1.0, 3.0, -1)

	checkComparer(t, DefaultComparer[uintptr](), 1, 3, -1)
	checkComparer(t, DefaultComparer[string](), `apple`, `dog`, -1)
	checkComparer(t, DefaultComparer[rune](), 'A', 'B', -1)
	checkComparer(t, DefaultComparer[byte](), 1, 3, -1)

	cc := DefaultComparer[*pseudoComparable]()
	pc0 := (*pseudoComparable)(nil)
	pc1 := &pseudoComparable{name: `apple`}
	pc2 := &pseudoComparable{name: `dog`}
	checkComparer(t, cc, pc0, pc0, 0)
	checkComparer(t, cc, pc0, pc1, -1)
	checkComparer(t, cc, pc0, pc2, -1)

	checkComparer(t, cc, pc1, pc0, 1)
	checkComparer(t, cc, pc1, pc1, 0)
	checkComparer(t, cc, pc1, pc2, -1)

	checkComparer(t, cc, pc2, pc0, 1)
	checkComparer(t, cc, pc2, pc1, 1)
	checkComparer(t, cc, pc2, pc2, 0)

	checkIsNil(t, true, true, DefaultComparer[[]string]())

	checkComparer(t, DefaultComparer[time.Duration](), time.Second, time.Minute, -1)
	checkComparer(t, DefaultComparer[time.Duration](), time.Hour, time.Minute, 1)
	checkComparer(t, DefaultComparer[time.Duration](), time.Second, time.Second, 0)
	checkComparer(t, DefaultComparer[time.Duration](), time.Hour, time.Hour, 0)

	time1, err := time.Parse(time.RFC822Z, `02 Jan 24 05:30 -0700`)
	if err != nil {
		panic(err)
	}
	time2, err := time.Parse(time.RFC822Z, `02 Jan 24 05:35 -0700`)
	if err != nil {
		panic(err)
	}
	checkComparer(t, DefaultComparer[time.Time](), time1, time2, -1)
	checkComparer(t, DefaultComparer[time.Time](), time1, time1, 0)
	checkComparer(t, DefaultComparer[time.Time](), time2, time2, 0)
	checkComparer(t, DefaultComparer[time.Time](), time2, time1, 1)
}

func Test_Utils_Keys(t *testing.T) {
	e1, e2, e3, e4 := 12, 34, 56, 78
	m1 := map[*int]string{&e1: `e1`, &e2: `e2`, &e3: `e3`, &e4: `e4`}
	keys := Keys(m1)
	parts := make([]string, len(keys))
	for i, key := range keys {
		parts[i] = m1[key]
	}
	slices.Sort(parts)
	result := strings.Join(parts, `, `)

	exp := `e1, e2, e3, e4`
	if result != exp {
		t.Errorf("\n"+
			"Unexpected value from Keys:\n"+
			"\tActual:   %s\n"+
			"\tExpected: %s\n", result, exp)
	}
}

func Test_Utils_Values(t *testing.T) {
	e1, e2, e3, e4 := 12, 34, 56, 78
	m1 := map[*int]string{&e1: `e1`, &e2: `e2`, &e3: `e3`, &e4: `e4`}
	values := Values(m1)
	parts := Strings(values)
	slices.Sort(parts)
	result := strings.Join(parts, `, `)

	exp := `e1, e2, e3, e4`
	if result != exp {
		t.Errorf("\n"+
			"Unexpected value from Values:\n"+
			"\tActual:   %s\n"+
			"\tExpected: %s\n", result, exp)
	}
}

func checkTypeOf[T any](t *testing.T, exp string) {
	r := TypeOf[T]()
	s := r.String()
	if exp != s {
		t.Errorf("\n"+
			"Unexpected value from TypeOf:\n"+
			"\tActual:   %s\n"+
			"\tExpected: %s\n", s, exp)
	}
}

func Test_Utils_TypeOf(t *testing.T) {
	checkTypeOf[int](t, `int`)
	checkTypeOf[*int](t, `*int`)
	checkTypeOf[**int](t, `**int`)
	checkTypeOf[error](t, `error`)
	checkTypeOf[any](t, `interface {}`)
	checkTypeOf[interface{ foo() float64 }](t, `interface { utils.foo() float64 }`)
}

func checkString[T any](t *testing.T, exp string, value T) {
	s := String(value)
	if exp != s {
		t.Errorf("\n"+
			"Unexpected value from String:\n"+
			"\tType:          %T\n"+
			"\tActual Value:  %v\n"+
			"\tActual String: %s\n"+
			"\tExpected:      %s\n", value, value, s, exp)
	}
}

type pseudoStringer struct{ text string }

func (ps pseudoStringer) String() string { return ps.text }

func Test_Utils_String(t *testing.T) {
	checkString(t, `12`, 12)
	checkString(t, `0.005612`, 56.12e-4)
	checkString(t, `<nil>`, (*int)(nil))
	checkString(t, `Cat`, `Cat`)
	checkString(t, `Foo`, fmt.Errorf(`Foo`))
	checkString(t, `Panda`, pseudoStringer{text: `Panda`})
	checkString[any](t, `<nil>`, nil)

	var v1 map[string]int
	checkString(t, `map[]`, v1)
	v1 = map[string]int{`A`: 1, `B`: 2, `C`: 3, `D`: 4}
	checkString(t, `map[A:1 B:2 C:3 D:4]`, v1)
}

func Test_Utils_Strings(t *testing.T) {
	actual := strings.Join(Strings([]int{1, 3, 4}), `|`)
	exp := `1|3|4`
	if actual != exp {
		t.Errorf("\n"+
			"Unexpected value from Strings:\n"+
			"\tActual:   %s\n"+
			"\tExpected: %s\n", actual, exp)
	}
}

func checkGetMaxStringLen(t *testing.T, exp int, str ...string) {
	len := GetMaxStringLen(str)
	if exp != len {
		t.Errorf("\n"+
			"Unexpected value from GetMaxStringLen:\n"+
			"\tActual:   %d\n"+
			"\tExpected: %d\n", len, exp)
	}
}

func Test_Utils_GetMaxStringLen(t *testing.T) {
	checkGetMaxStringLen(t, 5, `cat`, `world`, `four`)
	checkGetMaxStringLen(t, 0)
	checkGetMaxStringLen(t, 0, ``)
}

type pseudoEquatable struct{ success bool }

func (pe *pseudoEquatable) Equals(other any) bool { return pe.success }

func checkEqual(t *testing.T, a, b any, exp bool) {
	if Equal(a, b) != exp {
		t.Errorf("\n"+
			"Unexpected value from Equal:\n"+
			"\tValue 1:  %v (%T)\n"+
			"\tValue 2:  %v (%T)\n"+
			"\tExpected: %t\n", a, a, b, b, exp)
	}
}

func Test_Utils_Equal(t *testing.T) {
	checkEqual(t, true, true, true)
	checkEqual(t, false, true, false)
	checkEqual(t, true, false, false)
	checkEqual(t, false, false, true)

	checkEqual(t, 1, 1, true)
	checkEqual(t, 1, 2, false)
	checkEqual(t, 2, 1, false)
	checkEqual(t, 2, 2, true)

	e0 := (error)(nil)
	e1 := errors.New(`fred`)
	checkEqual(t, nil, nil, true)
	checkEqual(t, e0, e0, true)
	checkEqual(t, e1, e0, false)
	checkEqual(t, e0, e1, false)
	checkEqual(t, e1, e1, true)

	var v1 int = 0
	var v2 float64 = 0.0
	checkEqual(t, nil, v1, false)
	checkEqual(t, v2, v1, false)
	checkEqual(t, v1, v2, false)

	e2 := &pseudoEquatable{success: false}
	checkEqual(t, e2, nil, false)
	checkEqual(t, nil, e2, false)
	checkEqual(t, e2, e2, false)
	checkEqual(t, 4, e2, false)
	checkEqual(t, e2, 4, false)

	e3 := &pseudoEquatable{success: true}
	checkEqual(t, e3, nil, false)
	checkEqual(t, nil, e3, false)
	checkEqual(t, e3, e3, true)
	checkEqual(t, 4, e3, true)
	checkEqual(t, e3, 4, true)

	e4 := (*pseudoEquatable)(nil)
	checkEqual(t, e4, nil, false)
	checkEqual(t, nil, e4, false)
	checkEqual(t, e4, e0, false)
	checkEqual(t, e0, e4, false)
	checkEqual(t, e4, e4, true)

	e5 := ([]int)(nil)
	e6 := []int{}
	e7 := []int{1, 2, 3}
	e8 := []int{1, 4, 3}
	checkEqual(t, e5, nil, false)
	checkEqual(t, nil, e5, false)
	checkEqual(t, e5, e5, true)
	checkEqual(t, e5, e6, false)
	checkEqual(t, e5, e7, false)
	checkEqual(t, e5, e8, false)
	checkEqual(t, e6, e6, true)
	checkEqual(t, e6, e7, false)
	checkEqual(t, e6, e8, false)
	checkEqual(t, e7, e7, true)
	checkEqual(t, e7, e8, false)
	checkEqual(t, e8, e8, true)

	e9 := func() { print(`boom`) }
	var e10 func()
	checkEqual(t, e9, e9, false)
	checkEqual(t, e9, e10, false)
	checkEqual(t, e10, e9, false)
	checkEqual(t, e10, e10, true)
	checkEqual(t, e10, nil, false)
	checkEqual(t, nil, e10, false)
}
