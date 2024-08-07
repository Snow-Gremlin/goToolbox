package utils

// Ternary returns the `a` value if the given test is true,
// and the `b` value if the test is false.
//
// Since both the `a` and `b` values are evaluated prior to the ternary,
// this should only be used when neither takes long to compute or with functions.
//
// For example `value := Ternary(x, 1, -1)` or
//
//	   value := Ternary(x,
//		     func() int { return foo() - 1 },
//		     func() int { return bar()*6 + 2 },
//	   )()
func Ternary[T any](test bool, a, b T) T {
	if test {
		return a
	}
	return b
}

// Flip returns `[b, a]` if the test is true,
// otherwise the values are returned in the given order, `[a, b]`.
//
// Example: `max, min := Flip(x < y, x, y)`
func Flip[T any](test bool, a, b T) (T, T) {
	if test {
		return b, a
	}
	return a, b
}
