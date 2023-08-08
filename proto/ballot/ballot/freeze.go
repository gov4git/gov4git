package ballot

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/ballot/load"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
	"github.com/gov4git/lib4git/ns"
)

func Freeze(
	ctx context.Context,
	govAddr gov.OrganizerAddress,
	ballotName ns.NS,
) git.ChangeNoResult {

	govCloned := id.CloneOwner(ctx, id.OwnerAddress(govAddr))
	chg := FreezeStageOnly(ctx, govAddr, govCloned, ballotName)
	proto.Commit(ctx, govCloned.Public.Tree(), chg)
	govCloned.Public.Push(ctx)
	return chg
}

func FreezeStageOnly(
	ctx context.Context,
	govAddr gov.OrganizerAddress,
	govCloned id.OwnerCloned,
	ballotName ns.NS,
) git.ChangeNoResult {

	govTree := govCloned.Public.Tree()

	// verify ad and strategy are present
	ad, _ := load.LoadStrategy(ctx, govTree, ballotName, false)

	// verify ballot is not already frozen
	if ad.Frozen {
		must.Errorf(ctx, "ballot already frozen")
	}
	ad.Frozen = true

	// write updated ad
	openAdNS := common.OpenBallotNS(ballotName).Sub(common.AdFilebase)
	git.ToFileStage(ctx, govTree, openAdNS.Path(), ad)

	return git.NewChangeNoResult(fmt.Sprintf("Freeze ballot %v", ballotName), "ballot_freeze")
}
