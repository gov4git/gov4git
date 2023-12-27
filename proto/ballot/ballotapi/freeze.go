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
	govAddr gov.OwnerAddress,
	ballotName ballotproto.BallotName,

) git.ChangeNoResult {

	govCloned := gov.CloneOwner(ctx, govAddr)
	chg := Freeze_StageOnly(ctx, govCloned, ballotName)
	proto.Commit(ctx, govCloned.Public.Tree(), chg)
	govCloned.Public.Push(ctx)
	return chg
}

func Freeze_StageOnly(
	ctx context.Context,
	cloned gov.OwnerCloned,
	ballotName ballotproto.BallotName,

) git.ChangeNoResult {

	govTree := cloned.Public.Tree()

	ad, _ := ballotio.LoadStrategy(ctx, govTree, ballotName)

	must.Assertf(ctx, !ad.Closed, "ballot is closed")
	must.Assertf(ctx, !ad.Frozen, "ballot already frozen")

	ad.Frozen = true

	// write updated ad
	adNS := ballotproto.BallotPath(ballotName).Append(ballotproto.AdFilebase)
	git.ToFileStage(ctx, govTree, adNS, ad)

	trace.Log_StageOnly(ctx, cloned.PublicClone(), &trace.Event{
		Op:     "ballot_freeze",
		Args:   trace.M{"name": ballotName},
		Result: trace.M{"ad": ad},
	})

	return git.NewChangeNoResult(fmt.Sprintf("Freeze ballot %v", ballotName), "ballot_freeze")
}

func IsFrozen_Local(
	ctx context.Context,
	cloned gov.Cloned,
	ballotName ballotproto.BallotName,

) bool {

	ad, _ := ballotio.LoadStrategy(ctx, cloned.Tree(), ballotName)
	return ad.Frozen
}
