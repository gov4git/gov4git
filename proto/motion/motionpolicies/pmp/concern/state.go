package concern

import (
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/motion/motionpolicies/pmp"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
)

const StateFilebase = "state.json"

type ConcernState struct {
	PriorityPoll        ballotproto.BallotID `json:"priority_poll"`
	LatestPriorityScore float64              `json:"latest_priority_score"`
	EligibleProposals   motionproto.Refs     `json:"eligible_proposals"`
}

func NewConcernState(id motionproto.MotionID) *ConcernState {
	return &ConcernState{
		PriorityPoll: pmp.ConcernPollBallotName(id),
	}
}
