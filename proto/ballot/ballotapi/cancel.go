package ballotapi

import (
	"context"

	"github.com/gov4git/gov4git/v2/proto"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotio"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/history/trace"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func Cancel(
	ctx context.Context,
	govAddr gov.OwnerAddress,
	ballotName ballotproto.BallotName,
) git.Change[form.Map, ballotproto.Outcome] {

	govCloned := gov.CloneOwner(ctx, govAddr)
	chg := Cancel_StageOnly(ctx, govCloned, ballotName)
	proto.Commit(ctx, govCloned.Public.Tree(), chg)
	govCloned.Public.Push(ctx)
	return chg
}

func Cancel_StageOnly(
	ctx context.Context,
	cloned gov.OwnerCloned,
	ballotName ballotproto.BallotName,
) git.Change[form.Map, ballotproto.Outcome] {

	t := cloned.Public.Tree()

	// verify ad and strategy are present
	ad, strat := ballotio.LoadStrategy(ctx, t, ballotName)
	must.Assertf(ctx, !ad.Closed, "ballot already closed")

	tally := loadTally_Local(ctx, t, ballotName)

	var chg git.Change[map[string]form.Form, ballotproto.Outcome]
	chg = strat.Cancel(ctx, cloned, &ad, &tally)

	// write outcome
	openOutcomeNS := ballotproto.BallotPath(ballotName).Append(ballotproto.OutcomeFilebase)
	git.ToFileStage(ctx, t, openOutcomeNS, chg.Result)

	// write state
	ad.Closed = true
	ad.Cancelled = true
	openAdNS := ballotproto.BallotPath(ballotName).Append(ballotproto.AdFilebase)
	git.ToFileStage(ctx, t, openAdNS, ad)

	// log
	trace.Log_StageOnly(ctx, cloned.PublicClone(), &trace.Event{
		Op:     "ballot_cancel",
		Args:   trace.M{"name": ballotName},
		Result: trace.M{"ad": ad, "outcome": chg.Result},
	})

	return chg
}
