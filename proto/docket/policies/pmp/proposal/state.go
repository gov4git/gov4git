package proposal

import (
	"github.com/gov4git/gov4git/v2/proto/ballot/common"
	"github.com/gov4git/gov4git/v2/proto/docket/policies/pmp"
	"github.com/gov4git/gov4git/v2/proto/docket/schema"
)

const StateFilebase = "state.json"

type ProposalState struct {
	ApprovalPoll        common.BallotName `json:"approval_poll"`
	LatestApprovalScore float64           `json:"latest_approval_score"`
	EligibleConcerns    schema.Refs       `json:"eligible_concerns"`
}

func NewProposalState(id schema.MotionID) *ProposalState {
	return &ProposalState{
		ApprovalPoll: pmp.ProposalApprovalPollName(id),
	}
}
