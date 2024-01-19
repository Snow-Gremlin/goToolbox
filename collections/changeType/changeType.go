package changeType

// Change is a type of change
type ChangeType string

const (
	// invalid is the default to use for an invalid change.
	invalid ChangeType = `invalid`

	// Added indicates that one or more values were added, enqueued, pushed,
	// or inserted into a collection.
	//
	// All the values in the collection prior to this change will still exist
	// in the collection after, however the collection will have additional
	// new values.
	Added ChangeType = `Added`

	// Removed indicates that one or more values were removed, dequeued,
	// popped, deleted, or cleared from a collection.
	//
	// All the values in the collection after this change existed prior
	// to the change, however there will be less values in the collection.
	Removed ChangeType = `Removed`

	// Replaced indicates that one or more values were replaced.
	//
	// Values may have been added or removed without replacement if
	// a collection of values is replaced with a different sized collection.
	Replaced ChangeType = `Replaced`
)

// String gets the string value of this change.
func (c ChangeType) String() string {
	switch c {
	case Added, Removed, Replaced:
		return string(c)
	}
	return string(invalid)
}

// Valid indicates if the current change is a valid value for a change.
func (c ChangeType) Valid() bool {
	switch c {
	case Added, Removed, Replaced:
		return true
	}
	return false
}
