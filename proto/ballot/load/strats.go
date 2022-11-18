package load

import (
	"context"

	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/ballot/qv"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
	"github.com/gov4git/lib4git/ns"
)

func LoadStrategy(
	ctx context.Context,
	govTree *git.Tree,
	ballotName ns.NS,
) (common.Advertisement, common.Strategy) {

	// read ad
	openAdNS := common.OpenBallotNS(ballotName).Sub(common.AdFilebase)
	ad := git.FromFile[common.Advertisement](ctx, govTree, openAdNS.Path())

	// read strategy
	openStrategyNS := common.OpenBallotNS(ballotName).Sub(common.StrategyFilebase)
	switch ad.Strategy {
	case qv.PriorityPollName:
		return ad, git.FromFile[qv.PriorityPoll](ctx, govTree, openStrategyNS.Path())
	default:
		must.Errorf(ctx, "unkonwn ballot strategy %v", ad.Strategy)
		panic("unreachable")
	}
}
