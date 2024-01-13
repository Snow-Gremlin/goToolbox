package terror

import (
	"bytes"
	"sort"

	"github.com/Snow-Gremlin/goToolbox/internal/liteUtils"
	"github.com/Snow-Gremlin/goToolbox/terrors"
)

// ToString constructs a long detailed string for the whole
// tree of errors following this error.
func ToString(err error) string {
	w := newErrWriter()
	w.WriteError(err)
	return w.String()
}

type errWriter struct {
	buf *bytes.Buffer
}

func newErrWriter() errWriter {
	return errWriter{
		buf: &bytes.Buffer{},
	}
}

func (w errWriter) Write(text string) {
	_, _ = w.buf.WriteString(text)
}

func (w errWriter) WriteValue(value any) {
	w.Write(liteUtils.String(value))
}

func (w errWriter) WriteError(err error) {
	w.WriteMessage(err)
	w.WriteContext(err)
	w.WriteWrapped(err)
}

func (w errWriter) WriteMessage(err error) {
	if e, ok := err.(terrors.Messager); ok {
		w.Write(e.Message())
	} else {
		w.Write(err.Error())
	}
}

func (w errWriter) WriteContext(err error) {
	e, ok := err.(terrors.Contexture)
	if !ok {
		return
	}

	context := e.Context()
	if len(context) <= 0 {
		return
	}

	w.Write(` {`)
	keys := make([]string, len(context))
	index := 0
	for key := range context {
		keys[index] = key
		index++
	}
	sort.Strings(keys)
	for i, key := range keys {
		if i > 0 {
			w.Write(`, `)
		}
		w.Write(key)
		w.Write(`: `)
		w.WriteValue(context[key])
	}
	w.Write(`}`)
}

func (w errWriter) WriteWrapped(err error) {
	wrapped := Unwrap(err)
	count := len(wrapped)
	if count <= 0 {
		return
	}

	if count == 1 {
		w.Write(`: `)
		w.WriteError(wrapped[0])
		return
	}

	w.Write(`: [`)
	for i, e := range wrapped {
		if i > 0 {
			w.Write(`, `)
		}
		w.WriteError(e)
	}
	w.Write(`]`)
}

func (w errWriter) String() string {
	return w.buf.String()
}
