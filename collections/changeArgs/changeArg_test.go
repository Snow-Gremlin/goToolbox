package changeArgs

import (
	"testing"

	"github.com/Snow-Gremlin/goToolbox/collections/changeType"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

func Test_ChangeArg(t *testing.T) {
	check(t, NewAdded().Type(), changeType.Added)
	check(t, NewRemoved().Type(), changeType.Removed)
	check(t, NewReplaced().Type(), changeType.Replaced)
}

func check(t *testing.T, exp, actual any) {
	if !utils.Equal(actual, exp) {
		t.Errorf("Unexpected value:\n\tActual: %v\n\tExpected: %v\n",
			actual, exp)
	}
}
