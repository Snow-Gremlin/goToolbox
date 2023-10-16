package predicate

import (
	"errors"
	"fmt"
	"testing"

	"goToolbox/collections"
	"goToolbox/utils"
)

func Test_Predicate_IsNil(t *testing.T) {
	p1 := IsNil[*int]()
	checkPred(t, p1, nil, true)
	v := 42
	checkPred(t, p1, &v, false)

	p2 := IsNil[int]()
	checkPred(t, p2, 0, false)
	checkPred(t, p2, 1, false)
}

func Test_Predicate_IsNotNil(t *testing.T) {
	p1 := IsNotNil[*int]()
	checkPred(t, p1, nil, false)
	v := 42
	checkPred(t, p1, &v, true)

	p2 := IsNotNil[int]()
	checkPred(t, p2, 0, true)
	checkPred(t, p2, 1, true)
}

func Test_Predicate_IsZero(t *testing.T) {
	p1 := IsZero[int]()
	checkPred(t, p1, -1, false)
	checkPred(t, p1, 0, true)
	checkPred(t, p1, 1, false)

	p2 := IsZero[*int]()
	v1, v2 := 0, 42
	checkPred(t, p2, nil, true)
	checkPred(t, p2, &v1, false)
	checkPred(t, p2, &v2, false)

	p3 := IsZero[error]()
	checkPred(t, p3, nil, true)
	checkPred(t, p3, errors.New(`Boom`), false)

	type catLady struct{ cats int }
	p4 := IsZero[catLady]()
	v3 := catLady{cats: 0}
	v4 := catLady{cats: 36}
	checkPred(t, p4, v3, true)
	checkPred(t, p4, v4, false)
}

func Test_Predicate_IsNotZero(t *testing.T) {
	p1 := IsNotZero[int]()
	checkPred(t, p1, -1, true)
	checkPred(t, p1, 0, false)
	checkPred(t, p1, 1, true)

	p2 := IsNotZero[*int]()
	v1, v2 := 0, 42
	checkPred(t, p2, nil, false)
	checkPred(t, p2, &v1, true)
	checkPred(t, p2, &v2, true)

	p3 := IsNotZero[error]()
	checkPred(t, p3, nil, false)
	checkPred(t, p3, errors.New(`Boom`), true)

	type catLady struct{ cats int }
	p4 := IsNotZero[catLady]()
	v3 := catLady{cats: 0}
	v4 := catLady{cats: 36}
	checkPred(t, p4, v3, false)
	checkPred(t, p4, v4, true)
}

func Test_Predicate_IsTrue(t *testing.T) {
	p1 := IsTrue()
	checkPred(t, p1, false, false)
	checkPred(t, p1, true, true)
}

func Test_Predicate_IsFalse(t *testing.T) {
	p1 := IsFalse()
	checkPred(t, p1, true, false)
	checkPred(t, p1, false, true)
}

func Test_Predicate_OfType(t *testing.T) {
	p := OfType[interface{ Unwrap() error }, error]()
	checkPred(t, p, nil, false)
	e1 := errors.New(`oops`)
	checkPred(t, p, e1, false)
	e2 := fmt.Errorf(`Now: %w`, e1)
	checkPred(t, p, e2, true)
	var e3 interface {
		error
		Error() string
	}
	checkPred(t, p, e3, false)
	e4 := fmt.Errorf(`Both: %w and %w`, e1, e1)
	checkPred(t, p, e4, false)
}

func Test_Predicate_In(t *testing.T) {
	p := In(1, 3, 5, 7)
	checkPred(t, p, 0, false)
	checkPred(t, p, 1, true)
	checkPred(t, p, 2, false)
	checkPred(t, p, 3, true)
	checkPred(t, p, 4, false)
	checkPred(t, p, 5, true)
}

func Test_Predicate_AsString(t *testing.T) {
	p := AsString(Matches(`^\d(\d)?$`))
	checkPred(t, p, 2, true)
	checkPred(t, p, 12, true)
	checkPred(t, p, 120, false)
	checkPred(t, p, `a`, false)
	checkPred(t, p, `1`, true)
	checkPred(t, p, `42`, true)
	checkPred(t, p, `340`, false)
}

func Test_Predicate_Matches(t *testing.T) {
	p := Matches(`\W\w{2}\W`)
	checkPred(t, p, ``, false)
	checkPred(t, p, `Hello`, false)
	checkPred(t, p, `It can`, false)
	checkPred(t, p, `It is me`, true)
	checkPred(t, p, `It...is...me`, true)
	checkPred(t, p, `It was me`, false)

	checkPanic(t, func() {
		Matches(``)
	}, `may not used an empty pattern string`)

	checkPanic(t, func() {
		Matches(`:)`)
	}, `invalid regular expression pattern {pattern: :)}: `+
		`error parsing regexp: unexpected ): `+"`:)`")
}

func Test_Predicate_Eq(t *testing.T) {
	v1, v2, v3 := 1, 2, 3
	p1 := Eq(&v1)
	checkPred(t, p1, &v1, true)
	checkPred(t, p1, &v2, false)
	checkPred(t, p1, &v3, false)

	p2 := Eq(5)
	checkPred(t, p2, 4, false)
	checkPred(t, p2, 5, true)
	checkPred(t, p2, 6, false)
}

func Test_Predicate_NotEq(t *testing.T) {
	v1, v2, v3 := 1, 2, 3
	p1 := NotEq(&v1)
	checkPred(t, p1, &v1, false)
	checkPred(t, p1, &v2, true)
	checkPred(t, p1, &v3, true)

	p2 := NotEq(5)
	checkPred(t, p2, 4, true)
	checkPred(t, p2, 5, false)
	checkPred(t, p2, 6, true)
}

