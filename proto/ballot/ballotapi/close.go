package ballotapi

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/v2/proto"
	"github.com/gov4git/gov4git/v2/proto/account"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotio"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/history"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func Close(
	ctx context.Context,
	govAddr gov.OwnerAddress,
	ballotName ballotproto.BallotName,
	escrowTo account.AccountID,

) git.Change[form.Map, ballotproto.Outcome] {

	govCloned := gov.CloneOwner(ctx, govAddr)
	chg := Close_StageOnly(ctx, govCloned, ballotName, escrowTo)
	proto.Commit(ctx, govCloned.Public.Tree(), chg)
	govCloned.Public.Push(ctx)
	return chg
}

func Close_StageOnly(
	ctx context.Context,
	cloned gov.OwnerCloned,
	ballotName ballotproto.BallotName,
	escrowTo account.AccountID,

) git.Change[form.Map, ballotproto.Outcome] {

	t := cloned.Public.Tree()

	// verify ad and strategy are present
	ad, strat := ballotio.LoadStrategy(ctx, t, ballotName)
	must.Assertf(ctx, !ad.Closed, "ballot already closed")

	tally := LoadTally(ctx, t, ballotName)

	var chg git.Change[map[string]form.Form, ballotproto.Outcome]
	chg = strat.Close(ctx, cloned, &ad, &tally)

	// write outcome
	openOutcomeNS := ballotproto.BallotPath(ballotName).Append(ballotproto.OutcomeFilebase)
	git.ToFileStage(ctx, t, openOutcomeNS, chg.Result)

	// write state
	ad.Closed = true
	ad.Cancelled = false
	openAdNS := ballotproto.BallotPath(ballotName).Append(ballotproto.AdFilebase)
	git.ToFileStage(ctx, t, openAdNS, ad)

	// transfer escrow
	escrowAccountID := ballotproto.BallotEscrowAccountID(ballotName)
	escrowAssets := account.Get_Local(
		ctx,
		cloned.PublicClone(),
		escrowAccountID,
	).Assets
	for _, holding := range escrowAssets {
		account.Transfer_StageOnly(
			ctx,
			cloned.PublicClone(),
			escrowAccountID,
			escrowTo,
			holding,
			fmt.Sprintf("closing ballot %v", ballotName),
		)
	}

	// log
	history.Log_StageOnly(ctx, cloned.PublicClone(), &history.Event{
		Op: &history.Op{
			Op:     "ballot_close",
			Args:   history.M{"name": ballotName},
			Result: history.M{"ad": ad, "outcome": chg.Result},
		},
	})

	return chg
}
