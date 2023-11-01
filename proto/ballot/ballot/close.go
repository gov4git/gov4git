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

func Close(
	ctx context.Context,
	govAddr gov.OrganizerAddress,
	ballotName ns.NS,
	cancel bool,
) git.Change[form.Map, common.Outcome] {

	govCloned := id.CloneOwner(ctx, id.OwnerAddress(govAddr))
	chg := Close_StageOnly(ctx, govAddr, govCloned, ballotName, cancel)
	proto.Commit(ctx, govCloned.Public.Tree(), chg)
	govCloned.Public.Push(ctx)
	return chg
}

func Close_StageOnly(
	ctx context.Context,
	govAddr gov.OrganizerAddress,
	govCloned id.OwnerCloned,
	ballotName ns.NS,
	cancel bool,
) git.Change[form.Map, common.Outcome] {

	govTree := govCloned.Public.Tree()

	// verify ad and strategy are present
	ad, strat := load.LoadStrategy(ctx, govTree, ballotName)
	must.Assertf(ctx, !ad.Closed, "ballot already closed")

	tally := LoadTally(ctx, govTree, ballotName)

	var chg git.Change[map[string]form.Form, common.Outcome]
	if cancel {
		chg = strat.Cancel(ctx, govCloned, &ad, &tally)
	} else {
		chg = strat.Close(ctx, govCloned, &ad, &tally)
	}

	// write outcome
	openOutcomeNS := common.BallotPath(ballotName).Append(common.OutcomeFilebase)
	git.ToFileStage(ctx, govTree, openOutcomeNS, chg.Result)

	// write state
	ad.Closed = true
	ad.Cancelled = cancel
	openAdNS := common.BallotPath(ballotName).Append(common.AdFilebase)
	git.ToFileStage(ctx, govTree, openAdNS, ad)

	return chg
}
