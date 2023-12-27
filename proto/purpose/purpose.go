package purpose

import "github.com/gov4git/gov4git/v2/proto/history/metric"

type Purpose string

const (
	Unspecified Purpose = "unspecified"
	Concern     Purpose = "concern"
	Proposal    Purpose = "proposal"
)

func (p Purpose) MetricVotePurpose() metric.VotePurpose {
	switch p {
	case Concern:
		return metric.VotePurposeConcern
	case Proposal:
		return metric.VotePurposeProposal
	}
	return metric.VotePurposeUnspecified
}
