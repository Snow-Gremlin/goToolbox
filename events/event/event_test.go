package event

import (
	"bytes"
	"strconv"
	"testing"

	"github.com/Snow-Gremlin/goToolbox/events"
	"github.com/Snow-Gremlin/goToolbox/internal/liteUtils"
)

func Test_Event(t *testing.T) {
	buf := &bytes.Buffer{}
	obs1 := pseudoIntObserver{name: `one`, buf: buf}
	obs2 := pseudoIntObserver{name: `two`, buf: buf}
	obs3 := pseudoIntObserver{name: `three`, buf: buf}

	e1 := New[int]()
	checkEqual(t, true, e1.Add(obs1))
	checkEqual(t, true, e1.Add(obs2))
	checkEqual(t, false, e1.Add(nil))
	checkBuf(t, buf, `[one:Joined][two:Joined]`)

	e2 := New[int]()
	checkEqual(t, true, e2.Add(obs2))
	checkEqual(t, true, e2.Add(obs3))
	checkEqual(t, false, e2.Add(obs3))
	checkBuf(t, buf, `[two:Joined][three:Joined]`)

	e1.Invoke(42)
	checkBuf(t, buf, `[one:42][two:42]`)

	e2.Invoke(36)
	checkBuf(t, buf, `[two:36][three:36]`)

	e1.Invoke(8)
	checkBuf(t, buf, `[one:8][two:8]`)

	checkEqual(t, true, e1.Remove(obs1))
	checkEqual(t, true, e1.Remove(obs2))
	checkBuf(t, buf, `[one:Unjoined][two:Unjoined]`)
	checkEqual(t, false, e1.Remove(obs1))
	checkEqual(t, false, e1.Remove(obs2))
	checkEqual(t, false, e1.Remove(obs3))
	checkEqual(t, false, e1.Remove(nil))
	checkBuf(t, buf, ``)

	e2.Clear()
	checkBuf(t, buf, `[two:Unjoined][three:Unjoined]`)
	checkEqual(t, false, e2.Remove(obs1))
	checkEqual(t, false, e2.Remove(obs2))
	checkEqual(t, false, e2.Remove(obs3))
	checkBuf(t, buf, ``)

	e3 := Empty[int]()
	checkEqual(t, false, e3.Add(obs1))
	checkEqual(t, false, e3.Remove(obs1))
	e3.Clear()
	e3.Invoke(97)
}

type pseudoIntObserver struct {
	name string
	buf  *bytes.Buffer
}

func (pil pseudoIntObserver) Update(value int) {
	_, _ = pil.buf.WriteString(`[` + pil.name + `:` + strconv.Itoa(value) + `]`)
}

func (pil pseudoIntObserver) Joined(event events.Event[int]) {
	_, _ = pil.buf.WriteString(`[` + pil.name + `:Joined]`)
}

func (pil pseudoIntObserver) Unjoined(event events.Event[int]) {
	_, _ = pil.buf.WriteString(`[` + pil.name + `:Unjoined]`)
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
