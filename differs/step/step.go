package step

// Step indicates a step in the Levenshtein Path.
type Step int

const (
	// Equal indicates A and B entries are equal.
	Equal Step = iota

	// Added indicates A was added.
	// Meaning the value is in A but not in B so advance A.
	Added

	// Removed indicates A was removed.
	// Meaning the value is in B but not in A so advance B.
	Removed
)

// String gets the string for step type.
func (s Step) String() string {
	switch s {
	case Equal:
		return `=`
	case Added:
		return `+`
	case Removed:
		return `-`
	default:
		return `?`
	}
}
