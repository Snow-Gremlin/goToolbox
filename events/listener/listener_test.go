package listener

import (
	"bytes"
	"testing"

	"github.com/Snow-Gremlin/goToolbox/events/event"
	"github.com/Snow-Gremlin/goToolbox/internal/liteUtils"
)

func Test_Listener(t *testing.T) {
	buf := &bytes.Buffer{}
	lis1 := New(func(value string) {
		_, _ = buf.WriteString(`[1:` + value + `]`)
	})
	lis2 := New(func(value string) {
		_, _ = buf.WriteString(`[2:` + value + `]`)
	})
	e1 := event.New[string]()
	e2 := event.New[string]()
	e3 := event.New[string]()

	checkEqual(t, true, lis1.Subscribe(e1))
	checkEqual(t, true, lis1.Subscribe(e2))
	checkEqual(t, true, lis1.Subscribe(e3))
	checkEqual(t, false, lis1.Subscribe(e2))

	checkEqual(t, true, lis2.Subscribe(e1))
	checkEqual(t, true, lis2.Subscribe(e2))
	checkEqual(t, true, lis2.Subscribe(e3))
	checkBuf(t, buf, ``)

	e1.Invoke(`A`)
	e2.Invoke(`B`)
	e3.Invoke(`C`)
	checkBuf(t, buf, `[1:A][2:A][1:B][2:B][1:C][2:C]`)

	checkEqual(t, true, lis1.Unsubscribe(e1))
	checkEqual(t, true, lis1.Unsubscribe(e2))
	checkEqual(t, false, lis1.Unsubscribe(e2))
	checkEqual(t, true, lis1.Subscribe(e1))
	checkBuf(t, buf, ``)

	e1.Invoke(`A`)
	e2.Invoke(`B`)
	e3.Invoke(`C`)
	checkBuf(t, buf, `[2:A][1:A][2:B][1:C][2:C]`)

	lis1.Cancel()
	e1.Invoke(`A`)
	e2.Invoke(`B`)
	e3.Invoke(`C`)
	checkBuf(t, buf, `[2:A][2:B][2:C]`)

	e1.Clear()
	e2.Clear()
	e1.Invoke(`A`)
	e2.Invoke(`B`)
	e3.Invoke(`C`)
	checkBuf(t, buf, `[2:C]`)
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
