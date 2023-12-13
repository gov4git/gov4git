package proposal

import (
	"github.com/gov4git/gov4git/proto/account"
	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/docket/schema"
)

type CloseReport struct {
	Accepted            bool            `json:"accepted"`
	ApprovalPollOutcome common.Outcome  `json:"approval_poll_outcome"`
	Resolved            schema.Motions  `json:"resolved"`
	Bounty              account.Holding `json:"bounty"`
	BountyDonated       bool            `json:"bounty_donated"`
	Rewarded            Rewards         `json:"rewards"`
}

type CancelReport struct {
	ApprovalPollOutcome common.Outcome `json:"approval_poll_outcome"`
}
