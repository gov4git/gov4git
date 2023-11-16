package proposal

import (
	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/docket/policies/pmp"
	"github.com/gov4git/gov4git/proto/docket/schema"
)

const StateFilebase = "state.json"

type ProposalState struct {
	ApprovalReferendum common.BallotName `json:"approval_referendum_ballot"`
}

func NewProposalState(id schema.MotionID) *ProposalState {
	return &ProposalState{
		ApprovalReferendum: pmp.ProposalReferendumBallotName(id),
	}
}
