package metrics

import (
	"time"

	"github.com/gov4git/gov4git/v2/proto/history"
	"github.com/gov4git/gov4git/v2/proto/journal"
)

type Series struct {
	DailyNumJoins DailySeries
	//
	DailyNumMotionOpen   DailySeries
	DailyNumMotionClose  DailySeries
	DailyNumMotionCancel DailySeries
	//
	DailyNumVotes DailySeries
	//
	DailyCreditsIssued      DailySeries
	DailyCreditsBurned      DailySeries
	DailyCreditsTransferred DailySeries
	//
	DailyClearedBounties DailySeries
	DailyClearedRewards  DailySeries
	DailyClearedRefunds  DailySeries
	//
	DailyVoteCharges DailySeries
}

type DailySeries struct {
	X []time.Time
	Y []float64
}

func (ds DailySeries) Len() int {
	return len(ds.X)
}

func (ds DailySeries) Total() float64 {
	t := 0.0
	for _, v := range ds.Y {
		t += v
	}
	return t
}

func ComputeSeries(
	entries journal.Entries[*history.Event],
	earliest time.Time,
	latest time.Time,

) *Series {

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
	dailyVoteCharges := DailyBuckets{}

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
					dailyVoteCharges.Add(e.Stamp, r.Amount.Quantity)
				case history.ReceiptTypeRefund:
				case history.ReceiptTypeReward:
				}
			}
		}
	}

	// all daily series have the same x axis entries
	s := &Series{
		DailyNumJoins:           dailyNumJoins.XY(earliest, latest),
		DailyNumMotionOpen:      dailyNumMotionOpen.XY(earliest, latest),
		DailyNumMotionClose:     dailyNumMotionClose.XY(earliest, latest),
		DailyNumMotionCancel:    dailyNumMotionCancel.XY(earliest, latest),
		DailyNumVotes:           dailyNumVotes.XY(earliest, latest),
		DailyCreditsIssued:      dailyCreditIssued.XY(earliest, latest),
		DailyCreditsBurned:      dailyCreditBurned.XY(earliest, latest),
		DailyCreditsTransferred: dailyCreditTransferred.XY(earliest, latest),
		DailyClearedBounties:    dailyCreditInBounties.XY(earliest, latest),
		DailyClearedRewards:     dailyCreditInRewards.XY(earliest, latest),
		DailyClearedRefunds:     dailyCreditInRefunds.XY(earliest, latest),
		DailyVoteCharges:        dailyVoteCharges.XY(earliest, latest),
	}

	return s
}
