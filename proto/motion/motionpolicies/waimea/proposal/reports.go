package proposal

import (
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
)

type CloseReport struct {
	Accepted            bool                `json:"accepted"`
	AgainstPopular      bool                `json:"against_popular"`
	ApprovalPollOutcome ballotproto.Outcome `json:"approval_poll_outcome"`
	ApprovalScore       float64             `json:"approval_score"`
	Resolved            motionproto.Motions `json:"resolved"`
	// reviewers
	CostOfReview   float64 `json:"cost_of_review"`
	Rewarded       Rewards `json:"rewards"`
	RewardDonation float64 `json:"reward_donation"`
	// author
	CostOfPriority          float64 `json:"cost_of_priority"`
	ProjectedPriorityBounty float64 `json:"projected_priority_bounty"`
	ProjectedReviewBounty   float64 `json:"projected_review_bounty"`
	RealizedBounty          float64 `json:"realized_bounty"`
	BountyDonation          float64 `json:"bounty_donation"`
}

type CancelReport struct {
	ApprovalPollOutcome ballotproto.Outcome `json:"approval_poll_outcome"`
}
