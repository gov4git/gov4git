package ballotapi

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/v2/proto"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotio"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/id"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
)

func Erase(
	ctx context.Context,
	govAddr gov.OwnerAddress,
	ballotID ballotproto.BallotID,

) git.Change[form.Map, bool] {

	cloned := id.CloneOwner(ctx, id.OwnerAddress(govAddr))
	chg := Erase_StageOnly(ctx, cloned, ballotID)
	proto.Commit(ctx, cloned.Public.Tree(), chg)
	cloned.Public.Push(ctx)
	return chg
}

func Erase_StageOnly(
	ctx context.Context,
	cloned id.OwnerCloned,
	ballotID ballotproto.BallotID,

) git.Change[form.Map, bool] {

	t := cloned.Public.Tree()

	// verify ad and policy are present
	ballotio.LoadPolicy(ctx, t, ballotID)

	// erase
	ballotproto.BallotKV.Remove(ctx, ballotproto.BallotNS, cloned.PublicClone().Tree(), ballotID)

	return git.NewChange(
		fmt.Sprintf("Erased ballot %v", ballotID),
		"ballot_erase",
		form.Map{"name": ballotID},
		true,
		nil,
	)
}
