package optional

import (
	"goToolbox/terrors/terror"
	"goToolbox/utils"
)

func oneArg(values []int, default1 int, usage string) int {
	switch len(values) {
	case 0:
		return default1
	case 1:
		return values[0]
	default:
		panic(terror.InvalidArgCount(1, len(values), usage))
	}
}

func twoArgs(values []int, default1, default2 int, usage string) (int, int) {
	switch len(values) {
	case 0:
		return default1, default2
	case 1:
		return values[0], default2
	case 2:
		return values[0], values[1]
	default:
		panic(terror.InvalidArgCount(2, len(values), usage))
	}
}

// SizeAndCapacity deals with an optional size and capacity.
//
// This may have zero, one, or two values.
// The first value is the size which is clamped to zero.
// The second value is the capacity which is clamped to the size.
// This will panic if more than two values.
func SizeAndCapacity(sizes []int) (int, int) {
	size, cap := twoArgs(sizes, 0, 0, `size and capacity`)
	return max(0, size), max(0, size, cap)
}

// Capacity deals with an optional capacity.
//
// This may have zero or one values.
// The first value is the capacity which is clamped to zero.
// This will panic if more than one value.
func Capacity(capacity []int) int {
	return max(0, oneArg(capacity, 0, `capacity`))
}

// Size deals with an optional size.
//
// This may have zero or one values.
// The first value is the size which is clamped to zero.
// This will panic if more than one value.
func Size(size []int) int {
	return max(0, oneArg(size, 0, `size`))
}

// After deals with an optional after index.
//
// This may have zero or one values.
// The first value is the after index which is clamped to -1.
// This will panic if more than one value.
func After(size []int) int {
	return max(-1, oneArg(size, -1, `after index`))
}

// Comparer deals with an optional comparer.
//
// This may have zero or one comparer.
// If there is no comparer, this will try to determine the default comparer for the given type.
// If nil is passed in for the comparer, this will act like no comparers were given.
// This will panic if more than one comparer is given, or no comparer was given
// and their is no default comparer for this type.
func Comparer[T any](comps []utils.Comparer[T]) utils.Comparer[T] {
	if count := len(comps); count > 0 {
		if count > 1 {
			panic(terror.InvalidArgCount(1, count, `comparer`))
		}
		if cmp := comps[0]; !utils.IsNil(cmp) {
			return cmp
		}
	}

	if cmp := utils.DefaultComparer[T](); !utils.IsNil(cmp) {
		return cmp
	}

	panic(terror.New(`must provide a comparer to compare this type`).
		With(`type`, utils.TypeOf[T]()))
}
