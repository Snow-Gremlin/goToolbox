package wagner

import "goToolbox/differs/diff/internal"

type wagnerImp struct {
	costs []int
}

func (w *wagnerImp) allocateMatrix(size int) {
	w.costs = make([]int, size)
}

func (w *wagnerImp) NoResizeNeeded(cont internal.Container) bool {
	return len(w.costs) >= cont.ACount()*cont.BCount()
}

func (w *wagnerImp) Diff(cont internal.Container, col internal.Collector) {
	if size := cont.ACount() * cont.BCount(); len(w.costs) < size {
		w.allocateMatrix(size)
	}
	w.setCosts(cont)
	w.walkPath(cont, col)
}

// setCosts will populate the part of the cost matrix which is needed by the given container.
// The costs are based off of the equality of parts in the comparable in the given container.
func (w *wagnerImp) setCosts(cont internal.Container) {
	aLen := cont.ACount()
	bLen := cont.BCount()

	start := cont.SubstitutionCost(0, 0)
	w.costs[0] = start

	for i, value := 1, start; i < aLen; i++ {
		value = min(value+1,
			i+cont.SubstitutionCost(i, 0))
		w.costs[i] = value
	}

	for j, k, value := 1, aLen, start; j < bLen; j, k = j+1, k+aLen {
		value = min(value+1,
			j+cont.SubstitutionCost(0, j))
		w.costs[k] = value
	}

	for j, k, k2, k3 := 1, aLen+1, 1, 0; j < bLen; j, k, k2, k3 = j+1, k+1, k2+1, k3+1 {
		for i, value := 1, w.costs[k-1]; i < aLen; i, k, k2, k3 = i+1, k+1, k2+1, k3+1 {
			value = min(value+1,
				w.costs[k2]+1,
				w.costs[k3]+cont.SubstitutionCost(i, j))
			w.costs[k] = value
		}
	}
}

// getCost gets the cost value at the given indices.
// If the indices are out-of-bounds the edge cost will be returned.
func (w *wagnerImp) getCost(i, j, aLen int) int {
	switch {
	case i < 0:
		return j + 1
	case j < 0:
		return i + 1
	default:
		return w.costs[i+j*aLen]
	}
}

// walkPath will walk through the cost matrix backwards to find the minimum Levenshtein path.
// The steps for this path are added to the given collector.
func (w *wagnerImp) walkPath(cont internal.Container, col internal.Collector) {
	aLen := cont.ACount()
	walk := newWalker(cont, col)
	for walk.hasMore() {
		aCost := w.getCost(walk.i-1, walk.j, aLen)
		bCost := w.getCost(walk.i, walk.j-1, aLen)
		cCost := w.getCost(walk.i-1, walk.j-1, aLen)
		minCost := min(aCost, bCost, cCost)

		var curMove walkerStep
		if aCost == minCost {
			curMove = walk.moveA
		}
		if bCost == minCost {
			curMove = walk.moveB
		}
		if cCost == minCost {
			if cont.Equals(walk.i, walk.j) {
				curMove = walk.moveEqual
			} else if curMove == nil {
				curMove = walk.moveSubstitute
			}
		}

		curMove()
	}
	walk.finish()
}
