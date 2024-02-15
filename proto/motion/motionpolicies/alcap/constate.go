package alcap

import (
	"slices"

	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
)

const StateFilebase = "state.json"

type ConcernState struct {
	PriorityPoll      ballotproto.BallotID `json:"priority_poll"`
	EligibleProposals motionproto.Refs     `json:"eligible_proposals"`
	//
	PriorityScore float64 `json:"priority_score"`
}

func (x *ConcernState) Copy() *ConcernState {
	z := *x
	z.EligibleProposals = slices.Clone(x.EligibleProposals)
	return &z
}

func (x *ConcernState) ProjectedBounty() float64 {
	if x.PriorityScore < 0 {
		return 0
	}
	return x.PriorityScore // XXX: times match fraction
}

func NewConcernState(id motionproto.MotionID) *ConcernState {
	return &ConcernState{
		PriorityPoll: ConcernPollBallotName(id),
	}
}
