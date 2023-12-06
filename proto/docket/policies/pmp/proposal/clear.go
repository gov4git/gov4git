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
	"github.com/gov4git/gov4git/proto/docket/policies/pmp/concern"
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
	for _, ref := range prop.RefBy {
		if ref.Type != concern.ResolvesRefType {
			continue
		}
		concern := ops.LookupMotion_Local(ctx, cloned.PublicClone(), prop.ID)
		if !concern.IsConcern() {
			continue
		}
		if concern.ID == prop.ID {
			base.Errorf("bug: concern and proposal with same id")
			continue
		}
		if concern.Closed {
			continue
		}
		resolved = append(resolved, concern)
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
			pmp.ProposalBountyAccountID(prop.ID),
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

	// compute XXX
	rewardFund := 0.0
	totalCut := 0.0
	winnerCut := map[member.User]float64{}
	for user, choices := range adt.Tally.ScoresByUser {
		ss := choices[pmp.ProposalBallotChoice]
		if ss.Score == 0.0 {
			rewardFund += math.Abs(ss.Strength)
		} else {
			totalCut += ss.Score
			winnerCut[user] = ss.Score
		}
	}

	// payout winnings
	// XXX: precion problems, make sure less is withdrawn than the actual reward account
	for user, choices := range adt.Tally.ScoresByUser {
		ss := choices[pmp.ProposalBallotChoice]
		if ss.Score > 0.0 {
			payout := account.H(
				account.PluralAsset,
				math.Abs(ss.Strength)+winnerCut[user]/totalCut,
			)
			rewards = append(rewards,
				Reward{
					To:     user,
					Amount: payout,
				},
			)
			// XXX: transfer reward
			// account.Transfer_StageOnly(
			// 	ctx,
			// 	cloned,
			// 	XXX,
			// 	member.UserAccountID(user),
			// 	payout,
			// )
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
