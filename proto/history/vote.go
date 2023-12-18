package history

type VoteEvent struct {
	By       User     `json:"by"`
	Receipts Receipts `json:"receipts"`
}
