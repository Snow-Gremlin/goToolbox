package check

import (
	"errors"
	"fmt"
	"maps"
	"reflect"
	"sort"
	"strings"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/collections/predicate"
	"github.com/Snow-Gremlin/goToolbox/collections/readonlyVariantList"
	"github.com/Snow-Gremlin/goToolbox/comp"
	"github.com/Snow-Gremlin/goToolbox/internal/optional"
	"github.com/Snow-Gremlin/goToolbox/testers"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

// Nil creates a check that the actual type is nil.
//
// The actual value must a type which can be nil.
//
// Example: check.Nil(t).Assert(actual)
func Nil(t testers.Tester) (c testers.Check[any]) {
	defer handlePanic(t, &c)
	getHelper(t)()
	return newCheck(t, func(b *testee, actual any) {
		if isNil, ok := utils.TryIsNil(actual); !ok {
			b.Should(`be a nil-able type`)
		} else if !isNil {
			b.Should(`be nil`)
		}
	})
}

// NotNil creates a check that the actual type is not nil.
//
// The actual value must a type which can be nil.
//
// Example: check.NotNil(t).Assert(actual)
func NotNil(t testers.Tester) (c testers.Check[any]) {
	defer handlePanic(t, &c)
	getHelper(t)()
	return newCheck(t, func(b *testee, actual any) {
		if isNil, ok := utils.TryIsNil(actual); !ok {
			b.Should(`be a nil-able type`)
		} else if isNil {
			b.Should(`not be nil`)
		}
	})
}

// Zero creates a check that the actual value
// is the zero value of that type.
//
// Example: check.Zero(t).Assert(actual)
func Zero(t testers.Tester) (c testers.Check[any]) {
	defer handlePanic(t, &c)
	getHelper(t)()
	return newPred(t, predicate.IsZero[any](), `be a zero value`)
}

// NotZero creates a check that the actual value
// is not the zero value of that type.
//
// Example: check.NotZero(t).Assert(actual)
func NotZero(t testers.Tester) (c testers.Check[any]) {
	defer handlePanic(t, &c)
	getHelper(t)()
	return newPred(t, predicate.IsNotZero[any](), `not be a zero value`)
}

// True creates a check that the actual boolean value is true.
//
// Example: check.True(t).Assert(actual)
func True(t testers.Tester) (c testers.Check[bool]) {
	defer handlePanic(t, &c)
	getHelper(t)()
	return newPred(t, predicate.IsTrue(), `be true`)
}

// False creates a check that the actual boolean value is false.
//
// Example: check.False(t).Assert(actual)
func False(t testers.Tester) (c testers.Check[bool]) {
	defer handlePanic(t, &c)
	getHelper(t)()
	return newPred(t, predicate.IsFalse(), `be false`)
}

// Type creates a check that the actual value is the given type.
//
// Example: check.Type[int](t).Assert(actual)
func Type[T any](t testers.Tester) (c testers.Check[any]) {
	defer handlePanic(t, &c)
	getHelper(t)()
	return newPred[any](t, predicate.OfType[T, any](), `be the expected type`).
		With(`Expected Type`, utils.TypeOf[T]())
}

// NotType creates a check that the actual value is not the given type.
//
// Example: check.Type[int](t).Assert(actual)
func NotType[T any](t testers.Tester) (c testers.Check[any]) {
	defer handlePanic(t, &c)
	getHelper(t)()
	return newPred[any](t, predicate.Not(predicate.OfType[T, any]()),
		`not be the unexpected type`).
		With(`Unexpected Type`, utils.TypeOf[T]())
}

// Match creates a check that the given expected regular expression
// matches the actual string.
//
// The given pattern must be a string or a `*regexp.Regexp` instance.
// This uses `utils.String` to get the string from the actual value.
//
// Example: `check.Match(t, "^\w+$").Assert(actual)`
func Match(t testers.Tester, regex any) (c testers.Check[any]) {
	defer handlePanic(t, &c)
	getHelper(t)()
	re := getRegExp(t, regex)
	if re == nil {
		return (*checkImp[any])(nil)
	}
	return newPred(t, predicate.AsString(re.MatchString),
		`match the given regular expression pattern`).
		With(`Pattern`, re.String())
}

// String creates a check that the given expected string
// is the same of the string from the actual object.
//
// This uses `utils.String` to get the string.
//
// Example: `check.String(t, "foo").Assert(actual)`
func String(t testers.Tester, expected string) (c testers.Check[any]) {
	defer handlePanic(t, &c)
	getHelper(t)()
	return newPred(t, predicate.AsString(predicate.Eq(expected)),
		`have string be equal`).
		AsText().
		WithValue(`Expected String`, expected)
}

// StringAndReset creates a check that the given expected string
// is the same of the string from the actual object then calls
// `Reset` on the actual object if it exists.
//
// This uses `utils.String` to get the string.
// This is designed to work with `bytes.Buffer` so that a buffer can be used
// to collect changes, event, calls, etc then checked that those occurred.
// Resetting the buffer prepares it for the next collection of changes, etc.
//
// Example: `check.StringAndReset(t, "foo").Assert(actual)`
func StringAndReset(t testers.Tester, expected string) (c testers.Check[any]) {
	defer handlePanic(t, &c)
	getHelper(t)()
	p := predicate.AsString(predicate.Eq(expected))
	return newCheck(t, func(b *testee, actual any) {
		if !p(actual) {
			b.Should(`have string be equal`)
		}
		if r, ok := actual.(interface{ Reset() }); ok {
			// Overwrite value with the string so the reset
			// doesn't change the formatted string of the value.
			b.With(`Actual Value`, b.formatValue(r))
			r.Reset()
		}
	}).AsText().
		WithValue(`Expected String`, expected)
}

// Equal creates a check that the given expected value
// is equal to an actual value.
//
// This uses `utils.Equal` for the comparison.
//
// Example: `check.Equal(t, 42).Assert(actual)`
func Equal[T any](t testers.Tester, expected T) (c testers.Check[T]) {
	defer handlePanic(t, &c)
	getHelper(t)()
	return newPred(t, predicate.Eq(expected), `be equal`).
		WithValue(`Expected Value`, expected).
		WithType(`Expected Type`, expected)
}

// NotEqual creates a check that the given expected value
// is not equal to an actual value.
//
// This uses `utils.Equal` for the comparison.
//
// Example: `check.NotEqual(t, 42).Assert(actual)`
func NotEqual[T any](t testers.Tester, unexpected T) (c testers.Check[T]) {
	defer handlePanic(t, &c)
	getHelper(t)()
	return newPred(t, predicate.NotEq(unexpected), `not be equal`).
		WithValue(`Unexpected Value`, unexpected).
		WithType(`Unexpected Type`, unexpected)
}

// GreaterThan creates a check that the actual value
// is greater than the given expected value.
//
// May provide one optional comparer to use for this check.
// If the type being checked does not have a default comparer,
// then a comparer for that type must be provided.
//
// Example: check.GreaterThan(t, 14).Assert(actual)
func GreaterThan[T any](t testers.Tester, expected T, comparer ...comp.Comparer[T]) (c testers.Check[T]) {
	defer handlePanic(t, &c)
	getHelper(t)()
	return newPred(t, predicate.GreaterThan(expected, comparer...),
		`be greater than the expected value`).
		WithValue(`Minimum Value`, expected).
		WithType(`Minimum Type`, expected)
}

// GreaterEq creates a check that the actual value
// is greater than or equal to the given expected value.
//
// May provide one optional comparer to use for this check.
// If the type being checked does not have a default comparer,
// then a comparer for that type must be provided.
//
// Example: check.GreaterEq(t, 14).Assert(actual)
func GreaterEq[T any](t testers.Tester, expected T, comparer ...comp.Comparer[T]) (c testers.Check[T]) {
	defer handlePanic(t, &c)
	getHelper(t)()
	return newPred(t, predicate.GreaterEq(expected, comparer...),
		`be greater than or equal to the expected value`).
		WithValue(`Minimum Value`, expected).
		WithType(`Minimum Type`, expected)
}

// LessThan creates a check that the actual value
// is less than the given expected value.
//
// May provide one optional comparer to use for this check.
// If the type being checked does not have a default comparer,
// then a comparer for that type must be provided.
//
// Example: check.LessThan(t, 14).Assert(actual)
func LessThan[T any](t testers.Tester, expected T, comparer ...comp.Comparer[T]) (c testers.Check[T]) {
	defer handlePanic(t, &c)
	getHelper(t)()
	return newPred(t, predicate.LessThan(expected, comparer...),
		`be less than the expected value`).
		WithValue(`Maximum Value`, expected).
		WithType(`Maximum Type`, expected)
}

// LessEq creates a check that the actual value
// is less than or equal to the given expected value.
//
// May provide one optional comparer to use for this check.
// If the type being checked does not have a default comparer,
// then a comparer for that type must be provided.
//
// Example: check.LessEq(t, 14).Assert(actual)
func LessEq[T any](t testers.Tester, expected T, comparer ...comp.Comparer[T]) (c testers.Check[T]) {
	defer handlePanic(t, &c)
	getHelper(t)()
	return newPred(t, predicate.LessEq(expected, comparer...),
		`be less than or equal to the expected value`).
		WithValue(`Maximum Value`, expected).
		WithType(`Maximum Type`, expected)
}

// InRange creates a check that the actual value
// is between the given min and maximum inclusively.
//
// May provide one optional comparer to use for this check.
// If the type being checked does not have a default comparer,
// then a comparer for that type must be provided.
//
// Example: InRange.LessEq(t, 0, 359).Assert(actual)
func InRange[T any](t testers.Tester, min, max T, comparer ...comp.Comparer[T]) (c testers.Check[T]) {
	defer handlePanic(t, &c)
	getHelper(t)()
	return newPred(t, predicate.InRange(min, max, comparer...),
		`be between or equal to the given maximum and minimum`).
		WithValue(`Minimum Value`, min).
		WithValue(`Maximum Value`, max).
		WithType(`Range Type`, max)
}

// Epsilon creates a check that the actual value is equal to the given expected value
// or very close to the given expected value. The value must be within the given epsilon
// to be considered equal.
//
// The given epsilon must be greater than zero. An epsilon comparator should be used
// when comparing calculated floating point numbers since calculations may accrue small
// errors and make the actual value very close to the expected value but not exactly equal.
//
// Example: check.Epsilon(t, 14.0, 1.0e-9).Assert(actual)
func Epsilon[T utils.NumConstraint](t testers.Tester, expected, epsilon T) (c testers.Check[T]) {
	defer handlePanic(t, &c)
	getHelper(t)()
	if epsilon <= 0 {
		newTestee(t).WithValue(`Epsilon Value`, epsilon).
			WithType(`Epsilon Type`, epsilon).
			SetupMust(`provide an epsilon greater than zero`)
		return (*checkImp[T])(nil)
	}
	p := predicate.EpsilonEq(expected, epsilon)
	return newCheck(t, func(b *testee, actual T) {
		if !p(actual) {
			b.Should(`be within the epsilon of the expected value`).
				WithValue(`Expected Value`, expected).
				WithType(`Expected Type`, expected).
				WithValue(`Epsilon`, epsilon).
				WithValue(`Delta`, expected-actual)
		}
	})
}

// NotEpsilon creates a check that the actual value is not equal to the given
// unexpected value and not very close to the given unexpected value. The value
// must not be within the given epsilon to be considered not equal.
//
// The given epsilon must be greater than zero. An epsilon comparator should be used
// when comparing calculated floating point numbers since calculations may accrue small
// errors and make the actual value very close to the expected value but not exactly equal.
//
// Example: check.NotEpsilon(t, 14.0, 1.0e-9).Assert(actual)
func NotEpsilon[T utils.NumConstraint](t testers.Tester, unexpected, epsilon T) (c testers.Check[T]) {
	defer handlePanic(t, &c)
	getHelper(t)()
	if epsilon <= 0 {
		newTestee(t).WithValue(`Epsilon Value`, epsilon).
			WithType(`Epsilon Type`, epsilon).
			SetupMust(`provide an epsilon greater than zero`)
		return (*checkImp[T])(nil)
	}
	p := predicate.EpsilonNotEq(unexpected, epsilon)
	return newCheck(t, func(b *testee, actual T) {
		if !p(actual) {
			b.Should(`not be within the epsilon of the unexpected value`).
				WithValue(`Unexpected Value`, unexpected).
				WithType(`Expected Type`, unexpected).
				WithValue(`Epsilon`, epsilon).
				WithValue(`Delta`, unexpected-actual)
		}
	})
}

// Positive creates a check that the actual value is greater than zero.
//
// May provide one optional comparer to use for this check.
// If the type being checked does not have a default comparer,
// then a comparer for that type must be provided.
//
// Example: check.Positive[int](t).Assert(actual)
func Positive[T any](t testers.Tester, comparer ...comp.Comparer[T]) (c testers.Check[T]) {
	defer handlePanic(t, &c)
	getHelper(t)()
	return newPred(t, predicate.Pos(comparer...),
		`be a positive value`).
		With(`Expected Type`, utils.TypeOf[T]())
}

// Negative creates a check that the actual value is less than zero.
//
// May provide one optional comparer to use for this check.
// If the type being checked does not have a default comparer,
// then a comparer for that type must be provided.
//
// Example: check.Negative(t).Assert(actual)
func Negative[T any](t testers.Tester, comparer ...comp.Comparer[T]) (c testers.Check[T]) {
	defer handlePanic(t, &c)
	getHelper(t)()
	return newPred(t, predicate.Neg(comparer...),
		`be a negative value`).
		With(`Expected Type`, utils.TypeOf[T]())
}

// Is creates a check that the actual value causes the given predicate to return true.
//
// Example: check.Is(t, func(x thing) bool { return thing.Valid() }).Assert(actual)
func Is[T any](t testers.Tester, p collections.Predicate[T]) (c testers.Check[T]) {
	defer handlePanic(t, &c)
	getHelper(t)()
	return newPred(t, p, `be accepted by the given predicate`)
}

// IsNot creates a check that the actual value causes the given predicate to return false.
//
// Example: check.IsNot(t, func(x thing) bool { return thing.Valid() }).Assert(actual)
func IsNot[T any](t testers.Tester, p collections.Predicate[T]) (c testers.Check[T]) {
	defer handlePanic(t, &c)
	getHelper(t)()
	return newPred(t, predicate.Not(p), `not be accepted by the given predicate`)
}

// StartsWith creates a check that the given expected
// string or array is the prefix for the actual object.
//
// Example: `check.StartsWith(t, "foo").Assert(actual)`
// Example: `check.StartsWith(t, []int{3, 4, 5}).Assert(actual)`
func StartsWith(t testers.Tester, expected any) (c testers.Check[any]) {
	defer handlePanic(t, &c)
	getHelper(t)()
	expV := readonlyVariantList.Wrap(expected)
	if expV.Count() <= 0 {
		newTestee(t).WithType(`Expected Type`, expected).
			SetupMust(`have at least one expected prefix value`)
		return (*checkImp[any])(nil)
	}

	return newCheck(t, func(b *testee, actual any) {
		if actV := readonlyVariantList.Wrap(actual); !actV.StartsWith(expV) {
			b.WithValue(`Expected Prefix`, expected).
				WithType(`Expected Type`, expected).
				Should(`start with the given prefix`)
		}
	})
}

// EndsWith creates a check that the given expected
// string or array is the suffix for the actual object.
//
// Example: `check.EndsWith(t, "foo").Assert(actual)`
// Example: `check.EndsWith(t, []int{3, 4, 5}).Assert(actual)`
func EndsWith(t testers.Tester, expected any) (c testers.Check[any]) {
	defer handlePanic(t, &c)
	getHelper(t)()
	expV := readonlyVariantList.Wrap(expected)
	if expV.Count() <= 0 {
		newTestee(t).WithType(`Expected Type`, expected).
			SetupMust(`have at least one expected suffix value`)
		return (*checkImp[any])(nil)
	}

	return newCheck(t, func(b *testee, actual any) {
		if actV := readonlyVariantList.Wrap(actual); !actV.EndsWith(expV) {
			b.WithValue(`Expected Suffix`, expected).
				WithType(`Expected Type`, expected).
				Should(`end with the given suffix`)
		}
	})
}

// Empty creates a check that the length of the actual value is zero.
//
// This requires that the actual value is a string, slice, array, map,
// or anything that has a Len, Length, or Count method.
//
// Example: check.IsEmpty(t).Assert(actual)
func Empty(t testers.Tester) (c testers.Check[any]) {
	defer handlePanic(t, &c)
	getHelper(t)()
	return newLen[any](t, predicate.LessEq(0), `be empty`)
}

// NotEmpty creates a check that the length of the actual value is not zero.
//
// This requires that the actual value is a string, slice, array, map,
// or anything that has a Len, Length, or Count method.
//
// Example: check.IsNotEmpty(t).Assert(actual)
func NotEmpty(t testers.Tester) (c testers.Check[any]) {
	defer handlePanic(t, &c)
	getHelper(t)()
	return newLen[any](t, predicate.GreaterThan(0), `not be empty`)
}

// Single creates a check that the length of the actual value is one.
//
// This requires that the actual value is a string, slice, array, map,
// or anything that has a Len, Length, or Count method.
//
// Example: check.Single(t).Assert(actual)
func Single(t testers.Tester) (c testers.Check[any]) {
	defer handlePanic(t, &c)
	getHelper(t)()
	return newLen[any](t, predicate.Eq(1), `have one and only one value`)
}

// Length creates a check that the length of the actual value
// is equal to the given expected length.
//
// This requires that the actual value is a string, slice, array, map,
// or anything that has a Len, Length, or Count method.
//
// Example: check.Length(t, 5).Assert(actual)
func Length(t testers.Tester, expected int) (c testers.Check[any]) {
	defer handlePanic(t, &c)
	getHelper(t)()
	return newLen[any](t, predicate.Eq(expected), `be the expected length`).
		With(`Expected Length`, expected)
}

// ShorterThan creates a check that the length of the actual value
// is shorter than the given expected length.
//
// This requires that the actual value is a string, slice, array, map,
// or anything that has a Len, Length, or Count method.
//
// Example: check.ShorterThan(t, 5).Assert(actual)
func ShorterThan(t testers.Tester, expected int) (c testers.Check[any]) {
	defer handlePanic(t, &c)
	getHelper(t)()
	return newLen[any](t, predicate.LessThan(expected), `be shorter than the expected length`).
		With(`Maximum Length`, expected)
}

// LongerThan creates a check that the length of the actual value
// is longer than the given expected length.
//
// This requires that the actual value is a string, slice, array, map,
// or anything that has a Len, Length, or Count method.
//
// Example: check.LongerThan(t, 5).Assert(actual)
func LongerThan(t testers.Tester, expected int) (c testers.Check[any]) {
	defer handlePanic(t, &c)
	getHelper(t)()
	return newLen[any](t, predicate.GreaterThan(expected), `be longer than the expected length`).
		With(`Minimum Length`, expected)
}

// NoError creates a check that the actual error is not nil.
//
// Example: check.NoError(t).Assert(actual)
func NoError(t testers.Tester) (c testers.Check[error]) {
	defer handlePanic(t, &c)
	getHelper(t)()
	return newPred(t, predicate.IsNil[error](), `be no error`)
}

// MatchError creates a check that the given expected regular expression
// matches the actual error's Error() string.
//
// Example: check.MatchError(t, `^\w+$`).Assert(actual)
func MatchError(t testers.Tester, regex any) (c testers.Check[error]) {
	defer handlePanic(t, &c)
	getHelper(t)()
	re := getRegExp(t, regex)
	if re == nil {
		return (*checkImp[error])(nil)
	}
	return newCheck(t, func(b *testee, actual error) {
		if utils.IsNil(actual) {
			b.Should(`not be a nil error`)
			return
		}

		actualErr := actual.Error()
		if !re.MatchString(actualErr) {
			b.With(`Pattern`, re.String()).
				Should(`have error sting match the given regular expression pattern`)
		}
	})
}

// ErrorHas creates a check that the given error type is contained
// within the error tree of the actual error.
//
// This uses errors `As` method to find the contained type.
//
// Example: check.ErrorHas[Stacked](t).Assert(actual)
func ErrorHas[T any](t testers.Tester) (c testers.Check[error]) {
	defer handlePanic(t, &c)
	getHelper(t)()
	return newPred(t, func(actual error) bool {
		var target T
		return errors.As(actual, &target)
	}, `have an error of the target type in the error tree`).
		With(`Target Type`, utils.TypeOf[T]())
}

// Implements creates a check that the actual value implements the given type.
//
// The given target type must be an interface.
//
// Example: check.Implements[Stringer](t).Assert(actual)
func Implements[T any](t testers.Tester) (c testers.Check[any]) {
	defer handlePanic(t, &c)
	getHelper(t)()
	target := utils.TypeOf[T]()
	if target.Kind() != reflect.Interface {
		newTestee(t).With(`Type`, target).
			SetupMust(`provide an interface type`)
		return (*checkImp[any])(nil)
	}

	return newPred(t, func(actual any) bool {
		return reflect.TypeOf(actual).Implements(target)
	}, `implement the target type`).
		With(`Target Type`, target)
}

// ConvertibleTo creates a check that the actual value is conversable
// to the given expected type.
//
// Example: check.ConvertibleTo[int](t).Assert(actual)
func ConvertibleTo[T any](t testers.Tester) (c testers.Check[any]) {
	defer handlePanic(t, &c)
	getHelper(t)()
	target := utils.TypeOf[T]()
	return newPred(t, func(actual any) bool {
		return reflect.TypeOf(actual).ConvertibleTo(target)
	}, `be convertible to the target type`).
		With(`Target Type`, target)
}

// SameType creates a check that the actual value is the same type as the expected type.
//
// Example: check.SameType(t, 3.2).Assert(actual)
func SameType(t testers.Tester, expected any) (c testers.Check[any]) {
	defer handlePanic(t, &c)
	getHelper(t)()
	expT := reflect.TypeOf(expected)
	return newCheck(t, func(b *testee, actual any) {
		if expT != reflect.TypeOf(actual) {
			b.Should(`be the same expected type`).
				WithValue(`Expected Value`, expected).
				WithType(`Expected Type`, expected)
		}
	}).setTextHint(expected)
}

// NotSameType creates a check that the actual value is not the same type as the unexpected type.
//
// Example: check.NotSameType(t, 3.2).Assert(actual)
func NotSameType(t testers.Tester, unexpected any) (c testers.Check[any]) {
	defer handlePanic(t, &c)
	getHelper(t)()
	expT := reflect.TypeOf(unexpected)
	return newCheck(t, func(b *testee, actual any) {
		if expT == reflect.TypeOf(actual) {
			b.Should(`not be the unexpected type`).
				WithValue(`Unexpected Value`, unexpected).
				WithType(`Unexpected Type`, unexpected)
		}
	}).setTextHint(unexpected)
}

// Same creates a check that the given expected type is equal to the
// actual type using the `==` comparison. This can be used to ensure
// that two pointers point to the same object.
//
// Example: check.Same(t, expected).Assert(actual)
func Same[T comparable](t testers.Tester, expected T) (c testers.Check[T]) {
	defer handlePanic(t, &c)
	getHelper(t)()
	return newPred(t, func(actual T) bool {
		return actual == expected
	}, `be the same`).
		WithValue(`Expected Value`, expected).
		WithType(`Expected Type`, expected)
}

// NotSame creates a check that the given expected type is not equal to the
// actual type using the `==` comparison. This can be used to ensure
// that two pointers point to different objects.
//
// Example: check.NotSame(t, expected).Assert(actual)
func NotSame[T comparable](t testers.Tester, expected T) (c testers.Check[T]) {
	defer handlePanic(t, &c)
	getHelper(t)()
	return newPred(t, func(actual T) bool {
		return actual != expected
	}, `not be the same`).
		WithValue(`Expected Value`, expected).
		WithType(`Expected Type`, expected)
}

// Includes creates a check that the actual collection of values
// contains all of the given expected values.
//
// There must be at least one expected element. This doesn't check number of
// occurrences in the actual slice meaning that multiple expected elements
// has no effect and will simply match the same value in the slice.
//
// If a map is given to either the expected or actual values.
// The values being matched will be key/value tuples.
//
// Example: check.Contains(t, []int{3, 7, 10}).Assert(actual)
func Includes(t testers.Tester, expected any) (c testers.Check[any]) {
	defer handlePanic(t, &c)
	getHelper(t)()
	expV := readonlyVariantList.Wrap(expected)
	if expV.Count() <= 0 {
		newTestee(t).WithType(`Expected Type`, expected).
			SetupMust(`provide at least one expected value`)
		return (*checkImp[any])(nil)
	}

	return newCheck(t, func(b *testee, actual any) {
		actV := readonlyVariantList.Wrap(actual)
		missing := enumerator.Subtract(actV.Enumerate(), expV.Enumerate()).ToSlice()
		if len(missing) > 0 {
			b.setTextHint(actual).
				WithValue(`Expected Values`, expected).
				WithType(`Expected Type`, expected).
				With(`Missing Values`, missing).
				Should(`have the expected values`)
		}
	}).setTextHint(expected)
}

// OneOf creates a check that the actual value is one of the values
// in the given expected values.
//
// There must be at least one expected element.
// If a map is given to either the expected or actual values.
// The values being matched will be key/value tuples.
//
// Example: check.OneOf(t, []int{3, 7, 10}).Assert(actual)
func OneOf(t testers.Tester, expected any) (c testers.Check[any]) {
	defer handlePanic(t, &c)
	getHelper(t)()
	expV := readonlyVariantList.Wrap(expected)
	if expV.Count() <= 0 {
		newTestee(t).WithType(`Expected Type`, expected).
			SetupMust(`provide at least one expected value`)
		return (*checkImp[any])(nil)
	}

	return newCheck(t, func(b *testee, actual any) {
		if !expV.Contains(actual) {
			b.setTextHint(actual).
				WithValue(`Expected Values`, expected).
				WithType(`Expected Type`, expected).
				Should(`be one of the expected values`)
		}
	}).setTextHint(expected)
}

// Excludes creates a check that the actual collection of values
// does not contain any of the given unexpected values.
//
// There must be at least one unexpected element.
// If a map is given to either the expected or actual values.
// The values being matched will be key/value tuples.
//
// Example: check.Excludes(t, []int{3, 7, 10}).Assert(actual)
func Excludes(t testers.Tester, unexpected any) (c testers.Check[any]) {
	defer handlePanic(t, &c)
	getHelper(t)()
	expV := readonlyVariantList.Wrap(unexpected)
	if expV.Count() <= 0 {
		newTestee(t).WithType(`Expected Type`, unexpected).
			SetupMust(`provide at least one unexpected value`)
		return (*checkImp[any])(nil)
	}

	return newCheck(t, func(b *testee, actual any) {
		actV := readonlyVariantList.Wrap(actual)
		overlap := enumerator.Intersection(actV.Enumerate(), expV.Enumerate()).ToSlice()
		if len(overlap) > 0 {
			b.setTextHint(actual).
				WithValue(`Unexpected Values`, unexpected).
				WithType(`Expected Type`, unexpected).
				With(`Contained Values`, overlap).
				Should(`not have the any of the unexpected values`)
		}
	}).setTextHint(unexpected)
}

// Intersects creates a check that the actual collection of values
// contains at least one of the given expected values.
//
// There must be at least one expected element.
// If a map is given to either the expected or actual values.
// The values being matched will be key/value tuples.
//
// Example: check.Intersects(t, []int{3, 7, 10}).Assert(actual)
func Intersects(t testers.Tester, expected any) (c testers.Check[any]) {
	defer handlePanic(t, &c)
	getHelper(t)()
	expV := readonlyVariantList.Wrap(expected)
	if expV.Count() <= 0 {
		newTestee(t).WithType(`Expected Type`, expected).
			SetupMust(`provide at least one expected value`)
		return (*checkImp[any])(nil)
	}

	return newCheck(t, func(b *testee, actual any) {
		actV := readonlyVariantList.Wrap(actual)
		overlap := enumerator.Intersection(actV.Enumerate(), expV.Enumerate()).ToSlice()
		if len(overlap) <= 0 {
			b.setTextHint(actual).
				WithValue(`Expected Values`, expected).
				WithType(`Expected Type`, expected).
				Should(`have at least one of the expected values`)
		}
	}).setTextHint(expected)
}

// Sorted creates a check that the actual collection of values is sorted.
//
// This may take an optional comparer. The given type parameter is the
// type of elements in the actual value that are being compared.
//
// Example: check.Sorted(t).Assert(actual)
func Sorted[TElem any](t testers.Tester, comparer ...comp.Comparer[TElem]) (c testers.Check[any]) {
	defer handlePanic(t, &c)
	getHelper(t)()
	return newCheck(t, func(b *testee, actual any) {
		actV := readonlyVariantList.Cast[TElem](readonlyVariantList.Wrap(actual))
		if !actV.Enumerate().Sorted(comparer...) {
			b.setTextHint(actual).
				Should(`be in sorted order`)
		}
	})
}

// NotSorted creates a check that the actual collection of values is not sorted.
//
// This may take an optional comparer. The given type parameter is the
// type of elements in the actual value that are being compared.
//
// Example: check.NotSorted(t).Assert(actual)
func NotSorted[TElem any](t testers.Tester, comparer ...comp.Comparer[TElem]) (c testers.Check[any]) {
	defer handlePanic(t, &c)
	getHelper(t)()
	return newCheck(t, func(b *testee, actual any) {
		actV := readonlyVariantList.Cast[TElem](readonlyVariantList.Wrap(actual))
		if actV.Enumerate().Sorted(comparer...) {
			b.setTextHint(actual).
				Should(`not be in sorted order`)
		}
	})
}

// DescendingSorted creates a check that the actual collection of values is
// sorted on descending order.
//
// This may take an optional comparer. The given type parameter is the
// type of elements in the actual value that are being compared.
//
// Example: check.DescendingSorted(t).Assert(actual)
func DescendingSorted[TElem any](t testers.Tester, comparer ...comp.Comparer[TElem]) (c testers.Check[any]) {
	defer handlePanic(t, &c)
	getHelper(t)()
	cmp := comp.Descender(optional.Comparer(comparer))
	return newCheck(t, func(b *testee, actual any) {
		actV := readonlyVariantList.Cast[TElem](readonlyVariantList.Wrap(actual))
		if !actV.Enumerate().Sorted(cmp) {
			b.setTextHint(actual).
				Should(`be in descending sorted order`)
		}
	})
}

// NotDescendingSorted creates a check that the actual collection of values
// is not sorted in descending order.
//
// This may take an optional comparer. The given type parameter is the
// type of elements in the actual value that are being compared.
//
// Example: check.NotDescendingSorted(t).Assert(actual)
func NotDescendingSorted[TElem any](t testers.Tester, comparer ...comp.Comparer[TElem]) (c testers.Check[any]) {
	defer handlePanic(t, &c)
	getHelper(t)()
	cmp := comp.Descender(optional.Comparer(comparer))
	return newCheck(t, func(b *testee, actual any) {
		actV := readonlyVariantList.Cast[TElem](readonlyVariantList.Wrap(actual))
		if actV.Enumerate().Sorted(cmp) {
			b.setTextHint(actual).
				Should(`not be in descending sorted order`)
		}
	})
}

// Unique creates a check that the actual collection has all unique values.
//
// Example: check.Unique(t).Assert(actual)
func Unique[TElem comparable](t testers.Tester) (c testers.Check[any]) {
	defer handlePanic(t, &c)
	getHelper(t)()
	return newCheck(t, func(b *testee, actual any) {
		actV := readonlyVariantList.Cast[TElem](readonlyVariantList.Wrap(actual))
		if !enumerator.IsUnique[TElem](actV.Enumerate()) {
			b.setTextHint(actual).
				Should(`have unique values`)
		}
	})
}

// NotUnique creates a check that the actual collection does not have all unique values.
//
// Example: check.NotUnique(t).Assert(actual)
func NotUnique[TElem comparable](t testers.Tester) (c testers.Check[any]) {
	defer handlePanic(t, &c)
	getHelper(t)()
	return newCheck(t, func(b *testee, actual any) {
		actV := readonlyVariantList.Cast[TElem](readonlyVariantList.Wrap(actual))
		if enumerator.IsUnique[TElem](actV.Enumerate()) {
			b.setTextHint(actual).
				Should(`not have all unique values`)
		}
	})
}

// HasKeys creates a check that the actual map contains all of the given expected keys.
//
// There must be at least one expected key.
//
// Example: check.HasKeys[map[string]int](t, `Name`, `Pet`, `Car`).Assert(actual)
func HasKeys[M ~map[TKey]TValue, TKey comparable, TValue any](t testers.Tester, expected ...TKey) (c testers.Check[M]) {
	defer handlePanic(t, &c)
	getHelper(t)()
	if len(expected) <= 0 {
		newTestee(t).WithType(`Expected Type`, expected).
			SetupMust(`provide at least one expected key`)
		return (*checkImp[M])(nil)
	}

	return newCheck(t, func(b *testee, actual M) {
		missing := []TKey{}
		for _, expKey := range expected {
			if _, has := actual[expKey]; !has {
				missing = append(missing, expKey)
			}
		}
		if len(missing) > 0 {
			b.WithValue(`Expected Keys`, expected).
				WithType(`Expected Type`, expected).
				With(`Missing Keys`, missing).
				Should(`have the expected keys`)
		}
	})
}

// HasValues creates a check that the actual map contains all of the given expected values.
//
// There must be at least one expected key. This doesn't check number of
// occurrences in the actual map meaning that multiple expected values has
// no effect and will simply match the same value in the map.
//
// Example: check.HasValues[map[string]int](t, 4, 5, 6).Assert(actual)
func HasValues[M ~map[TKey]TValue, TKey, TValue comparable](t testers.Tester, expected ...TValue) (c testers.Check[M]) {
	defer handlePanic(t, &c)
	getHelper(t)()
	if len(expected) <= 0 {
		newTestee(t).WithType(`Expected Type`, expected).
			SetupMust(`provide at least one expected value`)
		return (*checkImp[M])(nil)
	}
	expV := enumerator.Enumerate(expected...)

	return newCheck(t, func(b *testee, actual M) {
		actV := enumerator.Enumerate(utils.Values(actual)...)
		missing := enumerator.Subtract(actV, expV).ToSlice()
		if len(missing) > 0 {
			b.WithValue(`Expected Values`, expected).
				WithType(`Expected Type`, expected).
				With(`Missing Values`, missing).
				Should(`have the expected values`)
		}
	})
}

// EqualElems creates a check that the actual slice contains all of the given
// expected elements and no others in any order while ignoring repeats.
//
// This doesn't check number of occurrences in the actual slice meaning that
// multiple expected elements has no effect and will simply match the
// same value in the slice.
//
// If a map is given to either the expected or actual values.
// The values being matched will be key/value tuples.
//
// Example: check.EqualElems(t, []int{3, 7, 10}).Assert(actual)
func EqualElems(t testers.Tester, expected any) (c testers.Check[any]) {
	defer handlePanic(t, &c)
	getHelper(t)()
	expV := readonlyVariantList.Wrap(expected)
	return newCheck(t, func(b *testee, actual any) {
		actV := readonlyVariantList.Wrap(actual)
		missing := enumerator.Subtract(actV.Enumerate(), expV.Enumerate()).ToSlice()
		extra := enumerator.Subtract(expV.Enumerate(), actV.Enumerate()).ToSlice()
		if len(missing) > 0 || len(extra) > 0 {
			b.setTextHint(actual).
				WithValue(`Expected Elements`, expected).
				WithType(`Expected Type`, expected).
				Should(`have the expected elements`)
			if len(missing) > 0 {
				b.With(`Missing Elements`, b.formatUniqueValues(missing))
			}
			if len(extra) > 0 {
				b.With(`Extra Elements`, b.formatUniqueValues(extra))
			}
		}
	}).setTextHint(expected)
}

// SameElems creates a check that the actual slice contains all of the given
// expected elements and no others in any order and in the same number.
// The types being compared to not need to match but the elements do.
//
// There must be at least one expected element. This doesn't check number of
// occurrences in the actual slice meaning that multiple expected elements
// has no effect and will simply match the same value in the slice.
//
// If a map is given to either the expected or actual values.
// The values being matched will be key/value tuples.
//
// Example: check.HasElems(t, []int{3, 7, 10}).Assert(actual)
func SameElems(t testers.Tester, expected any) (c testers.Check[any]) {
	defer handlePanic(t, &c)
	getHelper(t)()
	expV := readonlyVariantList.Wrap(expected)
	expCounts := enumerator.DuplicateCounts(expV.Enumerate())

	return newCheck(t, func(b *testee, actual any) {
		expCopy := maps.Clone(expCounts)
		actV := readonlyVariantList.Wrap(actual)
		actV.Backwards().Foreach(func(value any) {
			expCopy[value]--
		})
		b.setTextHint(actual)

		missing := []string{}
		extra := []string{}
		isWrong := false
		for key, diff := range expCopy {
			if diff == 0 {
				continue
			}
			isWrong = true
			str := b.formatValue(key)
			if count := abs(diff); count > 1 {
				str += fmt.Sprintf(`(x%d)`, count)
			}
			if diff > 0 {
				missing = append(missing, str)
			} else {
				extra = append(extra, str)
			}
		}

		if isWrong {
			b.WithValue(`Expected Elements`, expected).
				WithType(`Expected Type`, expected).
				Should(`have the expected elements`)
			if len(missing) > 0 {
				sort.Strings(missing)
				b.With(`Missing Elements`, `[`+strings.Join(missing, ` `)+`]`)
			}
			if len(extra) > 0 {
				sort.Strings(extra)
				b.With(`Extra Elements`, `[`+strings.Join(extra, ` `)+`]`)
			}
		}
	}).setTextHint(expected)
}
