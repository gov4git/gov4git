package ballot

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/ballot/load"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func Freeze(
	ctx context.Context,
	govAddr gov.OwnerAddress,
	ballotName common.BallotName,

) git.ChangeNoResult {

	govCloned := gov.CloneOwner(ctx, govAddr)
	chg := Freeze_StageOnly(ctx, govCloned, ballotName)
	proto.Commit(ctx, govCloned.Public.Tree(), chg)
	govCloned.Public.Push(ctx)
	return chg
}

func Freeze_StageOnly(
	ctx context.Context,
	govCloned gov.OwnerCloned,
	ballotName common.BallotName,

) git.ChangeNoResult {

	govTree := govCloned.Public.Tree()

	ad, _ := load.LoadStrategy(ctx, govTree, ballotName)

	must.Assertf(ctx, !ad.Closed, "ballot is closed")
	must.Assertf(ctx, !ad.Frozen, "ballot already frozen")

	ad.Frozen = true

	// write updated ad
	adNS := common.BallotPath(ballotName).Append(common.AdFilebase)
	git.ToFileStage(ctx, govTree, adNS, ad)

	return git.NewChangeNoResult(fmt.Sprintf("Freeze ballot %v", ballotName), "ballot_freeze")
}

func IsFrozen_Local(
	ctx context.Context,
	cloned gov.Cloned,
	ballotName common.BallotName,

) bool {

	ad, _ := load.LoadStrategy(ctx, cloned.Tree(), ballotName)
	return ad.Frozen
}
