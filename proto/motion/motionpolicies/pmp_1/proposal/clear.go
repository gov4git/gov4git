package proposal

import (
	"context"
	"fmt"
	"math"

	"github.com/gov4git/gov4git/v2/proto/account"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotapi"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/history/metric"
	"github.com/gov4git/gov4git/v2/proto/member"
	"github.com/gov4git/gov4git/v2/proto/motion/motionapi"
	"github.com/gov4git/gov4git/v2/proto/motion/motionpolicies/pmp_0"
	"github.com/gov4git/gov4git/v2/proto/motion/motionpolicies/pmp_1"
	"github.com/gov4git/gov4git/v2/proto/motion/motionpolicies/pmp_1/concern"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
	"github.com/gov4git/lib4git/base"
)

func loadResolvedConcerns(
	ctx context.Context,
	cloned gov.OwnerCloned,
	prop motionproto.Motion,

) (resolved motionproto.Motions, projectedBounties []float64, projectedBounty float64) {

	eligible := calcEligibleConcerns(ctx, cloned.PublicClone(), prop)
	for _, ref := range eligible {
		con := motionapi.LookupMotion_Local(ctx, cloned.PublicClone(), ref.To)
		conState := motionapi.LoadPolicyState_Local[*concern.ConcernState](ctx, cloned.PublicClone(), con.ID)
		//
		resolved = append(resolved, con)
		projectedBounties = append(projectedBounties, conState.ProjectedBounty())
	}

	projectedBounty = 0.0
	for _, pb := range projectedBounties {
		projectedBounty += pb
	}

	return resolved, projectedBounties, projectedBounty
}

func calcEligibleConcerns(ctx context.Context, cloned gov.Cloned, prop motionproto.Motion) motionproto.Refs {
	eligible := motionproto.Refs{}
	for _, ref := range prop.RefTo {
		if pmp_1.IsConcernProposalEligible(ctx, cloned, ref.To, prop.ID, ref.Type) {
			eligible = append(eligible, ref)
		}
	}
	eligible.Sort()
	return eligible
}

func closeResolvedConcerns(
	ctx context.Context,
	cloned gov.OwnerCloned,
	prop motionproto.Motion,
	cons motionproto.Motions,

) account.Holding {

	for _, con := range cons {
		// close resolved concerns, and transfer concern escrows to proposal-owned bounty account
		motionapi.CloseMotion_StageOnly(
			ctx,
			cloned,
			con.ID,
			motionproto.Accept,
			pmp_1.ProposalBountyAccountID(prop.ID), // account to send bounty to
			prop,                                   // proposal that resolves the issue
		)
	}

	return account.Get_Local(
		ctx,
		cloned.PublicClone(),
		pmp_1.ProposalBountyAccountID(prop.ID),
	).Assets.Balance(account.PluralAsset)
}

func loadPropApprovalPollTally(
	ctx context.Context,
	cloned gov.Cloned,
	prop motionproto.Motion,

) ballotproto.AdTally {

	pollName := pmp_1.ProposalApprovalPollName(prop.ID)
	return ballotapi.Show_Local(ctx, cloned.Tree(), pollName)
}

func calcReviewersRewards(
	ctx context.Context,
	cloned gov.OwnerCloned,
	prop motionproto.Motion,
	accepted bool,

) (Rewards, metric.Receipts, float64) {

	rewards := Rewards{}
	adt := loadPropApprovalPollTally(ctx, cloned.PublicClone(), prop)

	isWinner := func(score float64) bool {
		if accepted {
			return score > 0
		}
		return score < 0
	}

	// compute rewards
	idealPayout := map[member.User]float64{}
	idealFunds := 0.0
	realFunds := 0.0
	for user, choices := range adt.Tally.ScoresByUser {
		ss := choices[pmp_1.ProposalBallotChoice]
		if isWinner(ss.Score) {
			pay := math.Abs(2 * ss.Score)
			idealPayout[user] = pay
			idealFunds += pay
		}
		realFunds += math.Abs(ss.Strength)
	}

	// compute shrink factor
	if idealFunds <= 0.0 {
		return nil, nil, 0.0
	}
	payoutShrinkFactor := min(realFunds/idealFunds, 1.0)
	if payoutShrinkFactor < 1.0 {
		base.Infof("proposal %v has reviewer reward shrink factor %v", prop.ID, payoutShrinkFactor)
	}

	// disberse payouts
	for user, choices := range adt.Tally.ScoresByUser {
		ss := choices[pmp_1.ProposalBallotChoice]
		if isWinner(ss.Score) {
			payout := account.H(
				account.PluralAsset,
				payoutShrinkFactor*idealPayout[user],
			)
			rewards = append(rewards,
				Reward{
					To:     user,
					Amount: payout,
				},
			)
			// transfer reward
			account.Transfer_StageOnly(
				ctx,
				cloned.PublicClone(),
				pmp_1.ProposalRewardAccountID(prop.ID),
				member.UserAccountID(user),
				payout,
				fmt.Sprintf("reviewer reward for proposal %v", prop.ID),
			)
		}
	}

	// send remainder to matching fund
	receipts := metric.Receipts{}
	rewardAccount := account.Get_Local(ctx, cloned.PublicClone(), pmp_1.ProposalRewardAccountID(prop.ID))
	remainder := rewardAccount.Balance(account.PluralAsset).Quantity
	donation := account.H(account.PluralAsset, 0.0)
	if remainder > 0 {
		donation = account.H(
			account.PluralAsset,
			remainder,
		)
		account.Transfer_StageOnly(
			ctx,
			cloned.PublicClone(),
			pmp_1.ProposalRewardAccountID(prop.ID),
			pmp_0.MatchingPoolAccountID,
			donation,
			fmt.Sprintf("donation to matching fund for proposal %v", prop.ID),
		)
		receipts = append(
			receipts,
			metric.OneReceipt(
				pmp_0.MatchingPoolAccountID.MetricAccountID(),
				metric.ReceiptTypeDonation,
				donation.MetricHolding(),
			)...,
		)
	}

	rewards.Sort()
	return rewards, receipts, donation.Quantity
}

func calcRealBounty(
	priorityFunds float64,
	matchingFunds float64,
	projectedBounty float64,

) (
	projectedBountyFromPriorityFunds float64,
	projectedBountyFromMatchingFunds float64,
	donationFromPriorityFunds float64,

) {

	// ensure no negatives
	priorityFunds = max(priorityFunds, 0)
	matchingFunds = max(matchingFunds, 0)
	projectedBounty = max(projectedBounty, 0)

	projectedBountyFromPriorityFunds = min(projectedBounty, priorityFunds)
	bountyDeficit := max(projectedBounty-projectedBountyFromPriorityFunds, 0)
	if bountyDeficit == 0.0 {
		donationFromPriorityFunds = max(priorityFunds-projectedBountyFromPriorityFunds, 0)
	} else {
		projectedBountyFromMatchingFunds = min(bountyDeficit, matchingFunds)
	}

	return
}
