package metric

type AccountID string

type Asset string

type Holding struct {
	Asset    Asset   `json:"asset"`
	Quantity float64 `json:"quantity"`
}

type AccountEvent struct {
	Issue    *AccountIssueEvent    `json:"issue"`
	Burn     *AccountBurnEvent     `json:"burn"`
	Transfer *AccountTransferEvent `json:"transfer"`
}

type AccountIssueEvent struct {
	To     AccountID `json:"to"`
	Amount Holding   `json:"amount"`
}

type AccountBurnEvent struct {
	From   AccountID `json:"from"`
	Amount Holding   `json:"amount"`
}

type AccountTransferEvent struct {
	From   AccountID `json:"from"`
	To     AccountID `json:"to"`
	Amount Holding   `json:"amount"`
}
