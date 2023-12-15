package ballot

import (
	"context"

	"github.com/gov4git/gov4git/v2/proto"
	"github.com/gov4git/gov4git/v2/proto/ballot/common"
	"github.com/gov4git/gov4git/v2/proto/ballot/load"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/history"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func Cancel(
	ctx context.Context,
	govAddr gov.OwnerAddress,
	ballotName common.BallotName,
) git.Change[form.Map, common.Outcome] {

	govCloned := gov.CloneOwner(ctx, govAddr)
	chg := Cancel_StageOnly(ctx, govCloned, ballotName)
	proto.Commit(ctx, govCloned.Public.Tree(), chg)
	govCloned.Public.Push(ctx)
	return chg
}

func Cancel_StageOnly(
	ctx context.Context,
	cloned gov.OwnerCloned,
	ballotName common.BallotName,
) git.Change[form.Map, common.Outcome] {

	t := cloned.Public.Tree()

	// verify ad and strategy are present
	ad, strat := load.LoadStrategy(ctx, t, ballotName)
	must.Assertf(ctx, !ad.Closed, "ballot already closed")

	tally := LoadTally(ctx, t, ballotName)

	var chg git.Change[map[string]form.Form, common.Outcome]
	chg = strat.Cancel(ctx, cloned, &ad, &tally)

	// write outcome
	openOutcomeNS := common.BallotPath(ballotName).Append(common.OutcomeFilebase)
	git.ToFileStage(ctx, t, openOutcomeNS, chg.Result)

	// write state
	ad.Closed = true
	ad.Cancelled = true
	openAdNS := common.BallotPath(ballotName).Append(common.AdFilebase)
	git.ToFileStage(ctx, t, openAdNS, ad)

	// log
	history.Log_StageOnly(ctx, cloned.PublicClone(), &history.Event{
		Op: &history.Op{
			Op:     "ballot_cancel",
			Args:   history.M{"name": ballotName},
			Result: history.M{"ad": ad, "outcome": chg.Result},
		},
	})

	return chg
}
