package data

import (
	"regexp"
	"slices"

	"goToolbox/differs"
	"goToolbox/utils"
)

// New creates a comparer for comparing two collections of data.
//
// The given counts are the size of the A and B data. The given equals
// function can be called for any index between zero and the count.
// `aIndex` is the index in the A data between zero and `aCount`.
// `bIndex` is the index in the B data between zero and `bCount`.
// The index pairs may be called in any order and multiple times.
func New(aCount, bCount int, equals func(aIndex, bIndex int) bool) *imp {
	return &imp{
		aCount: aCount,
		bCount: bCount,
		equals: equals,
	}
}

// Comparable creates data for comparing comparable data types in two slices.
//
// This can be used for int64, uint64, or any other types that are comparable.
func Comparable[T comparable, S ~[]T](a, b S) differs.Data {
	return New(len(a), len(b), func(aIndex, bIndex int) bool {
		return a[aIndex] == b[bIndex]
	})
}

// Strings creates data for comparing strings in two slices of strings.
//
// This is they typical data used when since typically each string is
// a single line and the slices represent several lines in a whole document.
func Strings(a, b []string) differs.Data {
	return Comparable(a, b)
}

// Ints creates data for comparing integers in two slices.
func Ints(a, b []int) differs.Data {
	return Comparable(a, b)
}

// Bytes creates data for comparing bytes in two slices.
func Bytes(a, b []byte) differs.Data {
	return Comparable(a, b)
}

// Runes creates data for comparing runes in two rune slices.
func Runes(a, b []rune) differs.Data {
	return Comparable(a, b)
}

// Chars creates data for comparing characters in two strings.
//
// Warning: This does not handle escaped utf-8 sequences as a unit,
// this uses Go's default length and indexing for strings.
// For full unicode support use the Runes instead.
func Chars(a, b string) differs.Data {
	return New(len(a), len(b), func(aIndex, bIndex int) bool {
		return a[aIndex] == b[bIndex]
	})
}

// RuneSlice creates data for comparing rune slices in two slices of slices.
//
// This is similar to Strings except functions for unicode text.
func RuneSlice(a, b [][]rune) differs.Data {
	return New(len(a), len(b), func(aIndex, bIndex int) bool {
		return slices.Equal(a[aIndex], b[bIndex])
	})
}

// Any creates data for comparing any type in two slices using utils.Equal.
//
// Since utils.Equal may perform a reflection's deep equal, this may be
// slower than more specific data comparisons.
func Any[T any, S ~[]T](a, b S) differs.Data {
	return New(len(a), len(b), func(aIndex, bIndex int) bool {
		return utils.Equal(a[aIndex], b[bIndex])
	})
}

// Regex creates data for comparing a collection of parts against
// a collection of regular expressions.
//
// This is designed for checking a set of results which have some
// variability against expected patterns for the result.
func Regex[T any](patterns []string, parts []T) differs.Data {
	regs := make([]*regexp.Regexp, len(patterns))
	for i, pattern := range patterns {
		regs[i] = regexp.MustCompile(pattern)
	}
	lines := utils.Strings(parts)
	return New(len(regs), len(lines), func(regIndex, lineIndex int) bool {
		return regs[regIndex].MatchString(lines[lineIndex])
	})
}
