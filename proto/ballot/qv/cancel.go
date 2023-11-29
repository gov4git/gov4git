package qv

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/proto/balance"
	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
)

func (qv QV) Cancel(
	ctx context.Context,
	govOwner gov.OwnerCloned,
	ad *common.Advertisement,
	tally *common.Tally,
) git.Change[form.Map, common.Outcome] {

	// refund users
	for user, spent := range tally.Charges {
		balance.Add_StageOnly(ctx, govOwner.PublicClone(), user, VotingCredits, spent) //XXX: deprecated
		// account.Transfer_StageOnly(
		// 	ctx,
		// 	govOwner.PublicClone(),
		// 	common.BallotEscrowAccountID(ad.Name),
		// 	member.UserAccountID(user),
		// 	account.H(account.PluralAsset, spent),
		// )
	}

	return git.NewChange(
		fmt.Sprintf("cancelled ballot %v and refunded voters", ad.Name),
		"ballot_qv_cancel",
		form.Map{"ballot_name": ad.Name},
		common.Outcome{
			Summary: "cancelled",
			Scores:  tally.Scores,
		},
		nil,
	)
}
