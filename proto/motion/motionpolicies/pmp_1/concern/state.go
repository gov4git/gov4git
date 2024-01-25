package concern

import (
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/motion/motionpolicies/pmp_1"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
)

const StateFilebase = "state.json"

type ConcernState struct {
	PriorityPoll      ballotproto.BallotID `json:"priority_poll"`
	EligibleProposals motionproto.Refs     `json:"eligible_proposals"`
	//
	LatestIQFDeficit    float64 `json:"iqf_deficit"` // idealized quadratic funding deficit
	LatestPriorityScore float64 `json:"latest_priority_score"`
}

func NewConcernState(id motionproto.MotionID) *ConcernState {
	return &ConcernState{
		PriorityPoll: pmp_1.ConcernPollBallotName(id),
	}
}
