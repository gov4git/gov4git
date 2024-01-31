package proposal

import (
	"context"
	"fmt"
	"math"

	"github.com/gov4git/gov4git/v2/proto/account"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotapi"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/gov"
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

) (resolved motionproto.Motions, escrows []float64) {

	eligible := computeEligibleConcerns(ctx, cloned.PublicClone(), prop)
	for _, ref := range eligible {
		con := motionapi.LookupMotion_Local(ctx, cloned.PublicClone(), ref.To)
		conState := motionapi.LoadPolicyState_Local[*concern.ConcernState](ctx, cloned.PublicClone(), con.ID)
		//
		resolved = append(resolved, con)
		escrows = append(escrows, conState.Escrow())
	}
	return resolved, escrows
}

func computeEligibleConcerns(ctx context.Context, cloned gov.Cloned, prop motionproto.Motion) motionproto.Refs {
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

func disberseRewardsAccepted(
	ctx context.Context,
	cloned gov.OwnerCloned,
	prop motionproto.Motion,

) Rewards {

	rewards := Rewards{}
	adt := loadPropApprovalPollTally(ctx, cloned.PublicClone(), prop)

	// compute rewards
	idealPayout := map[member.User]float64{}
	idealFunds := 0.0
	realFunds := 0.0
	for user, choices := range adt.Tally.ScoresByUser {
		ss := choices[pmp_1.ProposalBallotChoice]
		if ss.Score > 0 {
			idealPayout[user] = 2 * ss.Score
			idealFunds += 2 * ss.Score
		}
		realFunds += math.Abs(ss.Strength)
	}

	// compute shrink factor
	if idealFunds <= 0.0 {
		return nil
	}
	payoutShrinkFactor := min(realFunds/idealFunds, 1.0)
	if payoutShrinkFactor < 1.0 {
		base.Infof("proposal %v has reviewer reward shrink factor %v", prop.ID, payoutShrinkFactor)
	}

	// disberse payouts
	for user, choices := range adt.Tally.ScoresByUser {
		ss := choices[pmp_1.ProposalBallotChoice]
		if ss.Score > 0.0 {
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
	rewardAccount := account.Get_Local(ctx, cloned.PublicClone(), pmp_1.ProposalRewardAccountID(prop.ID))
	remainder := rewardAccount.Balance(account.PluralAsset).Quantity
	if remainder > 0 {
		donation := account.H(
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
	}

	rewards.Sort()
	return rewards
}
