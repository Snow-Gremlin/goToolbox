package collector

import (
	"goToolbox/collections"
	"goToolbox/collections/tuple2"
	"goToolbox/differs"
	"goToolbox/differs/diff/internal/result"
	"goToolbox/differs/step"
)

type collectorImp struct {
	stack collections.Stack[collections.Tuple2[step.Step, int]]

	// aCount is the count of A values.
	aCount int

	// bCount is the count of B values.
	bCount int

	// total is the total number of parts represented by this collection.
	// The total sum of all the counts in each step.
	total int

	// addCount is the total number of added steps.
	addCount int

	// remCount is the total number of removed steps.
	remCount int

	// addedRun is the current amount of consecutive Added parts.
	addedRun int

	// removeRun is the current amount of consecutive Removed parts.
	removedRun int

	// equalRun is the current amount of consecutive Equal parts.
	equalRun int
}

// push pushes a new step into the collection.
func (c *collectorImp) push(step step.Step, count int) {
	c.stack.Push(tuple2.New(step, count))
	c.total += count
}

// pushAdd pushes an Added step if there is any Added parts currently collected.
func (c *collectorImp) pushAdded() {
	if c.addedRun > 0 {
		c.push(step.Added, c.addedRun)
		c.addCount += c.addedRun
		c.addedRun = 0
	}
}

// pushRemove pushes an Removed step if there is any Removed parts currently collected.
func (c *collectorImp) pushRemoved() {
	if c.removedRun > 0 {
		c.push(step.Removed, c.removedRun)
		c.remCount += c.removedRun
		c.removedRun = 0
	}
}

// pushEqual pushes an Add step if there is any Add parts currently collected.
func (c *collectorImp) pushEqual() {
	if c.equalRun > 0 {
		c.push(step.Equal, c.equalRun)
		c.equalRun = 0
	}
}

func (c *collectorImp) InsertAdded(count int) {
	if count > 0 {
		c.pushEqual()
		c.addedRun += count
	}
}

func (c *collectorImp) InsertRemoved(count int) {
	if count > 0 {
		c.pushEqual()
		c.removedRun += count
	}
}

func (c *collectorImp) InsertEqual(count int) {
	if count > 0 {
		c.pushAdded()
		c.pushRemoved()
		c.equalRun += count
	}
}

func (c *collectorImp) InsertSubstitute(count int) {
	if count > 0 {
		c.pushEqual()
		c.addedRun += count
		c.removedRun += count
	}
}

func (c *collectorImp) Finish() differs.Result {
	c.pushAdded()
	c.pushRemoved()
	c.pushEqual()
	return result.New(c.stack.Enumerate(), c.stack.Count(), c.aCount, c.bCount, c.total, c.addCount, c.remCount)
}
