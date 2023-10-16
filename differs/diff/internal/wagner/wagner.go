package wagner

import "goToolbox/differs/diff/internal"

// New creates a new Wagner–Fischer diff algorithm.
// (https://en.wikipedia.org/wiki/Wagner%E2%80%93Fischer_algorithm).
//
// The given size is the amount of matrix space, width * height, to preallocate
// for the Wagner-Fischer algorithm. Use zero or less to not preallocate any matrix.
func New(size int) internal.Algorithm {
	w := &wagnerImp{}
	if size > 0 {
		w.allocateMatrix(size)
	}
	return w
}
