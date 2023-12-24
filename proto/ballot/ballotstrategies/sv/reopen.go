package sv

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
)

func (qv SV) Reopen(
	ctx context.Context,
	govOwner gov.OwnerCloned,
	ad *ballotproto.Advertisement,
	tally *ballotproto.Tally,
) git.Change[form.Map, form.None] {

	return git.NewChange(
		fmt.Sprintf("reopened ballot %v", ad.Name),
		"ballot_qv_reopen",
		form.Map{"ballot_name": ad.Name},
		form.None{},
		nil,
	)
}
