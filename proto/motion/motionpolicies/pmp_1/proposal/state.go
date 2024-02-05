package proposal

import (
	"slices"

	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/motion/motionpolicies/pmp_1"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
)

const StateFilebase = "state.json"

type ProposalState struct {
	ApprovalPoll          ballotproto.BallotID `json:"approval_poll"`
	LatestApprovalScore   float64              `json:"latest_approval_score"`
	EligibleConcerns      motionproto.Refs     `json:"eligible_concerns"`
	InverseCostMultiplier float64              `json:"inverse_cost_multiplier"`
	Decision              motionproto.Decision `json:"decision,omitempty"` // set on close or cancel, to be picked up by clearance pass
}

func (x *ProposalState) Copy() *ProposalState {
	z := *x
	z.EligibleConcerns = slices.Clone(x.EligibleConcerns)
	return &z
}

func NewProposalState(id motionproto.MotionID) *ProposalState {
	return &ProposalState{
		ApprovalPoll:          pmp_1.ProposalApprovalPollName(id),
		InverseCostMultiplier: 1.0,
	}
}
