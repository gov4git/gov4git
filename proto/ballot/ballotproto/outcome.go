package ballotproto

import (
	"github.com/gov4git/gov4git/v2/proto/account"
	"github.com/gov4git/gov4git/v2/proto/history/metric"
	"github.com/gov4git/gov4git/v2/proto/member"
)

type Summary string

type Outcome struct {
	Summary      string                                      `json:"summary"`
	Scores       map[string]float64                          `json:"scores"`
	ScoresByUser map[member.User]map[string]StrengthAndScore `json:"scores_by_user"`
	Refunded     map[member.User]account.Holding             `json:"refunded"`
}

func (o Outcome) RefundedHistoryReceipts() metric.Receipts {
	r := metric.Receipts{}
	for user, h := range o.Refunded {
		r = append(r,
			metric.Receipt{
				To:     user.MetricAccountID(),
				Type:   metric.ReceiptTypeRefund,
				Amount: h.MetricHolding(),
			},
		)
	}
	return r
}
