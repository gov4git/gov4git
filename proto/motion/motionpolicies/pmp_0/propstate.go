package pmp_0

import (
	"slices"

	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
)

type ProposalState struct {
	ApprovalPoll        ballotproto.BallotID `json:"approval_poll"`
	LatestApprovalScore float64              `json:"latest_approval_score"`
	EligibleConcerns    motionproto.Refs     `json:"eligible_concerns"`
}

func NewProposalState(id motionproto.MotionID) *ProposalState {
	return &ProposalState{
		ApprovalPoll: ProposalApprovalPollName(id),
	}
}

func (x *ProposalState) Copy() *ProposalState {
	q := *x
	q.EligibleConcerns = slices.Clone(x.EligibleConcerns)
	return &q
}
