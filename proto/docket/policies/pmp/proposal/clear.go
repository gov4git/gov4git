package proposal

import (
	"context"
	"math"
	"sort"

	"github.com/gov4git/gov4git/proto/account"
	"github.com/gov4git/gov4git/proto/ballot/ballot"
	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/docket/ops"
	"github.com/gov4git/gov4git/proto/docket/policies/pmp"
	"github.com/gov4git/gov4git/proto/docket/schema"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/member"
	"github.com/gov4git/lib4git/base"
)

func loadResolvedConcerns(
	ctx context.Context,
	cloned gov.OwnerCloned,
	prop schema.Motion,

) schema.Motions {

	resolved := schema.Motions{}
	for _, ref := range prop.RefTo {
		if ref.Type != pmp.ResolvesRefType {
			continue
		}
		con := ops.LookupMotion_Local(ctx, cloned.PublicClone(), ref.To)
		if !con.IsConcern() {
			continue
		}
		if con.ID == prop.ID {
			base.Errorf("bug: concern and proposal with same id")
			continue
		}
		if con.Closed {
			continue
		}
		resolved = append(resolved, con)
	}
	return resolved
}

func closeResolvedConcerns(
	ctx context.Context,
	cloned gov.OwnerCloned,
	prop schema.Motion,
	cons schema.Motions,

) account.Holding {

	for _, con := range cons {
		// close resolved concerns, and transfer concern escrows to proposal-owned bounty account
		ops.CloseMotion_StageOnly(
			ctx,
			cloned,
			con.ID,
			pmp.ProposalBountyAccountID(prop.ID), // account to send bounty to
			prop,                                 // proposal that resolves the issue
		)
	}

	return account.Get_Local(
		ctx,
		cloned.PublicClone(),
		pmp.ProposalBountyAccountID(prop.ID),
	).Assets.Balance(account.PluralAsset)
}

func loadPropApprovalPollTally(
	ctx context.Context,
	cloned gov.Cloned,
	prop schema.Motion,

) common.AdTally {

	pollName := pmp.ProposalApprovalPollName(prop.ID)
	return ballot.Show_Local(ctx, cloned.Tree(), pollName)
}

func disberseRewards(
	ctx context.Context,
	cloned gov.OwnerCloned,
	prop schema.Motion,

) Rewards {

	rewards := Rewards{}
	adt := loadPropApprovalPollTally(ctx, cloned.PublicClone(), prop)

	// get reward account balance
	// totalWinnings := account.Get_Local(
	// 	ctx,
	// 	cloned.PublicClone(),
	// 	pmp.ProposalRewardAccountID(prop.ID),
	// ).Assets.Balance(account.PluralAsset).Quantity

	// compute reward distribution
	rewardFund := 0.0                      // total credits spent on negative votes
	totalCut := 0.0                        // sum of all positive votes
	winnerCut := map[member.User]float64{} // positive votes per user
	for user, choices := range adt.Tally.ScoresByUser {
		ss := choices[pmp.ProposalBallotChoice]
		if ss.Score <= 0.0 {
			// compute total credits spent on negative votes
			rewardFund += math.Abs(ss.Strength)
		} else {
			totalCut += ss.Score
			winnerCut[user] = ss.Score
		}
	}

	// payout winnings
	for user, choices := range adt.Tally.ScoresByUser {
		ss := choices[pmp.ProposalBallotChoice]
		if ss.Score > 0.0 {
			payout := account.H(
				account.PluralAsset,
				math.Abs(ss.Strength)+rewardFund*winnerCut[user]/totalCut,
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
				pmp.ProposalRewardAccountID(prop.ID),
				member.UserAccountID(user),
				payout,
			)
		}
	}

	rewards.Sort()
	return rewards
}

type Reward struct {
	To     member.User     `json:"to"`
	Amount account.Holding `json:"amount"`
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
