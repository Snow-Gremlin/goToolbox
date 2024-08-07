package comp

// Comparable is an object which can be compared against another object.
//
// The given T variable type is typically the type that is implementing this interface.
type Comparable[T any] interface {
	// CompareTo returns a comparison result of this object and the given object.
	//
	// The comparison results should be:
	// `< 0` if this is less than the given other,
	// `== 0` if this equals the given other,
	// `> 0` if this is greater than the given other.
	CompareTo(other T) int
}
