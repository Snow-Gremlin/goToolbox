package collector

import (
	"goToolbox/collections"
	"goToolbox/collections/capStack"
	"goToolbox/differs/diff/internal"
	"goToolbox/differs/step"
)

// New creates a new collector.
func New(aCount, bCount int) internal.Collector {
	return &collectorImp{
		stack:      capStack.New[collections.Tuple2[step.Step, int]](),
		aCount:     aCount,
		bCount:     bCount,
		total:      0,
		addedRun:   0,
		removedRun: 0,
		equalRun:   0,
	}
}
