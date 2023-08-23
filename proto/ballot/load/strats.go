package load

import (
	"context"

	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/ballot/qv"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func LoadStrategy(ctx context.Context, govTree *git.Tree, ballotName common.BallotName) (common.Advertisement, common.Strategy) {

	// read ad
	adNS := common.BallotPath(ballotName).Append(common.AdFilebase)
	ad := git.FromFile[common.Advertisement](ctx, govTree, adNS)

	// read strategy
	strategyNS := common.BallotPath(ballotName).Append(common.StrategyFilebase)
	switch ad.Strategy {
	case qv.QVStrategyName:
		return ad, git.FromFile[qv.QV](ctx, govTree, strategyNS)
	default:
		must.Errorf(ctx, "unkonwn ballot strategy %v", ad.Strategy)
		panic("unreachable")
	}
}
