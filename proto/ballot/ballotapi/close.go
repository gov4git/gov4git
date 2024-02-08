package ballotapi

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/v2/proto"
	"github.com/gov4git/gov4git/v2/proto/account"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotio"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/history/trace"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func Close(
	ctx context.Context,
	addr gov.OwnerAddress,
	id ballotproto.BallotID,
	escrowTo account.AccountID,

) git.Change[form.Map, ballotproto.Outcome] {

	cloned := gov.CloneOwner(ctx, addr)
	chg := Close_StageOnly(ctx, cloned, id, escrowTo)
	proto.Commit(ctx, cloned.Public.Tree(), chg)
	cloned.Public.Push(ctx)
	return chg
}

func Close_StageOnly(
	ctx context.Context,
	cloned gov.OwnerCloned,
	id ballotproto.BallotID,
	escrowTo account.AccountID,

) git.Change[form.Map, ballotproto.Outcome] {

	t := cloned.Public.Tree()

	// verify ad and policy are present
	ad, policy := ballotio.LoadAdPolicy_Local(ctx, t, id)
	must.Assertf(ctx, !ad.Closed, "ballot already closed")

	tally := loadTally_Local(ctx, t, id)

	var chg git.Change[map[string]form.Form, ballotproto.Outcome]
	chg = policy.Close(ctx, cloned, &ad, &tally)

	// write outcome
	git.ToFileStage(ctx, t, id.OutcomeNS(), chg.Result)

	// write state
	ad.Closed = true
	ad.Cancelled = false
	git.ToFileStage(ctx, t, id.AdNS(), ad)

	// transfer escrow
	escrowAccountID := ballotproto.BallotEscrowAccountID(id)
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
			fmt.Sprintf("closing ballot %v", id),
		)
	}

	// log
	trace.Log_StageOnly(ctx, cloned.PublicClone(), &trace.Event{
		Op:     "ballot_close",
		Args:   trace.M{"id": id},
		Result: trace.M{"ad": ad, "outcome": chg.Result},
	})

	return chg
}
