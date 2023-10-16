package hirschberg

import "goToolbox/differs/diff/internal"

// New creates a new Hirschberg diff algorithm
// (https://en.wikipedia.org/wiki/Hirschberg%27s_algorithm)
//
// This allows for an optional diff (usually Wagner) to use when possible to hybrid the algorithm,
// to not use the optional diff pass in nil. The hybrid is used if it has enough memory preallocated,
// NoResizeNeeded returns true, otherwise Hirschberg will continue to divide the space until
// the hybrid can be used without causing it to reallocate memory.
//
// The given length is the initial score vector size. If the vector is too small it will be
// reallocated to the larger size. Use -1 to not preallocate the vectors.
//
// The useReduce flag indicates if the equal padding edges should be checked
// at each step of the algorithm or not.
func New(hybrid internal.Algorithm, length int, useReduce bool) internal.Algorithm {
	return &hirschbergImp{
		scores:    newScores(length),
		hybrid:    hybrid,
		useReduce: useReduce,
	}
}
