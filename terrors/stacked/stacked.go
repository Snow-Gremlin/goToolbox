package stacked

import (
	"runtime/debug"
	"strings"

	"github.com/Snow-Gremlin/goToolbox/internal/liteUtils"
	"github.com/Snow-Gremlin/goToolbox/terrors"
)

// Stack gets the stack trace not including the call into this method.
// Offset is the number of frames to leave off of the top of the stack.
// Trim is the number of frames to remove from the end of the stack.
func Stack(offset, trim int) string {
	stack := strings.TrimSpace(string(debug.Stack()))
	lines := strings.Split(stack, "\n")
	total := len(lines)
	if total >= 5 {
		goroutine := lines[0]
		start := min(offset*2+5, total)
		stop := max(total-trim*2, 0)
		if start < stop {
			lines = lines[start:stop]
			stack = goroutine + "\n" + strings.Join(lines, "\n")
		}
	}
	return stack
}

// DeepestStacked finds the deepest error which has a stack.
// This may not get the deepest in the whole tree, only the deepest
// as found walking down one branch of the tree.
func DeepestStacked(err error) terrors.Stacked {
	s, _ := deepestStacked(err)
	return s
}

func deepestStacked(err error) (terrors.Stacked, bool) {
	if liteUtils.IsNil(err) {
		return nil, false
	}

	if e, ok := err.(terrors.MonoWrap); ok {
		if s, ok := deepestStacked(e.Unwrap()); ok {
			return s, true
		}
	} else if e, ok := err.(terrors.MultiWrap); ok {
		for _, w := range e.Unwrap() {
			if s, ok := deepestStacked(w); ok {
				return s, true
			}
		}
	}

	s, ok := err.(terrors.Stacked)
	return s, ok
}
