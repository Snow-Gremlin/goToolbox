package internal

const (
	// RemoveCost gives the cost to remove A at the given index.
	RemoveCost = 1

	// AddCost gives the cost to add B at the given index.
	AddCost = 1

	// SubstitutionCost gives the substitution cost for replacing A with B at the given indices.
	SubstitutionCost = 2

	// EqualCost gives the cost for A and B being equal.
	EqualCost = 0
)
