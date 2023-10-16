package check

import (
	"bytes"
	"fmt"
	"maps"
	"strings"

	"goToolbox/testers"
	"goToolbox/utils"
)

type testee struct {
	t       testers.Tester
	failed  bool
	action  string
	must    bool
	context map[string]string
}

func newTestee(t testers.Tester) *testee {
	return &testee{
		t:       t,
		failed:  false,
		action:  ``,
		must:    false,
		context: map[string]string{},
	}
}

func (b *testee) Copy() *testee {
	return &testee{
		t:       b.t,
		failed:  b.failed,
		action:  b.action,
		must:    b.must,
		context: maps.Clone(b.context),
	}
}

func (b *testee) With(key string, args ...any) *testee {
	b.context[key] = limitLines(fmt.Sprint(args...))
	return b
}

func (b *testee) Withf(key, format string, args ...any) *testee {
	b.context[key] = limitLines(fmt.Sprintf(format, args...))
	return b
}

func (b *testee) Required() *testee {
	b.must = true
	return b
}

func (b *testee) Should(action string) *testee {
	b.failed = true
	b.action = action
	return b
}

func (b *testee) SetupMust(action string) {
	b.Required().Should(action).Finish()
}

func (b *testee) Finish() {
	if !b.failed {
		return
	}

	imperative := `Should`
	if b.must {
		imperative = `Must`
	}
	message := fmt.Sprintf("\n%s %s", imperative, b.action)
	if len(b.context) > 0 {
		message += `:` + formatContext(b.context)
	}
	b.t.Error(message + "\n")
	if b.must {
		b.t.FailNow()
	}
}

// handlePanic handles a panic in a check.
func handlePanic[T any](t testers.Tester, pc *testers.Check[T]) {
	if r := recover(); r != nil {
		t.Error("\nError: " + utils.String(r) + "\n")
		t.FailNow()
		*pc = (*checkImp[T])(nil)
	}
}

// getHelper returns the helper function to call
// or a noop if no helper is found.
func getHelper(t testers.Tester) func() {
	if h, ok := t.(interface{ Helper() }); ok {
		return h.Helper
	}
	return func() {} // Noop
}

// limitLines caps the number of lines in the string to a specific number of lines.
func limitLines(value string) string {
	const (
		maxLines  = 100
		tailLines = 3
		newLine   = "\n"
		ellipsis  = `...`
		headEnd   = maxLines - tailLines - 1
	)
	if strings.Count(value, newLine) > maxLines {
		lines := strings.Split(value, newLine)
		if count := len(lines); count > maxLines {
			tailStart := count - tailLines
			lines = append(append(lines[:headEnd], ellipsis), lines[tailStart:]...)
			return strings.Join(lines, newLine)
		}
	}
	return value
}

// formatContext formats the given context into a string.
// The context is indented and the values, even multi-lined values are aligned.
func formatContext(context map[string]string) string {
	const newline = "\n"
	const indent = "\t"
	keys := utils.SortedKeys(context)
	maxWidth := utils.GetMaxStringLen(keys) + 2

	padding := newline + indent + strings.Repeat(` `, maxWidth)
	buf := &bytes.Buffer{}
	for _, key := range keys {
		value := strings.ReplaceAll(context[key], newline, padding)
		_, _ = buf.WriteString(newline)
		_, _ = buf.WriteString(indent)
		_, _ = buf.WriteString(fmt.Sprintf(`%-*s`, maxWidth, key+`: `))
		_, _ = buf.WriteString(value)
	}
	return buf.String()
}
