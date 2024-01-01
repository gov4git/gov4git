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
	Install(
		ctx,
		QVStrategyName,
		sv.SV{
			Kernel: sv.QVScoreKernel{},
		},
	)
}

func Install(ctx context.Context, name ballotproto.StrategyName, strategy ballotproto.Strategy) {
	strategyRegistry.Set(ctx, name, strategy)
}

func LoadStrategy(
	ctx context.Context,
	t *git.Tree,
	id ballotproto.BallotID,

) (ballotproto.Advertisement, ballotproto.Strategy) {

	ad := git.FromFile[ballotproto.Advertisement](ctx, t, id.AdNS())
	return ad, strategyRegistry.Get(ctx, ad.Strategy)
}

func LookupStrategy(
	ctx context.Context,
	id ballotproto.StrategyName,

) ballotproto.Strategy {

	return strategyRegistry.Get(ctx, id)
}
