package history

type VoteEvent struct {
	By       User        `json:"by"`
	Context  VoteContext `json:"context"`
	Receipts Receipts    `json:"receipts"`
}

type VoteContext string

const (
	VoteContextConcern  VoteContext = "concern"
	VoteContextProposal VoteContext = "proposal"
)
