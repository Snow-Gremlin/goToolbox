package data

type imp struct {
	aCount int
	bCount int
	equals func(aIndex, bIndex int) bool
}

func (i *imp) ACount() int {
	return i.aCount
}

func (i *imp) BCount() int {
	return i.bCount
}

func (i *imp) Equals(aIndex, bIndex int) bool {
	return i.equals(aIndex, bIndex)
}
