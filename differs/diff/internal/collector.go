package internal

import "github.com/Snow-Gremlin/goToolbox/differs"

// Collector is a tool for collecting the diff result.
type Collector interface {
	// InsertAdded inserts new Added parts into this collection.
	// This is expected to be inserted in reverse order from the expected result.
	InsertAdded(count int)

	// InsertRemoved inserts new Removed parts into this collection.
	// This is expected to be inserted in reverse order from the expected result.
	InsertRemoved(count int)

	// InsertEqual inserts new Equal parts into this collection.
	// This is expected to be inserted in reverse order from the expected result.
	InsertEqual(count int)

	// InsertSubstitute inserts new Added and Removed parts into this collection.
	// This is expected to be inserted in reverse order from the expected result.
	InsertSubstitute(count int)

	// Finish inserts any remaining parts which haven't been inserted yet.
	Finish() differs.Result
}
