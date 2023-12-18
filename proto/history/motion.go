package history

type MotionEvent struct {
	Open   *MotionOpen   `json:"open"`
	Close  *MotionClose  `json:"close"`
	Cancel *MotionCancel `json:"cancel"`
}

type MotionID string

type MotionOpen struct {
	ID   MotionID `json:"id"`
	Type string   `json:"type"`
}

type MotionClose struct {
	ID       MotionID `json:"id"`
	Type     string   `json:"type"`
	Receipts Receipts `json:"receipts"`
}

type MotionCancel struct {
	ID       MotionID `json:"id"`
	Type     string   `json:"type"`
	Receipts Receipts `json:"receipts"`
}

type Receipt struct {
	To     AccountID   `json:"account"`
	Type   ReceiptType `json:"type"`
	Amount Holding     `json:"amount"`
}

type Receipts []Receipt

type ReceiptType string

const (
	ReceiptTypeRefund ReceiptType = "refund"
	ReceiptTypeReward ReceiptType = "reward"
	ReceiptTypeBounty ReceiptType = "bounty"
)
