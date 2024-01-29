package proposal

import (
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/motion/motionpolicies/pmp_1"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
)

const StateFilebase = "state.json"

type ProposalState struct {
	ApprovalPoll        ballotproto.BallotID `json:"approval_poll"`
	LatestApprovalScore float64              `json:"latest_approval_score"`
	EligibleConcerns    motionproto.Refs     `json:"eligible_concerns"`
	CostMultiplier      float64              `json:"cost_multiplier"`
}

func NewProposalState(id motionproto.MotionID) *ProposalState {
	return &ProposalState{
		ApprovalPoll:   pmp_1.ProposalApprovalPollName(id),
		CostMultiplier: 1.0,
	}
}
