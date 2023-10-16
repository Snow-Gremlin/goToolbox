package collections

// Getter is an object which can get data at a given index or key.
//
// In other languages this would be an indexer. For example in C#
// this would be `public T this[int i]{ get; }`, but with Go's flair
// of allowing one or two returns, like when reading from a map in Go.
type Getter[TIn, TOut any] interface {

	// Get gets a value at the given index or key.
	//
	// Typically, if the index is out-of-bounds, this will panic.
	// However, if this is getting a value with a key that doesn't exist
	// it will typically return the zero value.
	Get(index TIn) TOut

	// TryGet gets a value at the given index or key.
	//
	// If the key exists or the index is in bounds then the found value
	// is returned with a true, otherwise if the key doesn't exist or
	// the index is out-of-bounds then the zero value is returned with false.
	TryGet(index TIn) (TOut, bool)
}
