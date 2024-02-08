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

func Freeze(
	ctx context.Context,
	addr gov.OwnerAddress,
	id ballotproto.BallotID,

) git.ChangeNoResult {

	cloned := gov.CloneOwner(ctx, addr)
	chg := Freeze_StageOnly(ctx, cloned, id)
	proto.Commit(ctx, cloned.Public.Tree(), chg)
	cloned.Public.Push(ctx)
	return chg
}

func Freeze_StageOnly(
	ctx context.Context,
	cloned gov.OwnerCloned,
	id ballotproto.BallotID,

) git.ChangeNoResult {

	t := cloned.Public.Tree()

	ad := ballotio.LoadAd_Local(ctx, t, id)

	must.Assertf(ctx, !ad.Closed, "ballot is closed")
	must.Assertf(ctx, !ad.Frozen, "ballot already frozen")

	ad.Frozen = true

	// write updated ad
	git.ToFileStage(ctx, t, id.AdNS(), ad)

	trace.Log_StageOnly(ctx, cloned.PublicClone(), &trace.Event{
		Op:     "ballot_freeze",
		Args:   trace.M{"id": id},
		Result: trace.M{"ad": ad},
	})

	return git.NewChangeNoResult(fmt.Sprintf("Freeze ballot %v", id), "ballot_freeze")
}

func IsFrozen_Local(
	ctx context.Context,
	cloned gov.Cloned,
	id ballotproto.BallotID,

) bool {

	ad := ballotio.LoadAd_Local(ctx, cloned.Tree(), id)
	return ad.Frozen
}
