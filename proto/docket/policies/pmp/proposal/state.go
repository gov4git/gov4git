package proposal

import (
	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/docket/policies/pmp"
	"github.com/gov4git/gov4git/proto/docket/schema"
)

const StateFilebase = "state.json"

type ProposalState struct {
	ApprovalPoll        common.BallotName `json:"approval_poll"`
	LatestApprovalScore float64           `json:"latest_approval_score"`
	ResolvingConcerns   schema.Refs       `json:"resolving_concerns"`
}

func NewProposalState(id schema.MotionID) *ProposalState {
	return &ProposalState{
		ApprovalPoll: pmp.ProposalApprovalPollName(id),
	}
}
