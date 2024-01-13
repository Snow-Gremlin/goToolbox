package check

import (
	"bytes"
	"fmt"
	"maps"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Snow-Gremlin/goToolbox/internal/simpleSet"
	"github.com/Snow-Gremlin/goToolbox/testers"
	"github.com/Snow-Gremlin/goToolbox/utils"
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
	intHint   int
	timeHint  string
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
		intHint:   10,
		timeHint:  ``,
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
	getHelper(b.t)()

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

func tryFormatString(value any) (string, bool) {
	switch t := value.(type) {
	case byte:
		return strconv.QuoteRuneToGraphic(rune(t)), true
	case rune:
		return strconv.QuoteRuneToGraphic(t), true
	case []byte:
		return strconv.QuoteToGraphic(string(t)), true
	case []rune:
		return strconv.QuoteToGraphic(string(t)), true
	}
	return ``, false
}

func formatValue(val, prefix string, digitGroup int) string {
	if strings.HasPrefix(val, `-`) {
		prefix = `-` + prefix
		val = val[1:]
	}
	tail := ``
	for len(val) > digitGroup {
		high := len(val) - digitGroup
		tail = `_` + val[high:] + tail
		val = val[:high]
	}
	return prefix + val + tail
}

func tryFormatInt(value any, radix int) (string, bool) {
	switch value.(type) {
	case int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64:
		switch radix {
		case 2:
			return formatValue(fmt.Sprintf(`%b`, value), `0b`, 4), true
		case 8:
			return formatValue(fmt.Sprintf(`%o`, value), `0o`, 4), true
		case 16:
			return formatValue(fmt.Sprintf(`%X`, value), `0x`, 4), true
		}
	}
	return ``, false
}

func (b *testee) formatValue(value any) string {
	switch t := value.(type) {
	case nil:
		return `<nil>`
	case string:
		return strconv.QuoteToGraphic(t)
	case time.Time:
		if len(b.timeHint) > 0 {
			return t.Format(b.timeHint)
		}
		return t.String()
	}

	if b.textHint {
		if str, ok := tryFormatString(value); ok {
			return str
		}
	}

	if b.intHint != 10 {
		if str, ok := tryFormatInt(value, b.intHint); ok {
			return str
		}
	}

	if str, ok := value.(utils.Stringer); ok {
		return strconv.QuoteToGraphic(str.String())
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

// getRegExp returns the regular expression from the given
// string pattern or regular expression instance.
func getRegExp(t testers.Tester, regex any) *regexp.Regexp {
	getHelper(t)()
	switch trex := regex.(type) {
	case string:
		if len(trex) <= 0 {
			newTestee(t).SetupMust(`provide a non-empty regular expression pattern`)
			return nil
		}
		re, err := regexp.Compile(trex)
		if err != nil {
			newTestee(t).With(`Pattern`, trex).
				SetupMust(`provide a valid regular expression pattern`)
			return nil
		}
		return re
	case *regexp.Regexp:
		if utils.IsNil(trex) {
			newTestee(t).SetupMust(`provide a non-nil regular expression instance`)
			return nil
		}
		return trex
	default:
		newTestee(t).WithType(`Given Type`, regex).
			SetupMust(`provide a string pattern or *regexp.Regexp`)
		return nil
	}
}
