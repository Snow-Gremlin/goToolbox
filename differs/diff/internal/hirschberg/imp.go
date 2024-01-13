package hirschberg

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/stack"
	"github.com/Snow-Gremlin/goToolbox/collections/tuple2"
	"github.com/Snow-Gremlin/goToolbox/differs/diff/internal"
)

type hirschbergImp struct {
	scores    *scores
	hybrid    internal.Algorithm
	useReduce bool
}

func (h *hirschbergImp) NoResizeNeeded(cont internal.Container) bool {
	return len(h.scores.back) >= cont.BCount()+1
}

func (h *hirschbergImp) Diff(cont internal.Container, col internal.Collector) {
	stack := stack.New[collections.Tuple2[internal.Container, int]]()
	stack.Push(tuple2.New(cont, 0))

	for !stack.Empty() {
		cur, remainder := stack.Pop().Values()
		col.InsertEqual(remainder)
		if cur == nil {
			continue
		}

		if h.useReduce {
			var before, after int
			cur, before, after = cur.Reduce()
			col.InsertEqual(after)
			stack.Push(tuple2.New[internal.Container](nil, before))
		}

		if cur.EndCase(col) {
			continue
		}

		if h.hybrid != nil && h.hybrid.NoResizeNeeded(cur) {
			h.hybrid.Diff(cur, col)
			continue
		}

		aLen, bLen := cur.ACount(), cur.BCount()
		aMid, bMid := h.scores.Split(cur)
		stack.Push(tuple2.New(cur.Sub(0, aMid, 0, bMid, false), 0))
		stack.Push(tuple2.New(cur.Sub(aMid, aLen, bMid, bLen, false), 0))
	}
}
