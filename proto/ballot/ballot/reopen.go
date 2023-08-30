package ballot

import (
	"context"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/ballot/load"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
	"github.com/gov4git/lib4git/ns"
)

func Reopen(
	ctx context.Context,
	govAddr gov.OrganizerAddress,
	ballotName ns.NS,
) git.Change[form.Map, form.None] {

	govCloned := id.CloneOwner(ctx, id.OwnerAddress(govAddr))
	chg := Reopen_StageOnly(ctx, govAddr, govCloned, ballotName)
	proto.Commit(ctx, govCloned.Public.Tree(), chg)
	govCloned.Public.Push(ctx)
	return chg
}

func Reopen_StageOnly(
	ctx context.Context,
	govAddr gov.OrganizerAddress,
	govCloned id.OwnerCloned,
	ballotName ns.NS,
) git.Change[form.Map, form.None] {

	govTree := govCloned.Public.Tree()

	// verify ad and strategy are present
	ad, strat := load.LoadStrategy(ctx, govTree, ballotName)
	must.Assertf(ctx, ad.Closed, "ballot is not closed")

	tally := LoadTally(ctx, govTree, ballotName)
	chg := strat.Reopen(ctx, govCloned, &ad, &tally)

	// remove prior outcome
	openOutcomeNS := common.BallotPath(ballotName).Sub(common.OutcomeFilebase)
	_, err := govTree.Remove(openOutcomeNS.Path())
	must.NoError(ctx, err)

	// write state
	ad.Closed = false
	ad.Cancelled = false
	openAdNS := common.BallotPath(ballotName).Sub(common.AdFilebase)
	git.ToFileStage(ctx, govTree, openAdNS.Path(), ad)

	return chg
}
