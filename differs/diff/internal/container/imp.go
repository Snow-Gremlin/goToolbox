package container

import (
	"goToolbox/differs"
	"goToolbox/differs/diff/internal"
)

func newSub(data differs.Data, aOffset, aCount, bOffset, bCount int, reverse bool) internal.Container {
	return &containerImp{
		data:    data,
		aOffset: aOffset,
		aCount:  aCount,
		bOffset: bOffset,
		bCount:  bCount,
		reverse: reverse,
	}
}

type containerImp struct {
	data    differs.Data
	aOffset int
	aCount  int
	bOffset int
	bCount  int
	reverse bool
}

func (cont *containerImp) ACount() int {
	return cont.aCount
}

func (cont *containerImp) BCount() int {
	return cont.bCount
}

func (cont *containerImp) Equals(aIndex, bIndex int) bool {
	if cont.reverse {
		return cont.data.Equals(
			cont.aCount-1-aIndex+cont.aOffset,
			cont.bCount-1-bIndex+cont.bOffset)
	}

	return cont.data.Equals(
		aIndex+cont.aOffset,
		bIndex+cont.bOffset)
}

func (cont *containerImp) SubstitutionCost(i, j int) int {
	if cont.Equals(i, j) {
		return internal.EqualCost
	}
	return internal.SubstitutionCost
}

func (cont *containerImp) Sub(aLow, aHigh, bLow, bHigh int, reverse bool) internal.Container {
	if cont.reverse {
		return newSub(cont.data,
			cont.aCount-aHigh+cont.aOffset, aHigh-aLow,
			cont.bCount-bHigh+cont.bOffset, bHigh-bLow,
			!reverse)
	}

	return newSub(cont.data,
		aLow+cont.aOffset, aHigh-aLow,
		bLow+cont.bOffset, bHigh-bLow,
		reverse)
}

func (cont *containerImp) Reduce() (internal.Container, int, int) {
	var before, after, i, j, width int

	width = min(cont.aCount, cont.bCount)
	for before, i, j = 0, cont.aOffset, cont.bOffset; before < width; before, i, j = before+1, i+1, j+1 {
		if !cont.data.Equals(i, j) {
			break
		}
	}

	width = width - before
	for after, i, j = 0, cont.aCount-1+cont.aOffset, cont.bCount-1+cont.bOffset; after < width; after, i, j = after+1, i-1, j-1 {
		if !cont.data.Equals(i, j) {
			break
		}
	}

	sub := newSub(cont.data,
		before+cont.aOffset, cont.aCount-after-before,
		before+cont.bOffset, cont.bCount-after-before,
		cont.reverse)

	if cont.reverse {
		return sub, after, before
	}
	return sub, before, after
}

func (cont *containerImp) EndCase(col internal.Collector) bool {
	if cont.aCount <= 1 {
		cont.aEdge(col)
		return true
	}

	if cont.bCount <= 1 {
		cont.bEdge(col)
		return true
	}

	return false
}

// aEdge handles when at the edge of the A source subset in the given container.
func (cont *containerImp) aEdge(col internal.Collector) {
	aLen, bLen := cont.aCount, cont.bCount

	if aLen <= 0 {
		col.InsertAdded(bLen)
		return
	}

	split := -1
	if cont.reverse {
		iRaw := cont.aCount - 1 + cont.aOffset
		jOffset := cont.bCount - 1 + cont.bOffset
		for j := 0; j < bLen; j++ {
			if cont.data.Equals(iRaw, jOffset-j) {
				split = j
				break
			}
		}
	} else {
		iRaw := cont.aOffset
		jOffset := cont.bOffset
		for j := 0; j < bLen; j++ {
			if cont.data.Equals(iRaw, j+jOffset) {
				split = j
				break
			}
		}
	}

	if split < 0 {
		col.InsertAdded(bLen)
		col.InsertRemoved(1)
	} else {
		col.InsertAdded(bLen - split - 1)
		col.InsertEqual(1)
		col.InsertAdded(split)
	}
}

// bEdge Handles when at the edge of the B source subset in the given container.
func (cont *containerImp) bEdge(col internal.Collector) {
	aLen, bLen := cont.aCount, cont.bCount

	if bLen <= 0 {
		col.InsertRemoved(aLen)
		return
	}

	split := -1
	if cont.reverse {
		iOffset := cont.aCount - 1 + cont.aOffset
		jRaw := cont.bCount - 1 + cont.bOffset
		for i := 0; i < aLen; i++ {
			if cont.data.Equals(iOffset-i, jRaw) {
				split = i
				break
			}
		}
	} else {
		iOffset := cont.aOffset
		jRaw := cont.bOffset
		for i := 0; i < aLen; i++ {
			if cont.data.Equals(i+iOffset, jRaw) {
				split = i
				break
			}
		}
	}

	if split < 0 {
		col.InsertAdded(1)
		col.InsertRemoved(aLen)
	} else {
		col.InsertRemoved(aLen - split - 1)
		col.InsertEqual(1)
		col.InsertRemoved(split)
	}
}
