package schema

import (
	"sort"

	"github.com/gov4git/lib4git/form"
)

type MotionView struct {
	Motion Motion   `json:"motion"`
	Policy form.Map `json:"policy"`
}

type MotionViews []MotionView

func (x MotionViews) Len() int {
	return len(x)
}

func (x MotionViews) Less(i, j int) bool {
	return x[i].Motion.ID < x[j].Motion.ID
}

func (x MotionViews) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

func (x MotionViews) Sort() {
	sort.Sort(x)
}
