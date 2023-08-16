package qv

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/proto/balance"
	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
)

func (qv QV) Cancel(
	ctx context.Context,
	govOwner id.OwnerCloned,
	ad *common.Advertisement,
	tally *common.Tally,
) git.Change[form.Map, common.Outcome] {

	// refund users
	for user, spent := range tally.Charges {
		balance.AddStageOnly(ctx, govOwner.Public.Tree(), user, VotingCredits, spent)
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
