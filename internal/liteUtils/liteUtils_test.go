package liteUtils

import (
	"errors"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"testing"
)

func checkIsZero[T any](t *testing.T, exp bool, value T) {
	if exp != IsZero(value) {
		t.Errorf("\n"+
			"Unexpected value from IsZero:\n"+
			"\tType:     %T\n"+
			"\tExpected: %t\n", value, exp)
	}

	v := reflect.ValueOf(value)
	if v.IsValid() && v.IsZero() != exp {
		t.Errorf("\n"+
			"Unexpected value from reflect.ValueOf.IsZero:\n"+
			"\tType:     %T\n"+
			"\tExpected: %t\n", value, exp)
	}
}

func Test_LiteUtils_Zero(t *testing.T) {
	checkIsZero(t, true, Zero[bool]())
	checkIsZero(t, false, true)

	checkIsZero(t, true, Zero[int]())
	checkIsZero(t, true, 0)
	checkIsZero(t, false, -1)
	checkIsZero(t, false, 42)

	checkIsZero(t, true, Zero[[]float64]())
	checkIsZero(t, true, ([]float64)(nil))
	checkIsZero(t, false, []float64{})
	checkIsZero(t, false, []float64{0.0})

	checkIsZero(t, true, Zero[string]())
	checkIsZero(t, true, ``)
	checkIsZero(t, false, `Hello`)

	checkIsZero(t, true, Zero[*int]())
	i := 0
	checkIsZero(t, false, &i)

	checkIsZero(t, true, Zero[map[string]int]())
	m := map[string]int{}
	checkIsZero(t, false, m)

	checkIsZero(t, true, Zero[interface{ Cats() int }]())

	checkIsZero(t, true, Zero[struct{ cats int }]())
	s := struct{ cats int }{cats: 9}
	checkIsZero(t, false, s)

	checkIsZero(t, true, Zero[func(...int) float64]())
	checkIsZero(t, true, Zero[func(*testing.T)]())
	checkIsZero(t, false, Test_LiteUtils_Zero)

	checkIsZero(t, true, Zero[chan *string]())
	c := make(chan *string)
	checkIsZero(t, false, c)

	checkIsZero(t, true, Zero[**string]())
	checkIsZero(t, true, Zero[any]())
	checkIsZero(t, true, Zero[error]())

	checkIsZero(t, true, Zero[testing.T]())
	checkIsZero(t, false, t)
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

func Test_LiteUtils_IsNil(t *testing.T) {
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

func Test_LiteUtils_String(t *testing.T) {
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

func Test_LiteUtils_Strings(t *testing.T) {
	actual := strings.Join(Strings([]int{1, 3, 4}), `|`)
	exp := `1|3|4`
	if actual != exp {
		t.Errorf("\n"+
			"Unexpected value from Strings:\n"+
			"\tActual:   %s\n"+
			"\tExpected: %s\n", actual, exp)
	}
}

type pseudoEquatable struct{ success bool }

func (pe *pseudoEquatable) Equals(_ any) bool { return pe.success }

func checkEqual(t *testing.T, a, b any, exp bool) {
	if Equal(a, b) != exp {
		t.Errorf("\n"+
			"Unexpected value from Equal:\n"+
			"\tValue 1:  %v (%T)\n"+
			"\tValue 2:  %v (%T)\n"+
			"\tExpected: %t\n", a, a, b, b, exp)
	}
}

func Test_LiteUtils_Equal(t *testing.T) {
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

func Test_LiteUtils_Keys(t *testing.T) {
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

func Test_LiteUtils_Values(t *testing.T) {
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
