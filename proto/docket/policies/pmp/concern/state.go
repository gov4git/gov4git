package concern

import (
	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/docket/policies/pmp"
	"github.com/gov4git/gov4git/proto/docket/schema"
)

const StateFilebase = "state.json"

type ConcernState struct {
	PriorityPoll common.BallotName `json:"priority_poll_ballot"`
}

func NewConcernState(id schema.MotionID) *ConcernState {
	return &ConcernState{
		PriorityPoll: pmp.ConcernPollBallotName(id),
	}
}
