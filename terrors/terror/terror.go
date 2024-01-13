package terror

import (
	"slices"

	"github.com/Snow-Gremlin/goToolbox/internal/liteUtils"
	"github.com/Snow-Gremlin/goToolbox/terrors"
	"github.com/Snow-Gremlin/goToolbox/terrors/stacked"
)

// New creates a new error with the given message
// and an optional error.
func New(msg string, errs ...error) terrors.TError {
	if len(msg) <= 0 {
		msg = `unknown error`
	}

	errs = slices.DeleteFunc(errs, liteUtils.IsZero)

	return &tErrorImp{
		msg:     msg,
		errs:    slices.Clip(errs),
		stack:   stacked.Stack(1, 0),
		context: nil,
	}
}

// RecoveredPanic creates an error for a recovered panic.
func RecoveredPanic(r any) terrors.TError {
	switch e := r.(type) {
	case nil:
		return nil
	case terrors.TError:
		return e
	case error:
		return New(`recovered panic`).
			WithError(e)
	default:
		return New(`recovered panic`).
			With(`recovered`, r)
	}
}

// OutOfBounds creates an index out of bounds error.
func OutOfBounds(index, count int) terrors.TError {
	return New(`index out of bounds`).
		With(`index`, index).
		With(`count`, count)
}

// EmptyCollection creates a no values in collection error.
func EmptyCollection(action string) terrors.TError {
	return New(`collection contains no values`).
		With(`action`, action)
}

// InvalidArgCount creates an invalid number of arguments error.
func InvalidArgCount(max, count int, usage string) terrors.TError {
	return New(`invalid number of arguments`).
		With(`maximum`, max).
		With(`count`, count).
		With(`usage`, usage)
}

// NilArg creates an invalid nil argument error.
func NilArg(name string) terrors.TError {
	return New(`argument may not be nil`).
		With(`name`, name)
}

// UnstableIteration creates an error for when a
// collection is modified in a way that could make
// continuing any iteration for that collection unstable.
func UnstableIteration() terrors.TError {
	return New(`Collection was modified; iteration may not continue`)
}
