package changeType

import "testing"

func Test_ChangeType(t *testing.T) {
	check(t, Added, `Added`, true)
	check(t, Removed, `Removed`, true)
	check(t, Replaced, `Replaced`, true)
	check(t, invalid, `invalid`, false)
	check(t, `hello`, `invalid`, false)
	check(t, ``, `invalid`, false)
}

func check(t *testing.T, c ChangeType, expStr string, expValid bool) {
	actualStr, actualValid := c.String(), c.Valid()
	if actualStr != expStr || actualValid != expValid {
		t.Errorf("unexpected result:\n\tactual:   %s (%t)\n\texpected: %s (%t)\n",
			actualStr, actualValid, expStr, expValid)
	}
}
