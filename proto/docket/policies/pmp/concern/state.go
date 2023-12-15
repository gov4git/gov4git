package concern

import (
	"github.com/gov4git/gov4git/v2/proto/ballot/common"
	"github.com/gov4git/gov4git/v2/proto/docket/policies/pmp"
	"github.com/gov4git/gov4git/v2/proto/docket/schema"
)

const StateFilebase = "state.json"

type ConcernState struct {
	PriorityPoll        common.BallotName `json:"priority_poll"`
	LatestPriorityScore float64           `json:"latest_priority_score"`
	EligibleProposals   schema.Refs       `json:"eligible_proposals"`
}

func NewConcernState(id schema.MotionID) *ConcernState {
	return &ConcernState{
		PriorityPoll: pmp.ConcernPollBallotName(id),
	}
}
