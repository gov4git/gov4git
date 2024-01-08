package sv

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/v2/proto/account"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/member"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
)

func (qv SV) Cancel(
	ctx context.Context,
	govOwner gov.OwnerCloned,
	ad *ballotproto.Ad,
	tally *ballotproto.Tally,
) git.Change[form.Map, ballotproto.Outcome] {

	// refund users
	refunded := map[member.User]account.Holding{}
	for user, spent := range tally.Charges {
		refund := account.H(account.PluralAsset, spent)
		account.Transfer_StageOnly(
			ctx,
			govOwner.PublicClone(),
			ballotproto.BallotEscrowAccountID(ad.ID),
			member.UserAccountID(user),
			refund,
			fmt.Sprintf("refund from cancelling ballot %v", ad.ID),
		)
		refunded[user] = refund
	}

	return git.NewChange(
		fmt.Sprintf("cancelled ballot %v and refunded voters", ad.ID),
		"ballot_qv_cancel",
		form.Map{"id": ad.ID},
		ballotproto.Outcome{
			Summary:      "cancelled",
			Scores:       tally.Scores,
			ScoresByUser: tally.ScoresByUser,
			Refunded:     refunded,
		},
		nil,
	)
}
