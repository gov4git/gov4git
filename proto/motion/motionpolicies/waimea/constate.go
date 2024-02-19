package waimea

import (
	"slices"

	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
)

const StateFilebase = "state.json"

type ConcernState struct {
	PriorityPoll        ballotproto.BallotID `json:"priority_poll"`
	EligibleProposals   motionproto.Refs     `json:"eligible_proposals"`
	CostOfPriority      float64              `json:"cost_of_priority"`
	PriorityScore       float64              `json:"priority_score"`
	CostOfPriorityMatch float64              `json:"cost_of_priority_match"` // copied from policy state
}

func (x *ConcernState) Copy() *ConcernState {
	z := *x
	z.EligibleProposals = slices.Clone(x.EligibleProposals)
	return &z
}

func (x *ConcernState) ProjectedBounty() float64 {
	if x.PriorityScore < 0 {
		return 0
	}
	return x.PriorityScore * x.CostOfPriorityMatch
}

func NewConcernState(id motionproto.MotionID, costOfPriorityMatch float64) *ConcernState {
	return &ConcernState{
		PriorityPoll:        ConcernPollBallotName(id),
		CostOfPriorityMatch: costOfPriorityMatch,
	}
}

type ConcernPolicyState struct {
	CostOfPriorityMatch   float64 `json:"cost_of_priority_match"`
	CostOfReviewForAuthor float64 `json:"cost_of_review_for_author"`
	TotalCostOfPriority   float64 `json:"total_cost_of_priority"`
}

var InitialPolicyState = &ConcernPolicyState{
	CostOfPriorityMatch:   1.0,
	CostOfReviewForAuthor: 0.2,
}
