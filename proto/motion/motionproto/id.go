package motionproto

import (
	"sort"
)

type MotionID string

func (x MotionID) String() string {
	return string(x)
}

type MotionIDs []MotionID

func (x MotionIDs) Len() int           { return len(x) }
func (x MotionIDs) Less(i, j int) bool { return x[i] < x[j] }
func (x MotionIDs) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
func (x MotionIDs) Sort()              { sort.Sort(x) }

type MotionIDSet map[MotionID]bool

func (x MotionIDSet) Add(id MotionID) {
	x[id] = true
}

func (x MotionIDSet) MotionIDs() MotionIDs {
	s := make(MotionIDs, 0, len(x))
	for id := range x {
		s = append(s, id)
	}
	s.Sort()
	return s
}
