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
	closed bool,
) (common.Advertisement, common.Strategy) {

	// read ad
	var adNS ns.NS
	if closed {
		adNS = common.ClosedBallotNS(ballotName).Sub(common.AdFilebase)
	} else {
		adNS = common.OpenBallotNS(ballotName).Sub(common.AdFilebase)
	}
	ad := git.FromFile[common.Advertisement](ctx, govTree, adNS.Path())

	// read strategy
	var strategyNS ns.NS
	if closed {
		strategyNS = common.ClosedBallotNS(ballotName).Sub(common.StrategyFilebase)
	} else {
		strategyNS = common.OpenBallotNS(ballotName).Sub(common.StrategyFilebase)
	}
	switch ad.Strategy {
	case qv.QVStrategyName:
		return ad, git.FromFile[qv.QV](ctx, govTree, strategyNS.Path())
	default:
		must.Errorf(ctx, "unkonwn ballot strategy %v", ad.Strategy)
		panic("unreachable")
	}
}
