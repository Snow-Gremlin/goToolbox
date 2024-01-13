package stacked

import (
	"errors"
	"fmt"
	"regexp"
	"testing"
)

func Test_Stack(t *testing.T) {
	s1 := Stack(0, 2)

	exp := `^goroutine\s\d+\s\[running\]:\s*\n` +
		`github\.com\/Snow-Gremlin\/goToolbox.terrors.stacked\.Test_Stack\(0x[0-9a-f]+\)\n` +
		`\t.+stacked_test\.go:\d+\s\+0x[0-9a-f]+$`
	re := regexp.MustCompile(exp)

	if !re.MatchString(s1) {
		t.Errorf("\nStack failed to match expected:\nStack: %q\nPattern: %s", s1, exp)
	}
}

type pseudoStackErr struct {
	msg string
	err error
}

func (e pseudoStackErr) Error() string {
	return e.msg
}

func (e pseudoStackErr) Stack() string {
	return e.msg
}

func (e pseudoStackErr) Unwrap() error {
	return e.err
}

func Test_DeepestStacked(t *testing.T) {
	var e0 error
	s0 := DeepestStacked(e0)
	if s0 != nil {
		t.Errorf("\nExpected a nil stacked from nil:\nActual: %v", s0)
	}

	e1 := errors.New(`First`)
	s1 := DeepestStacked(e1)
	if s1 != nil {
		t.Errorf("\nExpected a nil stacked from no stack error:\nActual: %v", s1)
	}

	e2 := pseudoStackErr{msg: `Second`, err: nil}
	e3 := pseudoStackErr{msg: `Third`, err: e1}
	e4 := fmt.Errorf(`Forth: [%w, %w]`, e2, e3)
	e5 := pseudoStackErr{msg: `Fifth`, err: e4}
	s5 := DeepestStacked(e5)
	if s5 != e2 {
		t.Errorf("\nExpected a deep stacked error:\nActual: %v\nExpected: %v", s5, e2)
	}
}