func Test_Predicate_GreaterThan(t *testing.T) {
	p := GreaterThan(5)
	checkPred(t, p, 4, false)
	checkPred(t, p, 5, false)
	checkPred(t, p, 6, true)
}

func Test_Predicate_GreaterEq(t *testing.T) {
	p := GreaterEq(5)
	checkPred(t, p, 4, false)
	checkPred(t, p, 5, true)
	checkPred(t, p, 6, true)
}

func Test_Predicate_LessThan(t *testing.T) {
	p := LessThan(5)
	checkPred(t, p, 4, true)
	checkPred(t, p, 5, false)
	checkPred(t, p, 6, false)
}

func Test_Predicate_LessEq(t *testing.T) {
	p := LessEq(5)
	checkPred(t, p, 4, true)
	checkPred(t, p, 5, true)
	checkPred(t, p, 6, false)
}

func Test_Predicate_InRange(t *testing.T) {
	p := InRange(5, 7)
	checkPred(t, p, 4, false)
	checkPred(t, p, 5, true)
	checkPred(t, p, 6, true)
	checkPred(t, p, 7, true)
	checkPred(t, p, 8, false)
}

func Test_Predicate_EpsilonEq(t *testing.T) {
	p1 := EpsilonEq(5.0, .2)
	checkPred(t, p1, 4.0, false)
	checkPred(t, p1, 4.7, false)
	checkPred(t, p1, 4.9, true)
	checkPred(t, p1, 5.0, true)
	checkPred(t, p1, 5.1, true)
	checkPred(t, p1, 5.3, false)
	checkPred(t, p1, 6.0, false)

	p2 := EpsilonEq(5, 1)
	checkPred(t, p2, 3, false)
	checkPred(t, p2, 4, true)
	checkPred(t, p2, 5, true)
	checkPred(t, p2, 6, true)
	checkPred(t, p2, 7, false)
}

func Test_Predicate_EpsilonNotEq(t *testing.T) {
	p1 := EpsilonNotEq(5.0, .2)
	checkPred(t, p1, 4.0, true)
	checkPred(t, p1, 4.7, true)
	checkPred(t, p1, 4.9, false)
	checkPred(t, p1, 5.0, false)
	checkPred(t, p1, 5.1, false)
	checkPred(t, p1, 5.3, true)
	checkPred(t, p1, 6.0, true)

	p2 := EpsilonNotEq(5, 1)
	checkPred(t, p2, 3, true)
	checkPred(t, p2, 4, false)
	checkPred(t, p2, 5, false)
	checkPred(t, p2, 6, false)
	checkPred(t, p2, 7, true)
}

func Test_Predicate_Pos(t *testing.T) {
	p := Pos[float64]()
	checkPred(t, p, 1.0e9, true)
	checkPred(t, p, 1.0, true)
	checkPred(t, p, 1.0e-9, true)
	checkPred(t, p, 0.0, false)
	checkPred(t, p, -1.0e-9, false)
	checkPred(t, p, -1.0, false)
	checkPred(t, p, -1.0e9, false)
}

func Test_Predicate_Neg(t *testing.T) {
	p := Neg[float64]()
	checkPred(t, p, 1.0e9, false)
	checkPred(t, p, 1.0, false)
	checkPred(t, p, 1.0e-9, false)
	checkPred(t, p, 0.0, false)
	checkPred(t, p, -1.0e-9, true)
	checkPred(t, p, -1.0, true)
	checkPred(t, p, -1.0e9, true)
}

func Test_Predicate_Not(t *testing.T) {
	p := Not(Eq(5))
	checkPred(t, p, 4, true)
	checkPred(t, p, 5, false)
	checkPred(t, p, 6, true)
}

func Test_Predicate_And(t *testing.T) {
	p := And(LessThan(7), GreaterThan(4))
	checkPred(t, p, 3, false)
	checkPred(t, p, 4, false)
	checkPred(t, p, 5, true)
	checkPred(t, p, 6, true)
	checkPred(t, p, 7, false)
	checkPred(t, p, 8, false)
}

func Test_Predicate_Or(t *testing.T) {
	p := Or(LessEq(4), GreaterEq(7))
	checkPred(t, p, 3, true)
	checkPred(t, p, 4, true)
	checkPred(t, p, 5, false)
	checkPred(t, p, 6, false)
	checkPred(t, p, 7, true)
	checkPred(t, p, 8, true)
}

func Test_Predicate_OnlyOne(t *testing.T) {
	p := OnlyOne(LessEq(4), LessEq(6), GreaterEq(8), GreaterEq(10))
	checkPred(t, p, 3, false)
	checkPred(t, p, 4, false)
	checkPred(t, p, 5, true)
	checkPred(t, p, 6, true)
	checkPred(t, p, 7, false)
	checkPred(t, p, 8, true)
	checkPred(t, p, 9, true)
	checkPred(t, p, 10, false)
	checkPred(t, p, 11, false)
}

func checkPred[T any](t *testing.T, p collections.Predicate[T], input T, exp bool) {
	t.Helper()
	actual := p(input)
	if !utils.Equal(actual, exp) {
		t.Errorf("\nUnexpected result from predicate:\n"+
			"\tInput:    %v\n"+
			"\tActual:   %t\n"+
			"\tExpected: %t", input, actual, exp)
	}
}

func checkPanic(t *testing.T, handle func(), exp string) {
	t.Helper()
	actual := func() (r string) {
		defer func() { r = utils.String(recover()) }()
		handle()
		return ``
	}()

	if !utils.Equal(actual, exp) {
		t.Errorf("\nUnexpected panic string from predicate creation:\n"+
			"\tActual:   %s\n"+
			"\tExpected: %s", actual, exp)
	}
}
