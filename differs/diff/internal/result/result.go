package result

import (
	"goToolbox/collections"
	"goToolbox/differs"
	"goToolbox/differs/step"
)

// New creates a new result from a diff.
//
//   - The given count must be the same count as the number of values in the enumerator.
//   - The given aCount must be the number of values in the A count.
//     The A count is also the sum of the added and same counts in the enumerator.
//   - The given bCount must be the number of values in the B count.
//     The B count is also the sum of the removed and same counts in the enumerator.
//   - The total must be the sum of all the counts in the enumerator.
//   - The given addCount is the sum of added counts in the enumerator.
//   - The given remCount is the sum of removed counts in the enumerator.
func New(enumerator collections.Enumerator[collections.Tuple2[step.Step, int]], count, aCount, bCount, total, addCount, remCount int) differs.Result {
	return &resultImp{
		enumerator: enumerator,
		count:      count,
		aCount:     aCount,
		bCount:     bCount,
		total:      total,
		addCount:   addCount,
		remCount:   remCount,
	}
}
