package utils

import (
	"errors"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/Snow-Gremlin/goToolbox/comp"
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
	v3 = &strconv.NumError{
		Func: `Oops`,
		Num:  `X`,
		Err:  nil,
	}
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
	length, ok := Length(value)
	if expLen != length || expOk != ok {
		t.Errorf("\n"+
			"Unexpected value from Length:\n"+
			"\tType:     %T\n"+
			"\tActual:   %d, %t\n"+
			"\tExpected: %d, %t\n", value, length, ok, expLen, expOk)
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

func checkSortedKeys[TKey comparable, TValue any, TMap ~map[TKey]TValue](t *testing.T, value TMap, exp string, cmp ...comp.Comparer[TKey]) {
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
		map[string]float64{`cat`: 4.3, `pig`: 2.16, `Dog`: 333.333, `apple`: 12.34},
		`Dog, apple, cat, pig`)

	checkSortedKeys(t,
		map[int]float64{3: 4.3, 1: 2.16, 6: 333.333, 5: 12.34},
		`6, 5, 3, 1`, comp.Descender(comp.Ordered[int]()))

	checkSortedKeys(t,
		map[float64]int{4.3: 3, 2.16: 1, 333.333: 6, 12.34: 5},
		`333.333, 12.34, 4.3, 2.16`, comp.Descender(comp.Ordered[float64]()))

	checkSortedKeys(t,
		map[rune]float64{'k': 4.3, 'a': 2.16, 'q': 333.333, 'H': 12.34},
		`113, 107, 97, 72`, comp.Descender(comp.Ordered[rune]()))

	checkSortedKeys(t,
		map[string]float64{`cat`: 4.3, `pig`: 2.16, `Dog`: 333.333, `apple`: 12.34},
		`pig, cat, apple, Dog`, comp.Descender(comp.Ordered[string]()))

	errStr := catchPanic(func() {
		SortedKeys(map[string]float64{`apple`: 12.34}, comp.Ordered[string](), comp.Ordered[string]())
	})
	checkEqual(t, `invalid number of arguments {count: 2, maximum: 1, usage: comparer}`, errStr, true)

	errStr = catchPanic(func() {
		SortedKeys(map[complex128]float64{}, nil)
	})
	checkEqual(t, `must provide a comparer to compare this type {type: complex128}`, errStr, true)
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
	checkString(t, `Foo`, errors.New(`Foo`))
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

func Test_Utils_LazyMatcher(t *testing.T) {
	hex := LazyMatcher(`^[0-9A-Fa-f]+$`)
	checkEqual(t, hex(`572A6F`), true, true)
	checkEqual(t, hex(`CAT`), false, true)

	exp := "regexp: Compile(`((((`): error parsing regexp: missing closing ): `((((`"
	hex = LazyMatcher(`((((`) // bad pattern but won't panic yet

	err := func() (r any) {
		defer func() { r = recover() }()
		// panic occurs on first use
		return hex(`CAT`)
	}()
	checkEqual(t, err, exp, true)

	err = func() (r any) {
		defer func() { r = recover() }()
		// panic is re-panicked on each following call
		return hex(`CAT`)
	}()
	checkEqual(t, err, exp, true)
}

func Test_Utils_Parse(t *testing.T) {
	v1, err := Parse[string](`Cat`)
	checkEqual(t, v1, `Cat`, true)
	checkEqual(t, err, nil, true)

	v2, err := Parse[bool](`true`)
	checkEqual(t, v2, true, true)
	checkEqual(t, err, nil, true)

	v2, err = Parse[bool](`Cat`)
	checkEqual(t, v2, false, true)
	checkEqual(t, err.Error(),
		`unable to parse value {input: Cat, type: bool}: `+
			`strconv.ParseBool: parsing "Cat": `+
			`invalid syntax: invalid syntax`, true)

	v3, err := Parse[int](`-24`)
	checkEqual(t, v3, -24, true)
	checkEqual(t, err, nil, true)

	v3, err = Parse[int](`-0xA5`)
	checkEqual(t, v3, -0xA5, true)
	checkEqual(t, err, nil, true)

	v4, err := Parse[int8](`-102`)
	checkEqual(t, v4, int8(-102), true)
	checkEqual(t, err, nil, true)

	v5, err := Parse[int16](`-2455`)
	checkEqual(t, v5, int16(-2455), true)
	checkEqual(t, err, nil, true)

	v6, err := Parse[int32](`-245512`)
	checkEqual(t, v6, int32(-245512), true)
	checkEqual(t, err, nil, true)

	v7, err := Parse[int64](`-24533512`)
	checkEqual(t, v7, int64(-24533512), true)
	checkEqual(t, err, nil, true)

	v8, err := Parse[uint](`42`)
	checkEqual(t, v8, uint(42), true)
	checkEqual(t, err, nil, true)

	v9, err := Parse[uint8](`102`)
	checkEqual(t, v9, uint8(102), true)
	checkEqual(t, err, nil, true)

	v10, err := Parse[uint16](`2455`)
	checkEqual(t, v10, uint16(2455), true)
	checkEqual(t, err, nil, true)

	v11, err := Parse[uint32](`245512`)
	checkEqual(t, v11, uint32(245512), true)
	checkEqual(t, err, nil, true)

	v12, err := Parse[uint64](`24533512`)
	checkEqual(t, v12, uint64(24533512), true)
	checkEqual(t, err, nil, true)

	v13, err := Parse[float32](`2453.3512`)
	checkEqual(t, v13, float32(2453.3512), true)
	checkEqual(t, err, nil, true)

	v14, err := Parse[float64](`-2458.351`)
	checkEqual(t, v14, float64(-2458.351), true)
	checkEqual(t, err, nil, true)

	v15, err := Parse[complex64](`234+41i`)
	checkEqual(t, v15, complex64(234+41i), true)
	checkEqual(t, err, nil, true)

	v16, err := Parse[complex128](`-1234.3+41.4i`)
	checkEqual(t, v16, -1234.3+41.4i, true)
	checkEqual(t, err, nil, true)
}

func checkGetMaxStringLen(t *testing.T, exp int, str ...string) {
	length := GetMaxStringLen(str)
	if exp != length {
		t.Errorf("\n"+
			"Unexpected value from GetMaxStringLen:\n"+
			"\tActual:   %d\n"+
			"\tExpected: %d\n", length, exp)
	}
}

func Test_Utils_GetMaxStringLen(t *testing.T) {
	checkGetMaxStringLen(t, 5, `cat`, `world`, `four`)
	checkGetMaxStringLen(t, 0)
	checkGetMaxStringLen(t, 0, ``)
}

func checkEqual(t *testing.T, a, b any, exp bool) {
	if comp.Equal(a, b) != exp {
		t.Errorf("\n"+
			"Unexpected value from Equal:\n"+
			"\tValue 1:  %v (%T)\n"+
			"\tValue 2:  %v (%T)\n"+
			"\tExpected: %t\n", a, a, b, b, exp)
	}
}

func Test_Utils_Ternary(t *testing.T) {
	checkEqual(t, Ternary(true, 12, 34), 12, true)
	checkEqual(t, Ternary(false, 12, 34), 34, true)
}

func Test_Utils_Flip(t *testing.T) {
	a, b := Flip(true, 12, 34)
	checkEqual(t, a, 34, true)
	checkEqual(t, b, 12, true)

	a, b = Flip(false, 12, 34)
	checkEqual(t, a, 12, true)
	checkEqual(t, b, 34, true)
}
