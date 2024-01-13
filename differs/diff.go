package differs

// Diff is an instance of a Levenshtein difference algorithm.
type Diff interface {
	// Diff performs a diff on the given data and
	// returns the resulting diff path.
	Diff(data Data) Result

	// PlusMinus gets the labelled difference between the two slices.
	// It formats the results by prepending a "+" to new strings in [b],
	// a "-" for any to removed strings from [a], and " " if the strings are the same.
	PlusMinus(a, b []string) []string

	// Merge gets the labelled difference between the two slices
	// using a similar output to the git merge differences output.
	Merge(a, b []string) []string
}
