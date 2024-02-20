package proposal

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/v2/proto/account"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotapi"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/history/metric"
	"github.com/gov4git/gov4git/v2/proto/member"
	"github.com/gov4git/gov4git/v2/proto/motion/motionapi"
	"github.com/gov4git/gov4git/v2/proto/motion/motionpolicies/waimea"
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

	approvalPollName := waimea.ProposalApprovalPollName(prop.ID)

	// ensure that eligible set in policy state is valid
	x.Update(ctx, cloned, prop)
	ballotapi.ReTally_StageOnly(ctx, cloned.PublicClone(), approvalPollName)
	x.Update(ctx, cloned, prop)

	propState := motionapi.LoadPolicyState_Local[*waimea.ProposalState](ctx, cloned.PublicClone(), prop.ID)

	// was the PR merged or not
	isMerged := decision.IsAccept()

	approvalTally := loadApprovalPoll(ctx, cloned.PublicClone(), prop)
	costOfReview := approvalTally.Tally.Capitalization()

	if isMerged {

		receipts := metric.Receipts{}

		// accepting a proposal against the popular vote?
		againstPopular := propState.ApprovalScore < 0

		// close the approval poll, move the funds to a reward account
		closeApprovalPoll := ballotapi.Close_StageOnly(
			ctx,
			cloned,
			approvalPollName,
			waimea.ProposalRewardAccountID(prop.ID),
		)

		// reward reviewers

		rewards, rewardDonationReceipts, rewardDonation := calcReviewersRewards(ctx, cloned, prop, true)
		receipts = append(receipts, rewards.MetricReceipts()...)
		receipts = append(receipts, rewardDonationReceipts...)

		// reward author

		// close all concerns by the motion, and
		// transfer their funds into the bounty account
		resolvedCons, _, projectedPriorityBounty := loadResolvedConcerns(ctx, cloned, prop)
		priorityFunds := max(0, closeResolvedConcerns(ctx, cloned, prop, resolvedCons))

		bountyAccountID := waimea.ProposalBountyAccountID(prop.ID)

		realizedBounty := 0.0 // award to author
		bountyDonation := 0.0 // leftovers to the penny jar
		projectedReviewBounty := propState.ProjectedApprovalBounty()

		if prop.Author.IsNone() {

			bountyDonation = priorityFunds
			if priorityFunds > 0 {
				donation := account.H(account.PluralAsset, priorityFunds)
				account.Transfer_StageOnly(
					ctx,
					cloned.PublicClone(),
					bountyAccountID,
					waimea.PennyAccountID,
					donation,
					fmt.Sprintf("bounty for proposal %v was donated to the penny jar", prop.ID),
				)
				receipts = append(receipts,
					metric.Receipt{
						To:     waimea.PennyAccountID.MetricAccountID(),
						Type:   metric.ReceiptTypeBounty,
						Amount: donation.MetricHolding(),
					},
				)
			}

		} else {

			authorAccount := member.UserAccountID(prop.Author)

			realizedBounty = max(0, projectedPriorityBounty+projectedReviewBounty)
			bountyDonation = 0

			if realizedBounty > 0 {
				to := authorAccount
				amt := account.H(account.PluralAsset, realizedBounty)
				account.Issue_StageOnly(
					ctx,
					cloned.PublicClone(),
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

		}

		// metrics
		metric.Log_StageOnly(ctx, cloned.PublicClone(), &metric.Event{
			Motion: &metric.MotionEvent{
				Close: &metric.MotionClose{
					ID:       metric.MotionID(prop.ID),
					Type:     "proposal",
					Policy:   metric.MotionPolicy(prop.Policy),
					Decision: decision.MetricDecision(),
					Receipts: receipts,
				},
			},
		})

		report := &CloseReport{
			Accepted:                true,
			AgainstPopular:          againstPopular,
			ApprovalPollOutcome:     closeApprovalPoll.Result,
			ApprovalScore:           propState.ApprovalScore,
			Resolved:                resolvedCons,
			CostOfReview:            costOfReview,
			Rewarded:                rewards,
			RewardDonation:          rewardDonation,
			CostOfPriority:          priorityFunds,
			ProjectedPriorityBounty: projectedPriorityBounty,
			ProjectedReviewBounty:   projectedReviewBounty,
			RealizedBounty:          realizedBounty,
			BountyDonation:          bountyDonation,
		}
		return report, closeNotice(ctx, prop, report)

	} else {

		// rejecting a proposal against the popular vote?
		againstPopular := propState.ApprovalScore > 0

		// close the approval poll, move the funds to a reward account
		closeApprovalPoll := ballotapi.Close_StageOnly(
			ctx,
			cloned,
			approvalPollName,
			waimea.ProposalRewardAccountID(prop.ID),
		)

		// reward reviewers
		rewards, donationReceipt, rewardDonation := calcReviewersRewards(ctx, cloned, prop, false)

		// metrics
		metric.Log_StageOnly(ctx, cloned.PublicClone(), &metric.Event{
			Motion: &metric.MotionEvent{
				Close: &metric.MotionClose{
					ID:       metric.MotionID(prop.ID),
					Type:     "proposal",
					Policy:   metric.MotionPolicy(prop.Policy),
					Decision: decision.MetricDecision(),
					Receipts: append(rewards.MetricReceipts(), donationReceipt...),
				},
			},
		})

		report := &CloseReport{
			Accepted:                false,
			AgainstPopular:          againstPopular,
			ApprovalPollOutcome:     closeApprovalPoll.Result,
			ApprovalScore:           propState.ApprovalScore,
			Resolved:                nil,
			CostOfReview:            costOfReview,
			Rewarded:                rewards,
			RewardDonation:          rewardDonation,
			CostOfPriority:          0.0,
			ProjectedPriorityBounty: 0.0,
			ProjectedReviewBounty:   0.0,
			RealizedBounty:          0.0,
			BountyDonation:          0.0,
		}
		return report, closeNotice(ctx, prop, report)

	}
}
