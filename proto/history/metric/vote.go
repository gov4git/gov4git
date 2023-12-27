package metric

type VoteEvent struct {
	By       User        `json:"by"`
	Purpose  VotePurpose `json:"context"`
	Receipts Receipts    `json:"receipts"`
}

type VotePurpose string

const (
	VotePurposeUnspecified VotePurpose = "unspecified"
	VotePurposeConcern     VotePurpose = "concern"
	VotePurposeProposal    VotePurpose = "proposal"
)
