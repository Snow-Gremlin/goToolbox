package collections

// Clippable is an interface for data which can have excess capacity removed.
//
// Not all clippable data structures will be able to have capacity.
type Clippable interface {
	// Clip removes any excess capacity.
	//
	// If there is no capacity, this will have no effect.
	Clip()
}
