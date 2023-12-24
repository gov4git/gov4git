package ballotio

import (
	"context"

	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotstrategies/sv"
	"github.com/gov4git/gov4git/v2/proto/mod"
	"github.com/gov4git/lib4git/git"
)

var strategyRegistry = mod.NewModuleRegistry[ballotproto.StrategyName, ballotproto.Strategy]()

const (
	QVStrategyName ballotproto.StrategyName = "qv"
)

func init() {
	ctx := context.Background()
	strategyRegistry.Set(ctx, QVStrategyName, sv.SV{Kernel: sv.QVScore{}})
}

func LoadStrategy(
	ctx context.Context,
	t *git.Tree,
	ballotName ballotproto.BallotName,
) (ballotproto.Advertisement, ballotproto.Strategy) {

	adNS := ballotproto.BallotPath(ballotName).Append(ballotproto.AdFilebase)
	ad := git.FromFile[ballotproto.Advertisement](ctx, t, adNS)

	return ad, strategyRegistry.Get(ctx, ad.Strategy)
}

func LookupStrategy(ctx context.Context, name ballotproto.StrategyName) ballotproto.Strategy {
	return strategyRegistry.Get(ctx, name)
}
