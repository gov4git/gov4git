package metrics

import (
	"time"

	"github.com/gov4git/gov4git/v2/proto/history/metric"
	"github.com/gov4git/gov4git/v2/proto/journal"
)

type Series struct {
	DailyNumJoins DailySeries
	//
	DailyNumMotionOpen   DailySeries
	DailyNumMotionClose  DailySeries
	DailyNumMotionCancel DailySeries
	//
	DailyCreditsIssued      DailySeries
	DailyCreditsBurned      DailySeries
	DailyCreditsTransferred DailySeries
	//
	DailyClearedBounties DailySeries
	DailyClearedRewards  DailySeries
	DailyClearedRefunds  DailySeries
	//
	DailyNumConcernVotes  DailySeries
	DailyNumProposalVotes DailySeries
	DailyNumOtherVotes    DailySeries
	//
	DailyConcernVoteCharges  DailySeries
	DailyProposalVoteCharges DailySeries
	DailyOtherVoteCharges    DailySeries
}

func ComputeSeries(
	entries journal.Entries[*metric.Event],
	earliest time.Time,
	latest time.Time,

) *Series {

	dailyNumJoins := DailyBuckets{}
	dailyNumMotionOpen := DailyBuckets{}
	dailyNumMotionClose := DailyBuckets{}
	dailyNumMotionCancel := DailyBuckets{}
	dailyCreditIssued := DailyBuckets{}
	dailyCreditBurned := DailyBuckets{}
	dailyCreditTransferred := DailyBuckets{}
	dailyCreditInBounties := DailyBuckets{}
	dailyCreditInRewards := DailyBuckets{}
	dailyCreditInRefunds := DailyBuckets{}
	dailyNumConcernVotes := DailyBuckets{}
	dailyNumProposalVotes := DailyBuckets{}
	dailyNumOtherVotes := DailyBuckets{}
	dailyConcernVoteCharges := DailyBuckets{}
	dailyProposalVoteCharges := DailyBuckets{}
	dailyOtherVoteCharges := DailyBuckets{}

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
					case metric.ReceiptTypeBounty:
						dailyCreditInBounties.Add(e.Stamp, r.Amount.Quantity)
					case metric.ReceiptTypeCharge:
					case metric.ReceiptTypeRefund:
						dailyCreditInRefunds.Add(e.Stamp, r.Amount.Quantity)
					case metric.ReceiptTypeReward:
						dailyCreditInRewards.Add(e.Stamp, r.Amount.Quantity)
					}
				}
			}
			if e.Payload.Motion.Cancel != nil {
				dailyNumMotionCancel.Add(e.Stamp, 1)
			}
		}
		if e.Payload.Vote != nil {
			switch e.Payload.Vote.Purpose {
			case metric.VotePurposeConcern:
				dailyNumConcernVotes.Add(e.Stamp, 1)
			case metric.VotePurposeProposal:
				dailyNumProposalVotes.Add(e.Stamp, 1)
			default:
				dailyNumOtherVotes.Add(e.Stamp, 1)
			}
			for _, r := range e.Payload.Vote.Receipts {
				switch r.Type {
				case metric.ReceiptTypeBounty:
				case metric.ReceiptTypeCharge:
					switch e.Payload.Vote.Purpose {
					case metric.VotePurposeConcern:
						dailyConcernVoteCharges.Add(e.Stamp, r.Amount.Quantity)
					case metric.VotePurposeProposal:
						dailyProposalVoteCharges.Add(e.Stamp, r.Amount.Quantity)
					default:
						dailyOtherVoteCharges.Add(e.Stamp, r.Amount.Quantity)
					}
				case metric.ReceiptTypeRefund:
				case metric.ReceiptTypeReward:
				}
			}
		}
	}

	// all daily series have the same x axis entries
	s := &Series{
		DailyNumJoins:            dailyNumJoins.XY(earliest, latest),
		DailyNumMotionOpen:       dailyNumMotionOpen.XY(earliest, latest),
		DailyNumMotionClose:      dailyNumMotionClose.XY(earliest, latest),
		DailyNumMotionCancel:     dailyNumMotionCancel.XY(earliest, latest),
		DailyCreditsIssued:       dailyCreditIssued.XY(earliest, latest),
		DailyCreditsBurned:       dailyCreditBurned.XY(earliest, latest),
		DailyCreditsTransferred:  dailyCreditTransferred.XY(earliest, latest),
		DailyClearedBounties:     dailyCreditInBounties.XY(earliest, latest),
		DailyClearedRewards:      dailyCreditInRewards.XY(earliest, latest),
		DailyClearedRefunds:      dailyCreditInRefunds.XY(earliest, latest),
		DailyNumConcernVotes:     dailyNumConcernVotes.XY(earliest, latest),
		DailyNumProposalVotes:    dailyNumProposalVotes.XY(earliest, latest),
		DailyNumOtherVotes:       dailyNumOtherVotes.XY(earliest, latest),
		DailyConcernVoteCharges:  dailyConcernVoteCharges.XY(earliest, latest),
		DailyProposalVoteCharges: dailyProposalVoteCharges.XY(earliest, latest),
		DailyOtherVoteCharges:    dailyOtherVoteCharges.XY(earliest, latest),
	}

	return s
}
