package result

import (
	"fmt"

	"goToolbox/collections"
	"goToolbox/collections/enumerator"
	"goToolbox/differs/step"
)

type resultImp struct {
	enumerator collections.Enumerator[collections.Tuple2[step.Step, int]]
	count      int
	aCount     int
	bCount     int
	total      int
	addCount   int
	remCount   int
}

func (i *resultImp) Enumerate() collections.Enumerator[collections.Tuple2[step.Step, int]] {
	return i.enumerator
}

func (i *resultImp) Count() int {
	return i.count
}

func (i *resultImp) ACount() int {
	return i.aCount
}

func (i *resultImp) BCount() int {
	return i.bCount
}

func (i *resultImp) Total() int {
	return i.total
}

func (i *resultImp) AddedCount() int {
	return i.addCount
}

func (i *resultImp) RemovedCount() int {
	return i.remCount
}

func (i *resultImp) HasDiff() bool {
	return i.addCount > 0 || i.remCount > 0
}

func (i *resultImp) String() string {
	return enumerator.Select(i.Enumerate(), func(v collections.Tuple2[step.Step, int]) string {
		return fmt.Sprintf(`%s%d`, v.Value1().String(), v.Value2())
	}).Join(` `)
}
