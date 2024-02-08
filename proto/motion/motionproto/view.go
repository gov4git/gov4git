package motionproto

import (
	"sort"

	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/lib4git/form"
)

type MotionView struct {
	Motion  Motion                   `json:"motion"`
	Ballots MotionBallots            `json:"ballots"`
	Policy  form.Form                `json:"policy"`
	Voter   *ballotproto.VoterStatus `json:"voter_status,omitempty"`
}

func (mv MotionView) IsMissingPolicy() bool {
	return mv.Policy == nil
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
