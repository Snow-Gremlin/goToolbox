package wagner

import "github.com/Snow-Gremlin/goToolbox/differs/diff/internal"

// walkerStep is the signature for steps which the walker can take.
type walkerStep func()

// walker is a structure for keeping track of a walk through
// the cost matrix for Wagner–Fischer.
type walker struct {
	// col is the collector to output the walk into.
	col internal.Collector

	// i is the current A offset
	i int

	// j is the current B offset
	j int
}

var (
	_ walkerStep = ((*walker)(nil)).moveA
	_ walkerStep = ((*walker)(nil)).moveB
	_ walkerStep = ((*walker)(nil)).moveEqual
	_ walkerStep = ((*walker)(nil)).moveSubstitute
)

// newWalker creates a new walker instance to help walk the cost matrix.
func newWalker(cont internal.Container, col internal.Collector) *walker {
	return &walker{
		col: col,
		i:   cont.ACount() - 1,
		j:   cont.BCount() - 1,
	}
}

// hasMore indicates there is more to walk.
func (w *walker) hasMore() bool {
	return w.i >= 0 && w.j >= 0
}

// moveA is a walker movement which steps only on the A input.
// This will insert a removed.
func (w *walker) moveA() {
	w.i--
	w.col.InsertRemoved(1)
}

// moveB is a walker movement which steps only on the B input.
// This will insert an added.
func (w *walker) moveB() {
	w.j--
	w.col.InsertAdded(1)
}

// moveEqual is a walker movement which steps both inputs.
// This will insert an equal.
func (w *walker) moveEqual() {
	w.i--
	w.j--
	w.col.InsertEqual(1)
}

// moveSubstitute is a walker movement which steps both inputs.
// This will insert a substitute.
func (w *walker) moveSubstitute() {
	w.i--
	w.j--
	w.col.InsertSubstitute(1)
}

// finish will be called when the walk is done to write any
// remaining adds and removes.
func (w *walker) finish() {
	w.col.InsertRemoved(w.i + 1)
	w.col.InsertAdded(w.j + 1)
}
