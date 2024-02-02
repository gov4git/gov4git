package proposal

import (
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
)

type CloseReport struct {
	Accepted            bool                `json:"accepted"`
	ApprovalPollOutcome ballotproto.Outcome `json:"approval_poll_outcome"`
	Resolved            motionproto.Motions `json:"resolved"`
	// reviewers
	Rewarded       Rewards `json:"rewards"`
	RewardDonation float64 `json:"reward_donation"`
	// author
	Bounty         float64 `json:"bounty"` // actual funds held in concerns
	Escrow         float64 `json:"escrow"` // target payment (may be different from bounty)
	Award          float64 `json:"award"`  // actual payment to author
	BountyDonation float64 `json:"bounty_donation"`
}

type CancelReport struct {
	ApprovalPollOutcome ballotproto.Outcome `json:"approval_poll_outcome"`
}
