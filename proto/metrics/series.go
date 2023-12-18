package metrics

import (
	"time"

	"github.com/gov4git/gov4git/v2/proto/history"
	"github.com/gov4git/gov4git/v2/proto/journal"
)

type Series struct {
	DailyNumJoins          DailySeries
	DailyNumMotionOpen     DailySeries
	DailyNumMotionClose    DailySeries
	DailyNumMotionCancel   DailySeries
	DailyNumVotes          DailySeries
	DailyCreditIssued      DailySeries
	DailyCreditBurned      DailySeries
	DailyCreditTransferred DailySeries
	DailyCreditInBounties  DailySeries
	DailyCreditInRewards   DailySeries
	DailyCreditInRefunds   DailySeries
	DailyCreditInVotes     DailySeries
}

type DailySeries struct {
	X     []time.Time
	Y     []float64
	Total float64
}

func (ds DailySeries) Len() int {
	return len(ds.X)
}

func ComputeSeries(entries journal.Entries[*history.Event]) *Series {

	// counts
	dailyNumJoins := DailyBuckets{}
	dailyNumMotionOpen := DailyBuckets{}
	dailyNumMotionClose := DailyBuckets{}
	dailyNumMotionCancel := DailyBuckets{}
	dailyNumVotes := DailyBuckets{}
	dailyCreditIssued := DailyBuckets{}
	dailyCreditBurned := DailyBuckets{}
	dailyCreditTransferred := DailyBuckets{}
	dailyCreditInBounties := DailyBuckets{}
	dailyCreditInRewards := DailyBuckets{}
	dailyCreditInRefunds := DailyBuckets{}
	dailyCreditInVotes := DailyBuckets{}

	for _, e := range entries {
		if e.Payload.Account != nil {
			if e.Payload.Account.Burn != nil {
				dailyCreditBurned.Add(e.Stamp, e.Payload.Account.Burn.Amount.Quantity)
			}
			if e.Payload.Account.Issue != nil {
				dailyCreditIssued.Add(e.Stamp, e.Payload.Account.Issue.Amount.Quantity)
			}
			if e.Payload.Account.Transfer != nil {
				dailyCreditTransferred.Add(e.Stamp, e.Payload.Account.Transfer.Amount.Quantity)
			}
		}
		if e.Payload.Join != nil {
			dailyNumJoins.Add(e.Stamp, 1)
		}
		if e.Payload.Motion != nil {
			if e.Payload.Motion.Open != nil {
				dailyNumMotionOpen.Add(e.Stamp, 1)
			}
			if e.Payload.Motion.Close != nil {
				dailyNumMotionClose.Add(e.Stamp, 1)
				for _, r := range e.Payload.Vote.Receipts {
					switch r.Type {
					case history.ReceiptTypeBounty:
						dailyCreditInBounties.Add(e.Stamp, r.Amount.Quantity)
					case history.ReceiptTypeCharge:
					case history.ReceiptTypeRefund:
						dailyCreditInRefunds.Add(e.Stamp, r.Amount.Quantity)
					case history.ReceiptTypeReward:
						dailyCreditInRewards.Add(e.Stamp, r.Amount.Quantity)
					}
				}
			}
			if e.Payload.Motion.Cancel != nil {
				dailyNumMotionCancel.Add(e.Stamp, 1)
			}
		}
		if e.Payload.Vote != nil {
			dailyNumVotes.Add(e.Stamp, 1)
			for _, r := range e.Payload.Vote.Receipts {
				switch r.Type {
				case history.ReceiptTypeBounty:
				case history.ReceiptTypeCharge:
					dailyCreditInVotes.Add(e.Stamp, r.Amount.Quantity)
				case history.ReceiptTypeRefund:
				case history.ReceiptTypeReward:
				}
			}
		}
	}

	s := &Series{
		DailyNumJoins:          dailyNumJoins.XY(),
		DailyNumMotionOpen:     dailyNumMotionOpen.XY(),
		DailyNumMotionClose:    dailyNumMotionClose.XY(),
		DailyNumMotionCancel:   dailyNumMotionCancel.XY(),
		DailyNumVotes:          dailyNumVotes.XY(),
		DailyCreditIssued:      dailyCreditIssued.XY(),
		DailyCreditBurned:      dailyCreditBurned.XY(),
		DailyCreditTransferred: dailyCreditTransferred.XY(),
		DailyCreditInBounties:  dailyCreditInBounties.XY(),
		DailyCreditInRewards:   dailyCreditInRewards.XY(),
		DailyCreditInRefunds:   dailyCreditInRefunds.XY(),
		DailyCreditInVotes:     dailyCreditInVotes.XY(),
	}

	return s
}
