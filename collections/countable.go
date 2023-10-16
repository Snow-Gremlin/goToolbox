package collections

// Countable is an object which has a countable number of values in it.
type Countable interface {

	// Count is the number of values in this object.
	Count() int
}
