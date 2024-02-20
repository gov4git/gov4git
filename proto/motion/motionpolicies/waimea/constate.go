package waimea

import (
	"slices"

	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
)

const StateFilebase = "state.json"

type ConcernState struct {
	PriorityPoll      ballotproto.BallotID `json:"priority_poll"`
	CostOfPriority    float64              `json:"cost_of_priority"`
	PriorityScore     float64              `json:"priority_score"`
	PriorityMatch     float64              `json:"priority_match"` // copied from policy state
	EligibleProposals motionproto.Refs     `json:"eligible_proposals"`
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
	return x.PriorityScore * x.PriorityMatch
}

func NewConcernState(id motionproto.MotionID, priorityMatch float64) *ConcernState {
	return &ConcernState{
		PriorityPoll:  ConcernPollBallotName(id),
		PriorityMatch: priorityMatch,
	}
}

type ConcernPolicyState struct {
	// parameters
	PriorityMatch float64 `json:"priority_match"`
	ReviewMatch   float64 `json:"review_match"`
	// state
	TotalCostOfPriority float64 `json:"total_cost_of_priority"`
	TotalCostOfReview   float64 `json:"total_cost_of_review"`
}

var InitialPolicyState = &ConcernPolicyState{
	PriorityMatch:       2.0,
	ReviewMatch:         2.0,
	TotalCostOfPriority: 0.0,
	TotalCostOfReview:   0.0,
}
