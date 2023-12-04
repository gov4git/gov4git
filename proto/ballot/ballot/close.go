package ballot

import (
	"context"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/ballot/load"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/history"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func Close(
	ctx context.Context,
	govAddr gov.OwnerAddress,
	ballotName common.BallotName,
	cancel bool,
) git.Change[form.Map, common.Outcome] {

	govCloned := gov.CloneOwner(ctx, govAddr)
	chg := Close_StageOnly(ctx, govCloned, ballotName, cancel)
	proto.Commit(ctx, govCloned.Public.Tree(), chg)
	govCloned.Public.Push(ctx)
	return chg
}

func Close_StageOnly(
	ctx context.Context,
	cloned gov.OwnerCloned,
	ballotName common.BallotName,
	cancel bool,
) git.Change[form.Map, common.Outcome] {

	govTree := cloned.Public.Tree()

	// verify ad and strategy are present
	ad, strat := load.LoadStrategy(ctx, govTree, ballotName)
	must.Assertf(ctx, !ad.Closed, "ballot already closed")

	tally := LoadTally(ctx, govTree, ballotName)

	var chg git.Change[map[string]form.Form, common.Outcome]
	if cancel {
		chg = strat.Cancel(ctx, cloned, &ad, &tally)
	} else {
		chg = strat.Close(ctx, cloned, &ad, &tally)
	}

	// write outcome
	openOutcomeNS := common.BallotPath(ballotName).Append(common.OutcomeFilebase)
	git.ToFileStage(ctx, govTree, openOutcomeNS, chg.Result)

	// write state
	ad.Closed = true
	ad.Cancelled = cancel
	openAdNS := common.BallotPath(ballotName).Append(common.AdFilebase)
	git.ToFileStage(ctx, govTree, openAdNS, ad)

	// log
	if cancel {
		history.Log_StageOnly(ctx, cloned.PublicClone(), &history.Event{
			Op: &history.Op{
				Op:     "ballot_cancel",
				Args:   history.M{"name": ballotName},
				Result: history.M{"ad": ad, "outcome": chg.Result},
			},
		})
	} else {
		history.Log_StageOnly(ctx, cloned.PublicClone(), &history.Event{
			Op: &history.Op{
				Op:     "ballot_close",
				Args:   history.M{"name": ballotName},
				Result: history.M{"ad": ad, "outcome": chg.Result},
			},
		})
	}

	return chg
}
