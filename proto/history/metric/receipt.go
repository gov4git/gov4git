package metric

type Receipt struct {
	To     AccountID   `json:"account"`
	Type   ReceiptType `json:"type"`
	Amount Holding     `json:"amount"`
}

type Receipts []Receipt

func OneReceipt(to AccountID, typ ReceiptType, amt Holding) Receipts {
	return Receipts{
		Receipt{To: to, Type: typ, Amount: amt},
	}
}

type ReceiptType string

const (
	ReceiptTypeRefund   ReceiptType = "refund"
	ReceiptTypeReward   ReceiptType = "reward"
	ReceiptTypeBounty   ReceiptType = "bounty"
	ReceiptTypeCharge   ReceiptType = "charge"
	ReceiptTypeDonation ReceiptType = "donation"
)
