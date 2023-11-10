package qv

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
)

func (qv QV) Close(
	ctx context.Context,
	govOwner gov.GovOwnerCloned,
	ad *common.Advertisement,
	tally *common.Tally,
) git.Change[form.Map, common.Outcome] {

	return git.NewChange(
		fmt.Sprintf("closed ballot %v", ad.Name),
		"ballot_qv_close",
		form.Map{"ballot_name": ad.Name},
		common.Outcome{
			Summary: "closed",
			Scores:  tally.Scores,
		},
		nil,
	)
}
