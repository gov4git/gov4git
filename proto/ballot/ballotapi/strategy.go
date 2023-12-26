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
	addr gov.OwnerAddress,
	name ballotproto.BallotName,

) SS {

	cloned := gov.CloneOwner(ctx, addr)
	return LoadStrategyState_Local[SS](ctx, cloned, name)
}

func LoadStrategyState_Local[SS form.Form](
	ctx context.Context,
	cloned gov.OwnerCloned,
	name ballotproto.BallotName,

) SS {

	t := cloned.Public.Tree()
	return git.FromFile[SS](ctx, t, name.StrategyNS())
}

func SaveStrategyState[SS form.Form](
	ctx context.Context,
	addr gov.OwnerAddress,
	name ballotproto.BallotName,
	strategyState SS,

) {

	cloned := gov.CloneOwner(ctx, addr)
	SaveStrategyState_StageOnly[SS](ctx, cloned, name, strategyState)
	proto.Commitf(ctx, cloned.Public, "ballot_save_strategy_state", "update ballot strategy state")
	cloned.Public.Push(ctx)
}

func SaveStrategyState_StageOnly[SS form.Form](
	ctx context.Context,
	cloned gov.OwnerCloned,
	name ballotproto.BallotName,
	strategyState SS,

) {

	t := cloned.Public.Tree()
	ad, _ := ballotio.LoadStrategy(ctx, t, name)
	must.Assertf(ctx, !ad.Closed, "ballot already closed")
	git.ToFileStage[SS](ctx, t, name.StrategyNS(), strategyState)
}
