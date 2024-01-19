package observer

import (
	"bytes"
	"testing"

	"github.com/Snow-Gremlin/goToolbox/events/event"
	"github.com/Snow-Gremlin/goToolbox/internal/liteUtils"
)

func Test_Observer(t *testing.T) {
	buf := &bytes.Buffer{}
	obv0 := New[string](nil)
	obv1 := New(func(value string) {
		_, _ = buf.WriteString(`[` + value + `]`)
	})

	e := event.New[string]()
	checkEqual(t, true, e.Add(obv0))
	checkEqual(t, false, e.Add(obv0))
	checkEqual(t, true, e.Add(obv1))
	checkEqual(t, false, e.Add(obv1))
	checkBuf(t, buf, ``)

	e.Invoke(`Hello`)
	e.Invoke(`World`)
	checkBuf(t, buf, `[Hello][World]`)

	checkEqual(t, true, e.Remove(obv0))
	checkEqual(t, false, e.Remove(obv0))
	checkEqual(t, true, e.Remove(obv1))
	checkEqual(t, false, e.Remove(obv1))
	checkBuf(t, buf, ``)

	e.Invoke(`Goodbye`)
	checkBuf(t, buf, ``)
}

func checkEqual(t *testing.T, exp, actual any) {
	t.Helper()
	if !liteUtils.Equal(exp, actual) {
		t.Errorf("unexpected result:\n\tactual:   %v\n\texpected: %v\n", actual, exp)
	}
}

func checkBuf(t *testing.T, buf *bytes.Buffer, exp string) {
	t.Helper()
	checkEqual(t, exp, buf.String())
	buf.Reset()
}
