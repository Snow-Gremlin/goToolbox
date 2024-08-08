package comp

// Comparer returns the comparison result between the two given values.
//
// The comparison results should be:
// `< 0` if x is less than y,
// `== 0` if x equals y,
// `> 0` if x is greater than y.
type Comparer[T any] func(x, y T) int

// Reverse gets a comparer that performs the opposite comparison.
// This makes a typical ascending sort into a descending sort.
func (cmp Comparer[T]) Reverse() Comparer[T] {
	return Descender(cmp)
}

// Pend will wait to perform the given comparison until the returned
// method is called. This is designed to help with `Or`.
func (cmp Comparer[T]) Pend(x, y T) func() int {
	return func() int {
		return cmp(x, y)
	}
}

// ToLess gets an IsLessThan for this comparer.
func (cmp Comparer[T]) ToLess() IsLessThan[T] {
	return func(x, y T) bool {
		return cmp(x, y) < 0
	}
}
