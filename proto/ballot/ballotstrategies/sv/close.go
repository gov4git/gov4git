package sv

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
)

func (qv SV) Close(
	ctx context.Context,
	govOwner gov.OwnerCloned,
	ad *ballotproto.Advertisement,
	tally *ballotproto.Tally,
) git.Change[form.Map, ballotproto.Outcome] {

	return git.NewChange(
		fmt.Sprintf("closed ballot %v", ad.Name),
		"ballot_qv_close",
		form.Map{"ballot_name": ad.Name},
		ballotproto.Outcome{
			Summary:      "closed",
			Scores:       tally.Scores,
			ScoresByUser: tally.ScoresByUser,
			Refunded:     nil,
		},
		nil,
	)
}
