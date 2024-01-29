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
	IQDeficit     float64 `json:"iq_deficit"` // idealized quadratic funding deficit
	PriorityScore float64 `json:"priority_score"`
}

func NewConcernState(id motionproto.MotionID) *ConcernState {
	return &ConcernState{
		PriorityPoll: pmp_1.ConcernPollBallotName(id),
	}
}

type PolicyState struct {
	WithheldEscrowFraction float64 `json:"withheld_escrow_fraction"`
	MatchDeficit           float64 `json:"match_deficit"`
}

var InitialPolicyState = &PolicyState{
	WithheldEscrowFraction: 0.1,
	MatchDeficit:           0.0,
}
