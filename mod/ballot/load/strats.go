package load

import (
	"context"

	"github.com/gov4git/gov4git/mod/ballot/proto"
	"github.com/gov4git/gov4git/mod/ballot/qv"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
	"github.com/gov4git/lib4git/ns"
)

func LoadStrategy(
	ctx context.Context,
	govTree *git.Tree,
	ballotName ns.NS,
) (proto.Advertisement, proto.Strategy) {

	// read ad
	openAdNS := proto.OpenBallotNS(ballotName).Sub(proto.AdFilebase)
	ad := git.FromFile[proto.Advertisement](ctx, govTree, openAdNS.Path())

	// read strategy
	openStrategyNS := proto.OpenBallotNS(ballotName).Sub(proto.StrategyFilebase)
	switch ad.Strategy {
	case qv.PriorityPollName:
		return ad, git.FromFile[qv.PriorityPoll](ctx, govTree, openStrategyNS.Path())
	default:
		must.Errorf(ctx, "unkonwn ballot strategy %v", ad.Strategy)
		panic("unreachable")
	}
}
