package check

import (
	"bytes"
	"fmt"
	"maps"
	"sort"
	"strconv"
	"strings"

	"goToolbox/internal/simpleSet"
	"goToolbox/testers"
	"goToolbox/utils"
)

const (
	defaultMaxLines  = 100
	defaultTailLines = 3
	newLine          = "\n"
	indent           = "\t"
	ellipsis         = `...`
)

type testee struct {
	t         testers.Tester
	failed    bool
	action    string
	must      bool
	maxLines  int
	tailLines int
	textHint  bool
	pContext  map[string]string
	fContext  map[string]any
}

func newTestee(t testers.Tester) *testee {
	return &testee{
		t:         t,
		failed:    false,
		action:    ``,
		must:      false,
		maxLines:  defaultMaxLines,
		tailLines: defaultTailLines,
		textHint:  false,
		pContext:  map[string]string{},
		fContext:  map[string]any{},
	}
}

func (b *testee) Copy() *testee {
	b2 := *b
	b2.pContext = maps.Clone(b.pContext)
	b2.fContext = maps.Clone(b.fContext)
	return &b2
}

func (b *testee) With(key string, args ...any) *testee {
	b.pContext[key] = b.limitLines(fmt.Sprint(args...))
	return b
}

func (b *testee) Withf(key, format string, args ...any) *testee {
	b.pContext[key] = b.limitLines(fmt.Sprintf(format, args...))
	return b
}

func (b *testee) WithType(key string, valueForType any) *testee {
	return b.Withf(key, `%T`, valueForType)
}

func (b *testee) WithValue(key string, value any) *testee {
	b.fContext[key] = value
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
	if context := b.formatContext(); len(context) > 0 {
		message += `:` + context
	}
	b.t.Error(message + "\n")
	if b.must {
		b.t.FailNow()
	}
}

func (b *testee) formatContext() string {
	context := maps.Clone(b.pContext)
	for key, value := range b.fContext {
		context[key] = b.limitLines(b.formatValue(value))
	}
	keys := utils.SortedKeys(context)
	maxWidth := utils.GetMaxStringLen(keys) + 2
	padding := newLine + indent + strings.Repeat(` `, maxWidth)
	buf := &bytes.Buffer{}
	for _, key := range keys {
		value := strings.ReplaceAll(context[key], newLine, padding)
		_, _ = buf.WriteString(newLine)
		_, _ = buf.WriteString(indent)
		_, _ = buf.WriteString(fmt.Sprintf(`%-*s`, maxWidth, key+`: `))
		_, _ = buf.WriteString(value)
	}
	return buf.String()
}

func (b *testee) limitLines(value string) string {
	maxLines := max(b.maxLines, 3)
	if strings.Count(value, newLine) <= maxLines {
		return value
	}

	lines := strings.Split(value, newLine)
	count := len(lines)
	result := []string{}
	tailLines := min(max(b.tailLines, 0), maxLines-2)
	if headEnd := maxLines - tailLines - 1; headEnd > 0 {
		result = append(result, lines[:headEnd]...)
	}

	result = append(result, ellipsis)
	if tailStart := count - tailLines; tailStart < count {
		result = append(result, lines[tailStart:]...)
	}

	return strings.Join(result, newLine)
}

func (b *testee) setTextHint(value any) *testee {
	if _, ok := value.(string); ok {
		b.textHint = true
	}
	return b
}

func (b *testee) formatValue(value any) string {
	switch t := value.(type) {
	case nil:
		return `<nil>`
	case byte:
		if b.textHint {
			return strconv.QuoteRuneToGraphic(rune(t))
		}
	case rune:
		if b.textHint {
			return strconv.QuoteRuneToGraphic(t)
		}
	case string:
		return strconv.QuoteToGraphic(t)
	case utils.Stringer:
		return strconv.QuoteToGraphic(t.String())
	}
	return utils.String(value)
}

func (b *testee) formatUniqueValues(values []any) string {
	values = simpleSet.With(values...).ToSlice()
	parts := make([]string, len(values))
	for i, value := range values {
		parts[i] = b.formatValue(value)
	}
	sort.Strings(parts)
	return `[` + strings.Join(parts, ` `) + `]`
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

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
