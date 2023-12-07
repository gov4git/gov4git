package ballot

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/account"
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
	escrowTo account.AccountID,

) git.Change[form.Map, common.Outcome] {

	govCloned := gov.CloneOwner(ctx, govAddr)
	chg := Close_StageOnly(ctx, govCloned, ballotName, escrowTo)
	proto.Commit(ctx, govCloned.Public.Tree(), chg)
	govCloned.Public.Push(ctx)
	return chg
}

func Close_StageOnly(
	ctx context.Context,
	cloned gov.OwnerCloned,
	ballotName common.BallotName,
	escrowTo account.AccountID,

) git.Change[form.Map, common.Outcome] {

	t := cloned.Public.Tree()

	// verify ad and strategy are present
	ad, strat := load.LoadStrategy(ctx, t, ballotName)
	must.Assertf(ctx, !ad.Closed, "ballot already closed")

	tally := LoadTally(ctx, t, ballotName)

	var chg git.Change[map[string]form.Form, common.Outcome]
	chg = strat.Close(ctx, cloned, &ad, &tally)

	// write outcome
	openOutcomeNS := common.BallotPath(ballotName).Append(common.OutcomeFilebase)
	git.ToFileStage(ctx, t, openOutcomeNS, chg.Result)

	// write state
	ad.Closed = true
	ad.Cancelled = false
	openAdNS := common.BallotPath(ballotName).Append(common.AdFilebase)
	git.ToFileStage(ctx, t, openAdNS, ad)

	// transfer escrow
	escrowAccountID := common.BallotEscrowAccountID(ballotName)
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
