package load

import (
	"context"

	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/ballot/strategies/sv"
	"github.com/gov4git/gov4git/proto/mod"
	"github.com/gov4git/lib4git/git"
)

var strategyRegistry = mod.NewModuleRegistry[common.StrategyName, common.Strategy]()

const (
	QVStrategyName common.StrategyName = "qv"
)

func init() {
	ctx := context.Background()
	strategyRegistry.Set(ctx, QVStrategyName, sv.SV{Scorer: sv.QVScore})
}

func LoadStrategy(
	ctx context.Context,
	t *git.Tree,
	ballotName common.BallotName,
) (common.Advertisement, common.Strategy) {

	adNS := common.BallotPath(ballotName).Append(common.AdFilebase)
	ad := git.FromFile[common.Advertisement](ctx, t, adNS)

	return ad, strategyRegistry.Get(ctx, ad.Strategy)
}
