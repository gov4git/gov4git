package pmp_1

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
	IQDeficit     float64 `json:"iq_deficit"`     // idealized quadratic funding deficit
	PriorityScore float64 `json:"priority_score"` // is "escrow"
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
	return x.PriorityScore
}

func NewConcernState(id motionproto.MotionID) *ConcernState {
	return &ConcernState{
		PriorityPoll: ConcernPollBallotName(id),
	}
}

type ConcernPolicyState struct {
	WithheldEscrowFraction float64 `json:"withheld_escrow_fraction"`
	MatchDeficit           float64 `json:"match_deficit"`
}

var InitialPolicyState = &ConcernPolicyState{
	WithheldEscrowFraction: 0.1,
	MatchDeficit:           0.0,
}
