package ballotapi

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/v2/proto"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotio"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/history/trace"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func Unfreeze(
	ctx context.Context,
	addr gov.OwnerAddress,
	id ballotproto.BallotID,

) git.ChangeNoResult {

	cloned := gov.CloneOwner(ctx, addr)
	chg := Unfreeze_StageOnly(ctx, cloned, id)
	proto.Commit(ctx, cloned.Public.Tree(), chg)
	cloned.Public.Push(ctx)
	return chg
}

func Unfreeze_StageOnly(
	ctx context.Context,
	cloned gov.OwnerCloned,
	id ballotproto.BallotID,

) git.ChangeNoResult {

	t := cloned.Public.Tree()

	// verify ad is present
	ad := ballotio.LoadAd_Local(ctx, t, id)

	must.Assertf(ctx, !ad.Closed, "ballot is closed")
	must.Assertf(ctx, ad.Frozen, "ballot is not frozen")

	ad.Frozen = false

	// write updated ad
	git.ToFileStage(ctx, t, id.AdNS(), ad)

	trace.Log_StageOnly(ctx, cloned.PublicClone(), &trace.Event{
		Op:     "ballot_freeze",
		Args:   trace.M{"id": id},
		Result: trace.M{"ad": ad},
	})

	return git.NewChangeNoResult(fmt.Sprintf("Unfreeze ballot %v", id), "ballot_unfreeze")
}
