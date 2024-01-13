package differs

// Data is the interface used to determine a diff for a variety of data types.
type Data interface {
	// ACount is the length of the A data.
	ACount() int

	// BCount is the length of the B data.
	BCount() int

	// Equals determines if the data in A at the index `aIndex`
	// is equal to B at the index `bIndex`.
	Equals(aIndex, bIndex int) bool
}
