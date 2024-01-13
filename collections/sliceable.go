package collections

// Sliceable is an object which can get the data as a slice.
type Sliceable[T any] interface {
	// ToSlice returns the values as a slice.
	ToSlice() []T

	// CopyToSlice copies as much of the values fit into the given slice.
	CopyToSlice(s []T)
}
