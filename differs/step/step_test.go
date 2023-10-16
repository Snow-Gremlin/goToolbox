package step

import (
	"fmt"
	"testing"
)

func stepString(t *testing.T, s Step, exp string) {
	if result := s.String(); result != exp {
		t.Error(fmt.Sprint("Unexpected step string:",
			"\n   Actual:    ", result,
			"\n   Expected: ", exp))
	}
}

func Test_Type(t *testing.T) {
	stepString(t, -1, `?`)
	stepString(t, Equal, `=`)
	stepString(t, Added, `+`)
	stepString(t, Removed, `-`)
	stepString(t, 4, `?`)
}
