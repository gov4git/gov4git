package proposal

import (
	"github.com/gov4git/gov4git/v2/proto/account"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
)

type CloseReport struct {
	Accepted            bool                `json:"accepted"`
	ApprovalPollOutcome ballotproto.Outcome `json:"approval_poll_outcome"`
	Resolved            motionproto.Motions `json:"resolved"`
	Bounty              account.Holding     `json:"bounty"`
	BountyDonated       bool                `json:"bounty_donated"`
	Rewarded            Rewards             `json:"rewards"`
}

type CancelReport struct {
	ApprovalPollOutcome ballotproto.Outcome `json:"approval_poll_outcome"`
}
