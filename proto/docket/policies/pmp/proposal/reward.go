package proposal

import (
	"sort"

	"github.com/gov4git/gov4git/v2/proto/account"
	"github.com/gov4git/gov4git/v2/proto/history"
	"github.com/gov4git/gov4git/v2/proto/member"
)

type Reward struct {
	To     member.User     `json:"to"`
	Amount account.Holding `json:"amount"`
}

func (x Reward) HistoryReceipt() history.Receipt {
	return history.Receipt{
		To:     history.AccountID(member.UserAccountID(x.To)),
		Type:   history.ReceiptTypeReward,
		Amount: x.Amount.HistoryHolding(),
	}
}

type Rewards []Reward

func (x Rewards) Len() int {
	return len(x)
}

func (x Rewards) Less(i, j int) bool {
	return x[i].To < x[j].To
}

func (x Rewards) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

func (x Rewards) Sort() {
	sort.Sort(x)
}

func (x Rewards) HistoryReceipts() history.Receipts {
	r := make(history.Receipts, len(x))
	for i := range x {
		r[i] = x[i].HistoryReceipt()
	}
	return r
}
