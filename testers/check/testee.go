package check

import (
	"bytes"
	"fmt"
	"maps"
	"strconv"
	"strings"

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
	textHint  bool
	maxLines  int
	tailLines int
	context   map[string]any
}

func newTestee(t testers.Tester) *testee {
	return &testee{
		t:         t,
		failed:    false,
		action:    ``,
		must:      false,
		textHint:  false,
		maxLines:  defaultMaxLines,
		tailLines: defaultTailLines,
		context:   map[string]any{},
	}
}

func (b *testee) Copy() *testee {
	b2 := *b
	b2.context = maps.Clone(b.context)
	return &b2
}

func (b *testee) With(key string, args ...any) *testee {
	// when context is `[]any` it still needs to be formatted
	b.context[key] = args
	return b
}

func (b *testee) Withf(key, format string, args ...any) *testee {
	// when context is `string` it has been preformatted
	b.context[key] = fmt.Sprintf(format, args...)
	return b
}

func (b *testee) WithType(key string, valueForType any) *testee {
	return b.Withf(key, `%T`, valueForType)
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
		message += `:` + b.formatContext()
	}
	b.t.Error(message + "\n")
	if b.must {
		b.t.FailNow()
	}
}

func (b *testee) formatContext() string {
	keys := utils.SortedKeys(b.context)
	maxWidth := utils.GetMaxStringLen(keys) + 2
	padding := newLine + indent + strings.Repeat(` `, maxWidth)
	buf := &bytes.Buffer{}
	for _, key := range keys {
		_, _ = buf.WriteString(newLine)
		_, _ = buf.WriteString(indent)
		_, _ = buf.WriteString(fmt.Sprintf(`%-*s`, maxWidth, key+`: `))
		_, _ = buf.WriteString(b.valueToString(key, padding))
	}
	return buf.String()
}

func (b *testee) valueToString(key, padding string) string {
	value := b.context[key]
	if str, ok := value.(string); ok {
		return b.limitLines(str)
	}

	parts := value.([]any)
	switch len(parts) {
	case 0:
		return `<empty>`
	case 1:
		return fmt.Sprint(b.formatValuePart(parts[0], true))
	default:
		for i, part := range parts {
			parts[i] = b.formatValuePart(part, false)
		}
		str := b.limitLines(fmt.Sprint(parts...))
		return strings.ReplaceAll(str, newLine, padding)
	}
}

func (b *testee) formatValuePart(part any, quoteHint bool) any {
	switch t := part.(type) {
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
		if quoteHint {
			return strconv.QuoteToGraphic(t)
		}
		return t
	case error:
		return t.Error()
	case utils.Stringer:
		if quoteHint {
			return strconv.QuoteToGraphic(t.String())
		}
		return t.String()
	}
	return part
}

func (b *testee) limitLines(value string) string {
	maxLines := max(b.maxLines, 3)
	if strings.Count(value, newLine) <= maxLines {
		return value
	}

	lines := strings.Split(value, newLine)
	count := len(lines)
	if count <= maxLines {
		return value
	}

	result := []string{}
	tailLines := min(max(b.tailLines, 0), maxLines-2)
	if headEnd := maxLines - tailLines - 1; headEnd <= 0 {
		result = append(result, lines[:headEnd]...)
	}

	result = append(result, ellipsis)
	if tailStart := count - tailLines; tailStart < count {
		result = append(result, lines[tailStart:]...)
	}

	return strings.Join(result, newLine)
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
