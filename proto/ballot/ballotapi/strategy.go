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

func LoadStrategyState[SS form.Form](
	ctx context.Context,
	addr gov.Address,
	id ballotproto.BallotID,

) SS {

	cloned := gov.Clone(ctx, addr)
	return LoadStrategyState_Local[SS](ctx, cloned, id)
}

func LoadStrategyState_Local[SS form.Form](
	ctx context.Context,
	cloned gov.Cloned,
	id ballotproto.BallotID,

) SS {

	t := cloned.Tree()
	return git.FromFile[SS](ctx, t, id.StrategyNS())
}

func SaveStrategyState[SS form.Form](
	ctx context.Context,
	addr gov.Address,
	id ballotproto.BallotID,
	strategyState SS,

) {

	cloned := gov.Clone(ctx, addr)
	SaveStrategyState_StageOnly[SS](ctx, cloned, id, strategyState)
	proto.Commitf(ctx, cloned, "ballot_save_strategy_state", "update ballot strategy state")
	cloned.Push(ctx)
}

func SaveStrategyState_StageOnly[SS form.Form](
	ctx context.Context,
	cloned gov.Cloned,
	id ballotproto.BallotID,
	strategyState SS,

) {

	t := cloned.Tree()
	ad, _ := ballotio.LoadStrategy(ctx, t, id)
	must.Assertf(ctx, !ad.Closed, "ballot already closed")
	git.ToFileStage[SS](ctx, t, id.StrategyNS(), strategyState)
}
