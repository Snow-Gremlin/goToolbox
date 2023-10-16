package internal

import "goToolbox/differs"

// Container is a wrapper for the data.
// It is used to determine subset and revered reading of the data.
type Container interface {
	differs.Data

	// SubstitutionCost determines the substitution cost for the given indices.
	SubstitutionCost(i, j int) int

	// Sub creates a new data container for a subset and
	// reverse relative to this container's settings.
	// The high values are exclusive, the low is inclusive.
	Sub(aLow, aHigh, bLow, bHigh int, reverse bool) Container

	// Reduce determines how much of the edges of this container are equal.
	// The amount before and after which are equal are returned and
	// the reduced sub-container is returned.
	Reduce() (Container, int, int)

	// EndCase determines if the given container is small enough to be simply added
	// into the collector without any diff algorithm. This will add into the given
	// collector and return true if done, otherwise it will return false.
	EndCase(col Collector) bool
}
