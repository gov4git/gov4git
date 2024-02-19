package waimea

import (
	"slices"

	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
)

type ProposalState struct {
	ApprovalPoll     ballotproto.BallotID `json:"approval_poll"`
	ApprovalScore    float64              `json:"approval_score"`
	EligibleConcerns motionproto.Refs     `json:"eligible_concerns"`
	Decision         motionproto.Decision `json:"decision,omitempty"` // set on close or cancel, to be picked up by clearance pass
}

func (x *ProposalState) Copy() *ProposalState {
	z := *x
	z.EligibleConcerns = slices.Clone(x.EligibleConcerns)
	return &z
}

func NewProposalState(id motionproto.MotionID) *ProposalState {
	return &ProposalState{
		ApprovalPoll: ProposalApprovalPollName(id),
	}
}
