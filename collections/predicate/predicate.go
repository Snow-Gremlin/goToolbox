package predicate

import (
	"regexp"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/internal/optional"
	"github.com/Snow-Gremlin/goToolbox/internal/simpleSet"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

// IsNil is a predicate which returns true if the given value is nil.
// This will return false if the value is not nil-able.
func IsNil[T any]() collections.Predicate[T] {
	return func(value T) bool {
		return utils.IsNil(value)
	}
}

// IsNotNil is a predicate which returns true if the given value is not nil.
// This will return false if the value is not nil-able.
func IsNotNil[T any]() collections.Predicate[T] {
	return Not(IsNil[T]())
}

// IsZero is a predicate which returns true if the the value is zero.
func IsZero[T any]() collections.Predicate[T] {
	return utils.IsZero[T]
}

// IsNotZero is a predicate which returns true if the the value is zero.
func IsNotZero[T any]() collections.Predicate[T] {
	return Not(IsZero[T]())
}

// IsTrue is a predicate which returns true if the the value is true.
func IsTrue() collections.Predicate[bool] {
	return func(value bool) bool {
		return value
	}
}

// IsFalse is a predicate which returns true if the the value is false.
func IsFalse() collections.Predicate[bool] {
	return func(value bool) bool {
		return !value
	}
}

// OfType is a predicate which checks if the given value is the given target type.
func OfType[Target, T any]() collections.Predicate[T] {
	return func(value T) bool {
		_, ok := any(value).(Target)
		return ok
	}
}

// Matches is a predicate which checks if a string matches the given
// regular expression pattern.
//
// If the pattern is empty or not a valid regular expression,
// then this will panic.
func Matches(pattern string) collections.Predicate[string] {
	if len(pattern) <= 0 {
		panic(terror.New(`may not used an empty pattern string`))
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		panic(terror.New(`invalid regular expression pattern`).
			With(`pattern`, pattern).
			WithError(err))
	}
	return re.MatchString
}

// InMap is a predicate which returns true if the given value exists as a key in the given map.
// The map is used as a set so the values in the map are not used.
func InMap[TKey comparable, TValue any, M ~map[TKey]TValue](m M) collections.Predicate[TKey] {
	return func(value TKey) bool {
		_, exists := m[value]
		return exists
	}
}

// In is a predicate which returns true if the given value is in the
// set of values given to create the predicate.
//
// This will create a map for fast lookup.
func In[T comparable](values ...T) collections.Predicate[T] {
	return InMap(simpleSet.With(values...))
}

// AsString is a predicate which calls the given predicate with
// the string of the given value.
func AsString(p collections.Predicate[string]) collections.Predicate[any] {
	return func(value any) bool {
		return p(utils.String(value))
	}
}

// Eq is a predicate which returns true if the given value
// is equal to the value passed into the predicate.
func Eq[T any](value T) collections.Predicate[T] {
	return func(query T) bool {
		return utils.Equal(query, value)
	}
}

// NotEq is a predicate which returns true if the given value
// is not equal to the value passed into the predicate.
func NotEq[T any](value T) collections.Predicate[T] {
	return Not(Eq[T](value))
}

// GreaterThan is a predicate which returns true if the value
// passed into the predicate is greater than the given value.
func GreaterThan[T any](value T, comparer ...utils.Comparer[T]) collections.Predicate[T] {
	cmp := optional.Comparer(comparer)
	return func(query T) bool {
		return cmp(query, value) > 0
	}
}

// GreaterEq is a predicate which returns true if the value
// passed into the predicate is greater than or equal to the given value.
func GreaterEq[T any](value T, comparer ...utils.Comparer[T]) collections.Predicate[T] {
	cmp := optional.Comparer(comparer)
	return func(query T) bool {
		return cmp(query, value) >= 0
	}
}

// LessThan is a predicate which returns true if the value
// passed into the predicate is less than the given value.
func LessThan[T any](value T, comparer ...utils.Comparer[T]) collections.Predicate[T] {
	cmp := optional.Comparer(comparer)
	return func(query T) bool {
		return cmp(query, value) < 0
	}
}

// LessEq is a predicate which returns true if the value
// passed into the predicate is less than or equal to the given value.
func LessEq[T any](value T, comparer ...utils.Comparer[T]) collections.Predicate[T] {
	cmp := optional.Comparer(comparer)
	return func(query T) bool {
		return cmp(query, value) <= 0
	}
}

// InRange is a predicate which returns true if the value passed
// into the predicate is between the given min and maximum inclusively.
func InRange[T any](min, max T, comparer ...utils.Comparer[T]) collections.Predicate[T] {
	cmp := optional.Comparer(comparer)
	return func(query T) bool {
		return cmp(min, query) <= 0 && cmp(query, max) <= 0
	}
}

// EpsilonEq is a predicate which returns true if the value passed into the
// predicate is within an epsilon value (inclusively) of the other value.
//
// This is useful for finding floating-point values which have lost precision
// via calculations and will not be equal a literal, but will be very close to it.
func EpsilonEq[T utils.NumConstraint](value, epsilon T) collections.Predicate[T] {
	cmp := utils.EpsilonComparer(epsilon)
	return func(query T) bool {
		return cmp(query, value) == 0
	}
}

// EpsilonNotEq is a predicate which returns true if the value passed into the
// predicate is not within an epsilon value (inclusively) of the other value.
//
// This is useful for finding floating-point values which have lost precision
// via calculations and will not be equal a literal, but will be very close to it.
func EpsilonNotEq[T utils.NumConstraint](value, epsilon T) collections.Predicate[T] {
	cmp := utils.EpsilonComparer(epsilon)
	return func(query T) bool {
		return cmp(query, value) != 0
	}
}

// Pos is a predicate which returns true if the value
// passed into the predicate is greater than zero, i.e. positive.
func Pos[T any](comparer ...utils.Comparer[T]) collections.Predicate[T] {
	cmp, zero := optional.Comparer(comparer), utils.Zero[T]()
	return func(query T) bool {
		return cmp(query, zero) > 0
	}
}

// Neg is a predicate which returns true if the value
// passed into the predicate is less than zero, i.e. negative,
func Neg[T any](comparer ...utils.Comparer[T]) collections.Predicate[T] {
	cmp, zero := optional.Comparer(comparer), utils.Zero[T]()
	return func(query T) bool {
		return cmp(query, zero) < 0
	}
}

// Not negates the result of the given predicate.
func Not[T any](p collections.Predicate[T]) collections.Predicate[T] {
	return func(value T) bool {
		return !p(value)
	}
}

// And is a predicate which returns true only if
// all the given predicates are true.
func And[T any](p ...collections.Predicate[T]) collections.Predicate[T] {
	count := len(p)
	return func(value T) bool {
		for i := 0; i < count; i++ {
			if !p[i](value) {
				return false
			}
		}
		return true
	}
}

// Or is a predicate which returns true if
// any of the given predicated are true.
func Or[T any](p ...collections.Predicate[T]) collections.Predicate[T] {
	count := len(p)
	return func(value T) bool {
		for i := 0; i < count; i++ {
			if p[i](value) {
				return true
			}
		}
		return false
	}
}

// OnlyOne is a predicate which returns true if
// only one of the given predicates are true.
func OnlyOne[T any](p ...collections.Predicate[T]) collections.Predicate[T] {
	count := len(p)
	return func(value T) bool {
		found := false
		for i := 0; i < count; i++ {
			if p[i](value) {
				if found {
					return false
				}
				found = true
			}
		}
		return found
	}
}
