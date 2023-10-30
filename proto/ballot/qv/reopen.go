package qv

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
)

func (qv QV) Reopen(
	ctx context.Context,
	govOwner id.OwnerCloned,
	ad *common.Advertisement,
	tally *common.Tally,
) git.Change[form.Map, form.None] {

	// XXX: reopening a cancelled issue must charge user accounts again?

	return git.NewChange(
		fmt.Sprintf("reopened ballot %v", ad.Name),
		"ballot_qv_reopen",
		form.Map{"ballot_name": ad.Name},
		form.None{},
		nil,
	)
}
