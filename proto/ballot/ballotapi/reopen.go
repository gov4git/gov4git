package ballotapi

import (
	"context"

	"github.com/gov4git/gov4git/v2/proto"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotio"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func Reopen(
	ctx context.Context,
	govAddr gov.OwnerAddress,
	ballotName ballotproto.BallotName,
) git.Change[form.Map, form.None] {

	govCloned := gov.CloneOwner(ctx, govAddr)
	chg := Reopen_StageOnly(ctx, govCloned, ballotName)
	proto.Commit(ctx, govCloned.Public.Tree(), chg)
	govCloned.Public.Push(ctx)
	return chg
}

func Reopen_StageOnly(
	ctx context.Context,
	govCloned gov.OwnerCloned,
	ballotName ballotproto.BallotName,
) git.Change[form.Map, form.None] {

	govTree := govCloned.Public.Tree()

	// verify ad and strategy are present
	ad, strat := ballotio.LoadStrategy(ctx, govTree, ballotName)
	must.Assertf(ctx, ad.Closed, "ballot is not closed")
	must.Assertf(ctx, !ad.Cancelled, "ballot was cancelled")

	tally := loadTally_Local(ctx, govTree, ballotName)
	chg := strat.Reopen(ctx, govCloned, &ad, &tally)

	// remove prior outcome
	openOutcomeNS := ballotproto.BallotPath(ballotName).Append(ballotproto.OutcomeFilebase)
	_, err := git.TreeRemove(ctx, govTree, openOutcomeNS)
	must.NoError(ctx, err)

	// write state
	ad.Closed = false
	ad.Cancelled = false
	openAdNS := ballotproto.BallotPath(ballotName).Append(ballotproto.AdFilebase)
	git.ToFileStage(ctx, govTree, openAdNS, ad)

	return chg
}