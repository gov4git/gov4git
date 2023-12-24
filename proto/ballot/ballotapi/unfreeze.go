package ballotapi

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/v2/proto"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotio"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/history"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func Unfreeze(
	ctx context.Context,
	govAddr gov.OwnerAddress,
	ballotName ballotproto.BallotName,
) git.ChangeNoResult {

	govCloned := gov.CloneOwner(ctx, govAddr)
	chg := Unfreeze_StageOnly(ctx, govCloned, ballotName)
	proto.Commit(ctx, govCloned.Public.Tree(), chg)
	govCloned.Public.Push(ctx)
	return chg
}

func Unfreeze_StageOnly(
	ctx context.Context,
	cloned gov.OwnerCloned,
	ballotName ballotproto.BallotName,
) git.ChangeNoResult {

	govTree := cloned.Public.Tree()

	// verify ad and strategy are present
	ad, _ := ballotio.LoadStrategy(ctx, govTree, ballotName)

	must.Assertf(ctx, !ad.Closed, "ballot is closed")
	must.Assertf(ctx, ad.Frozen, "ballot is not frozen")

	ad.Frozen = false

	// write updated ad
	adNS := ballotproto.BallotPath(ballotName).Append(ballotproto.AdFilebase)
	git.ToFileStage(ctx, govTree, adNS, ad)

	history.Log_StageOnly(ctx, cloned.PublicClone(), &history.Event{
		Op: &history.Op{
			Op:     "ballot_freeze",
			Args:   history.M{"name": ballotName},
			Result: history.M{"ad": ad},
		},
	})

	return git.NewChangeNoResult(fmt.Sprintf("Unfreeze ballot %v", ballotName), "ballot_unfreeze")
}
