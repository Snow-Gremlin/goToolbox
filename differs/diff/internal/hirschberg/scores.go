package hirschberg

import "github.com/Snow-Gremlin/goToolbox/differs/diff/internal"

// scores is the Hirschberg scores used for diffing two comparable sources.
type scores struct {
	// front is the score vector at the front of the score calculation.
	front []int

	// back is the score vector at the back of the score calculation.
	back []int

	// other is the score vector to store off a result vector to.
	other []int
}

// newScores creates a new path builder. The given length must be one greater
// than the maximum B length that will be passed into these scores.
func newScores(length int) *scores {
	s := &scores{
		front: nil,
		back:  nil,
		other: nil,
	}
	if length > 0 {
		s.allocateVectors(length)
	}
	return s
}

// allocateVectors will create the slices used for the score vectors.
func (s *scores) allocateVectors(length int) {
	s.front = make([]int, length)
	s.back = make([]int, length)
	s.other = make([]int, length)
}

// swap swaps the front and back score vectors.
func (s *scores) swap() {
	s.back, s.front = s.front, s.back
}

// store swaps the back and other score vectors.
func (s *scores) store() {
	s.back, s.other = s.other, s.back
}

// calculate calculates the Needleman-Wunsch score.
// At the end of this calculation the score is in the back vector.
func (s *scores) calculate(cont internal.Container) {
	aLen := cont.ACount()
	bLen := cont.BCount()
	if len(s.back) < bLen+1 {
		s.allocateVectors(bLen + 1)
	}

	s.back[0] = 0
	for j := 1; j <= bLen; j++ {
		s.back[j] = s.back[j-1] + internal.AddCost
	}

	for i := 1; i <= aLen; i++ {
		s.front[0] = s.back[0] + internal.RemoveCost
		for j := 1; j <= bLen; j++ {
			s.front[j] = min(
				s.back[j-1]+cont.SubstitutionCost(i-1, j-1),
				s.back[j]+internal.RemoveCost,
				s.front[j-1]+internal.AddCost)
		}

		s.swap()
	}
}

// findPivot finds the pivot between the other score and the reverse of the back score.
// The pivot is the index of the maximum sum of each element in the two scores.
func (s *scores) findPivot(bLength int) int {
	index := 0
	minValue := s.other[0] + s.back[bLength]
	for j := 1; j <= bLength; j++ {
		value := s.other[j] + s.back[bLength-j]
		if value < minValue {
			minValue = value
			index = j
		}
	}
	return index
}

// Split will find the A and B mid points to split the container at.
func (s *scores) Split(cont internal.Container) (int, int) {
	aLen := cont.ACount()
	bLen := cont.BCount()

	aMid := aLen / 2
	s.calculate(cont.Sub(0, aMid, 0, bLen, false))
	s.store()
	s.calculate(cont.Sub(aMid, aLen, 0, bLen, true))
	bMid := s.findPivot(bLen)

	return aMid, bMid
}
