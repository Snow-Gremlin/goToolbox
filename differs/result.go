package differs

import (
	"goToolbox/collections"
	"goToolbox/differs/step"
)

// Result contains the steps to take to walk through the diff.
//
// The result is kept as step groups. A step group is a continuous run
// of the same type of step represented by the step type and the number
// of steps of that type.
type Result interface {

	// Count get the number of step groups.
	collections.Countable

	// ACount is the length of the A data.
	ACount() int

	// BCount is the length of the B data.
	BCount() int

	// Total gets the sum of all the steps in each group.
	Total() int

	// AddedCount gets the sum of all the added steps.
	AddedCount() int

	// RemovedCount gets the sum of all the removed steps.
	RemovedCount() int

	// HasDiff indicates if there were any differences found.
	HasDiff() bool

	// Enumerates the step groups where each group is the step type
	// and the number of steps to take of that type.
	collections.Enumerable[collections.Tuple2[step.Step, int]]
}
