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

		receipts := metric.Receipts{}

		// accepting a proposal against the popular vote?
		againstPopular := adt.Tally.Scores[pmp_1.ProposalBallotChoice] < 0

		// close the approval poll, move the funds to a reward account
		closeApprovalPoll := ballotapi.Close_StageOnly(
			ctx,
			cloned,
			approvalPollName,
			pmp_1.ProposalRewardAccountID(prop.ID),
		)

		// reward reviewers

		rewards, rewardDonationReceipts, rewardDonation := disberseRewards(ctx, cloned, prop, true)
		receipts = append(receipts, rewards.MetricReceipts()...)
		receipts = append(receipts, rewardDonationReceipts...)

		// reward author

		// close all concerns consResolved by the motion, and
		// transfer their funds into the bounty account
		consResolved, consEscrows := loadResolvedConcerns(ctx, cloned, prop)
		consFunds := closeResolvedConcerns(ctx, cloned, prop, consResolved)

		// total escrow is the target pay for the author
		escrowTotal := 0.0
		for _, conEscrow := range consEscrows {
			escrowTotal += conEscrow
		}

		bountyAccount := pmp_1.ProposalBountyAccountID(prop.ID)

		award := 0.0 // award to author
		bountyDonation := 0.0

		if prop.Author.IsNone() {

			account.Transfer_StageOnly(
				ctx,
				cloned.PublicClone(),
				bountyAccount,
				pmp_0.MatchingPoolAccountID,
				consFunds,
				fmt.Sprintf("bounty for proposal %v was donated", prop.ID),
			)
			bountyDonation = consFunds.Quantity
			receipts = append(receipts,
				metric.Receipt{
					To:     pmp_0.MatchingPoolAccountID.MetricAccountID(),
					Type:   metric.ReceiptTypeBounty,
					Amount: consFunds.MetricHolding(),
				},
			)

		} else {

			authorAccount := member.UserAccountID(prop.Author)

			matchAccount := account.Get_Local(ctx, cloned.PublicClone(), pmp_0.MatchingPoolAccountID)
			matchFunds := matchAccount.Balance(account.PluralAsset).Quantity

			awardFromCon, awardFromMatch, donateFromCon := calcAuthorAward(consFunds.Quantity, matchFunds, escrowTotal)
			award = awardFromCon + awardFromMatch
			bountyDonation = donateFromCon

			if awardFromCon > 0 {
				to := authorAccount
				amt := account.H(account.PluralAsset, awardFromCon)
				account.Transfer_StageOnly(
					ctx,
					cloned.PublicClone(),
					bountyAccount,
					to,
					amt,
					fmt.Sprintf("bounty to proposal %v author", prop.ID),
				)
				receipts = append(receipts,
					metric.Receipt{
						To:     to.MetricAccountID(),
						Type:   metric.ReceiptTypeBounty,
						Amount: amt.MetricHolding(),
					},
				)
			}

			if awardFromMatch > 0 {
				to := authorAccount
				amt := account.H(account.PluralAsset, awardFromMatch)
				account.Transfer_StageOnly(
					ctx,
					cloned.PublicClone(),
					pmp_0.MatchingPoolAccountID,
					to,
					amt,
					fmt.Sprintf("matched bounty to proposal %v author", prop.ID),
				)
				receipts = append(receipts,
					metric.Receipt{
						To:     to.MetricAccountID(),
						Type:   metric.ReceiptTypeBounty,
						Amount: amt.MetricHolding(),
					},
				)
			}

			if donateFromCon > 0 {
				to := pmp_0.MatchingPoolAccountID
				amt := account.H(account.PluralAsset, donateFromCon)
				account.Transfer_StageOnly(
					ctx,
					cloned.PublicClone(),
					bountyAccount,
					to,
					amt,
					fmt.Sprintf("donation of unused bounty for proposal %v to matching pool", prop.ID),
				)
				receipts = append(receipts,
					metric.Receipt{
						To:     to.MetricAccountID(),
						Type:   metric.ReceiptTypeDonation,
						Amount: amt.MetricHolding(),
					},
				)
			}

		}

		// metrics
		metric.Log_StageOnly(ctx, cloned.PublicClone(), &metric.Event{
			Motion: &metric.MotionEvent{
				Close: &metric.MotionClose{
					ID:       metric.MotionID(prop.ID),
					Type:     "proposal-v1",
					Policy:   metric.MotionPolicy(prop.Policy),
					Decision: decision.MetricDecision(),
					Receipts: receipts,
				},
			},
		})

		return &CloseReport{
				Accepted:            true,
				ApprovalPollOutcome: closeApprovalPoll.Result,
				Resolved:            consResolved,
				Rewarded:            rewards,
				RewardDonation:      rewardDonation,
				Bounty:              consFunds.Quantity,
				Escrow:              escrowTotal,
				Award:               award,
				BountyDonation:      bountyDonation,
			}, closeNotice(
				ctx,
				prop,
				isMerged,
				againstPopular,
				closeApprovalPoll.Result,
				consResolved,
				consFunds,
				bountyDonation,
				rewards,
				rewardDonation,
			)

	} else {

		// rejecting a proposal against the popular vote?
		againstPopular := adt.Tally.Scores[pmp_1.ProposalBallotChoice] > 0

		// close the approval poll, move the funds to a reward account
		closeApprovalPoll := ballotapi.Close_StageOnly(
			ctx,
			cloned,
			approvalPollName,
			pmp_1.ProposalRewardAccountID(prop.ID),
		)

		// reward reviewers
		rewards, donationReceipt, rewardDonation := disberseRewards(ctx, cloned, prop, false)

		// metrics
		metric.Log_StageOnly(ctx, cloned.PublicClone(), &metric.Event{
			Motion: &metric.MotionEvent{
				Close: &metric.MotionClose{
					ID:       metric.MotionID(prop.ID),
					Type:     "proposal-v1",
					Policy:   metric.MotionPolicy(prop.Policy),
					Decision: decision.MetricDecision(),
					Receipts: append(rewards.MetricReceipts(), donationReceipt...),
				},
			},
		})

		return &CloseReport{
				Accepted:            false,
				ApprovalPollOutcome: closeApprovalPoll.Result,
				Resolved:            nil,
				Rewarded:            rewards,
				RewardDonation:      rewardDonation,
				Bounty:              0.0,
				Escrow:              0.0,
				Award:               0.0,
				BountyDonation:      0.0,
			}, closeNotice(
				ctx,
				prop,
				isMerged,
				againstPopular,
				closeApprovalPoll.Result,
				nil,
				account.H(account.PluralAsset, 0.0),
				0.0,
				rewards,
				rewardDonation,
			)

	}
}
