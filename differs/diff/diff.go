package diff

import (
	"goToolbox/differs"
	"goToolbox/differs/diff/internal/hirschberg"
	"goToolbox/differs/diff/internal/wagner"
)

// DefaultWagnerThreshold is the point at which the algorithms switch from Hirschberg
// to Wagner-Fischer. When both length of the comparable are smaller than this value
// Wagner-Fischer is used. The Wagner matrix will never be larger than this value of entries.
// If this is less than 4 the Wagner algorithm will not be used.
const DefaultWagnerThreshold = 500

// Hirschberg creates a new Hirschberg algorithm instance for performing a diff.
//
// The given length is the initial score vector size. If the vector is too small it will be
// reallocated to the larger size. Use -1 to not preallocate the vectors.
// The useReduce flag indicates if the equal padding edges should be checked
// at each step of the algorithm or not.
func Hirschberg(length int, useReduce bool) differs.Diff {
	return wrap(hirschberg.New(nil, length, useReduce))
}

// Wagner creates a new Wagner-Fischer algorithm instance for performing a diff.
//
// The given size is the amount of matrix space, width * height, to preallocate
// for the Wagner-Fischer algorithm. Use -1 to not preallocate any matrix.
func Wagner(size int) differs.Diff {
	return wrap(wagner.New(size))
}

// Hybrid creates a new hybrid Hirschberg with Wagner-Fischer cutoff for performing a diff.
//
// The given length is the initial score vector size of the Hirschberg algorithm. If the vector
// is too small it will be reallocated to the larger size. Use -1 to not preallocate the vectors.
// The useReduce flag indicates if the equal padding edges should be checked
// at each step of the algorithm or not.
//
// The given size is the amount of matrix space, width * height, to use for the Wagner-Fischer.
// This must be greater than 4 fo use the cutoff. The larger the size, the more memory is used
// creating the matrix but the earlier the Wagner-Fischer algorithm can take over.
func Hybrid(length int, useReduce bool, size int) differs.Diff {
	return wrap(hirschberg.New(wagner.New(size), length, useReduce))
}

// Default creates the default diff algorithm with default configuration.
//
// The default is a hybrid Hirschberg with Wagner-Fischer using a reduction
// at each step and the default Wagner threshold.
func Default() differs.Diff {
	return Hybrid(-1, true, DefaultWagnerThreshold)
}
