package concern

import (
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/motion/motionpolicies/pmp_0"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
)

const StateFilebase = "state.json"

type ConcernState struct {
	PriorityPoll        ballotproto.BallotID `json:"priority_poll"`
	LatestPriorityScore float64              `json:"latest_priority_score"`
	EligibleProposals   motionproto.Refs     `json:"eligible_proposals"`
	CostOfPriority      float64              `json:"cost_of_priority"`
}

func NewConcernState(id motionproto.MotionID) *ConcernState {
	return &ConcernState{
		PriorityPoll: pmp_0.ConcernPollBallotName(id),
	}
}

func (x *ConcernState) ProjectedBounty() float64 {
	return x.CostOfPriority
}
