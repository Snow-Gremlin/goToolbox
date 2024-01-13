package terror

import (
	"maps"
	"slices"

	"github.com/Snow-Gremlin/goToolbox/internal/liteUtils"
	"github.com/Snow-Gremlin/goToolbox/terrors"
	"github.com/Snow-Gremlin/goToolbox/terrors/stacked"
)

type tErrorImp struct {
	msg     string
	errs    []error
	stack   string
	context map[string]any
}

func (e *tErrorImp) Message() string {
	return e.msg
}

func (e *tErrorImp) Context() map[string]any {
	return maps.Clone(e.context)
}

func (e *tErrorImp) Error() string {
	return ToString(e)
}

func (e *tErrorImp) String() string {
	return e.Error()
}

func (e *tErrorImp) Stack() string {
	return e.stack
}

func (e *tErrorImp) Unwrap() []error {
	return e.errs
}

func (e *tErrorImp) With(key string, value any) terrors.TError {
	if e.context == nil {
		e.context = make(map[string]any)
	}
	e.context[key] = value
	return e
}

func (e *tErrorImp) WithError(err error) terrors.TError {
	if liteUtils.IsNil(err) {
		return e
	}

	it := Walk(err)
	for it.Next() {
		if it.Current() == e {
			// Loop found so ignore error
			return e
		}
	}

	e.errs = append(e.errs, err)
	return e
}

func (e *tErrorImp) ResetStack(offset int) terrors.TError {
	e.stack = stacked.Stack(offset, 0)
	return e
}

func (e *tErrorImp) Clone() terrors.TError {
	return &tErrorImp{
		msg:     e.msg,
		errs:    slices.Clone(e.errs),
		stack:   e.stack,
		context: maps.Clone(e.context),
	}
}

func (e *tErrorImp) Equals(other any) bool {
	e2, ok := other.(terrors.TError)
	if !ok || e.msg != e2.Message() ||
		!maps.Equal(e.context, e2.Context()) {
		return false
	}

	w2 := e2.Unwrap()
	if len(e.errs) != len(w2) {
		return false
	}

	for i, e := range e.errs {
		if !liteUtils.Equal(e, w2[i]) {
			return false
		}
	}

	// Don't compare stacks
	return true
}
