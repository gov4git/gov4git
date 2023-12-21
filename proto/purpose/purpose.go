package purpose

import (
	"github.com/gov4git/gov4git/v2/proto/history"
)

type Purpose string

const (
	Unspecified Purpose = "unspecified"
	Concern     Purpose = "concern"
	Proposal    Purpose = "proposal"
)

func (p Purpose) HistoryVotePurpose() history.VotePurpose {
	switch p {
	case Concern:
		return history.VotePurposeConcern
	case Proposal:
		return history.VotePurposeProposal
	}
	return history.VotePurposeUnspecified
}
