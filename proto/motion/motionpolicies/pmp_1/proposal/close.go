package proposal

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/v2/proto/account"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotapi"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/history/metric"
	"github.com/gov4git/gov4git/v2/proto/member"
	"github.com/gov4git/gov4git/v2/proto/motion/motionpolicies/pmp_0"
	"github.com/gov4git/gov4git/v2/proto/motion/motionpolicies/pmp_1"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
	"github.com/gov4git/gov4git/v2/proto/notice"
)

func (x proposalPolicy) Close(
	ctx context.Context,
	cloned gov.OwnerCloned,
	prop motionproto.Motion,
	decision motionproto.Decision,
	_ ...any,

) (motionproto.Report, notice.Notices) {

	// ensure that eligible set in policy state is valid
	x.Update(ctx, cloned, prop)

	// was the PR merged or not
	isMerged := decision.IsAccept()

	approvalPollName := pmp_1.ProposalApprovalPollName(prop.ID)
	adt := loadPropApprovalPollTally(ctx, cloned.PublicClone(), prop)

	if isMerged {

		// accepting a proposal against the popular vote?
		againstPopular := adt.Tally.Scores[pmp_1.ProposalBallotChoice] < 0

		// close the approval poll, move the funds to a reward account
		approvalPollName := pmp_1.ProposalApprovalPollName(prop.ID)
		closeApprovalPoll := ballotapi.Close_StageOnly(
			ctx,
			cloned,
			approvalPollName,
			pmp_1.ProposalRewardAccountID(prop.ID),
		)

		// close all concerns consResolved by the motion, and
		// transfer their funds into the bounty account
		consResolved, consEscrows := loadResolvedConcerns(ctx, cloned, prop)
		consFunds := closeResolvedConcerns(ctx, cloned, prop, consResolved)

		// XXX: reward reviewers
		_ = consEscrows

		rewards := disberseRewards(ctx, cloned, prop)

		// XXX: reward author

		var bountyDonated bool
		bountyReceipt := metric.Receipt{
			Type:   metric.ReceiptTypeBounty,
			Amount: consFunds.MetricHolding(),
		}
		if prop.Author.IsNone() {
			account.Transfer_StageOnly(
				ctx,
				cloned.PublicClone(),
				pmp_1.ProposalBountyAccountID(prop.ID),
				pmp_0.MatchingPoolAccountID,
				consFunds,
				fmt.Sprintf("bounty for proposal %v", prop.ID),
			)
			bountyDonated = true
			bountyReceipt.To = pmp_0.MatchingPoolAccountID.HistoryAccountID()
		} else {
			account.Transfer_StageOnly(
				ctx,
				cloned.PublicClone(),
				pmp_1.ProposalBountyAccountID(prop.ID),
				member.UserAccountID(prop.Author),
				consFunds,
				fmt.Sprintf("bounty for proposal %v", prop.ID),
			)
			bountyReceipt.To = member.UserAccountID(prop.Author).HistoryAccountID()
		}

		// metrics
		metric.Log_StageOnly(ctx, cloned.PublicClone(), &metric.Event{
			Motion: &metric.MotionEvent{
				Close: &metric.MotionClose{
					ID:       metric.MotionID(prop.ID),
					Type:     "proposal-v1",
					Policy:   metric.MotionPolicy(prop.Policy),
					Decision: decision.MetricDecision(),
					Receipts: append(rewards.MetricReceipts(), bountyReceipt),
				},
			},
		})

		return &CloseReport{
			Accepted:            true,
			ApprovalPollOutcome: closeApprovalPoll.Result,
			Resolved:            consResolved,
			Bounty:              consFunds,
			BountyDonated:       bountyDonated,
			Rewarded:            rewards,
		}, closeNotice(ctx, prop, againstPopular, closeApprovalPoll.Result, consResolved, consFunds, bountyDonated, rewards)

	} else {

		// rejecting a proposal against the popular vote?
		againstPopular := adt.Tally.Scores[pmp_1.ProposalBallotChoice] > 0

		// cancel the referendum for the motion (refunds voters)
		cancelApprovalPoll := ballotapi.Cancel_StageOnly(
			ctx,
			cloned,
			approvalPollName,
		)

		// metrics
		metric.Log_StageOnly(ctx, cloned.PublicClone(), &metric.Event{
			Motion: &metric.MotionEvent{
				Close: &metric.MotionClose{
					ID:       metric.MotionID(prop.ID),
					Type:     "proposal-v1",
					Policy:   metric.MotionPolicy(prop.Policy),
					Decision: decision.MetricDecision(),
					Receipts: cancelApprovalPoll.Result.RefundedHistoryReceipts(),
				},
			},
		})

		return &CloseReport{
			Accepted:            false,
			ApprovalPollOutcome: cancelApprovalPoll.Result,
			Resolved:            nil,
			Bounty:              account.H(account.PluralAsset, 0.0),
			BountyDonated:       false,
			Rewarded:            nil,
		}, cancelNotice(ctx, prop, againstPopular, cancelApprovalPoll.Result)

	}
}
