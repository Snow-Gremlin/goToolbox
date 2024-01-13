package collector

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/capStack"
	"github.com/Snow-Gremlin/goToolbox/differs/diff/internal"
	"github.com/Snow-Gremlin/goToolbox/differs/step"
)

// New creates a new collector.
func New(aCount, bCount int) internal.Collector {
	return &collectorImp{
		stack:      capStack.New[collections.Tuple2[step.Step, int]](),
		aCount:     aCount,
		bCount:     bCount,
		total:      0,
		addCount:   0,
		remCount:   0,
		addedRun:   0,
		removedRun: 0,
		equalRun:   0,
	}
}
