package waimea

import (
	"slices"

	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
)

type ProposalState struct {
	ApprovalPoll            ballotproto.BallotID `json:"approval_poll"`
	ApprovalScore           float64              `json:"approval_score"`
	ReviewMatch             float64              `json:"review_match"` // copied from policy state
	CostOfReview            float64              `json:"cost_of_review"`
	ProjectedPriorityBounty float64              `json:"projected_priority_bounty"`
	Decision                motionproto.Decision `json:"decision,omitempty"` // set on close or cancel, to be picked up by clearance pass
	EligibleConcerns        motionproto.Refs     `json:"eligible_concerns"`
}

func (x *ProposalState) Copy() *ProposalState {
	z := *x
	z.EligibleConcerns = slices.Clone(x.EligibleConcerns)
	return &z
}

func (x *ProposalState) ProjectedApprovalBounty() float64 {
	if x.ApprovalScore < 0 {
		return 0
	}
	return x.ApprovalScore * x.ReviewMatch
}

func NewProposalState(id motionproto.MotionID, reviewMatch float64) *ProposalState {
	return &ProposalState{
		ApprovalPoll: ProposalApprovalPollName(id),
		ReviewMatch:  reviewMatch,
	}
}
