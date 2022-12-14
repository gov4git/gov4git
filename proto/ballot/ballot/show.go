package ballot

import (
	"context"

	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/ballot/load"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
	"github.com/gov4git/lib4git/ns"
)

func Show(
	ctx context.Context,
	govAddr gov.GovAddress,
	ballotName ns.NS,
	closed bool,
) common.AdStrategyTally {

	return ShowLocal(ctx, govAddr, git.Clone(ctx, git.Address(govAddr)).Tree(), ballotName, closed)
}

func ShowLocal(
	ctx context.Context,
	govAddr gov.GovAddress,
	govTree *git.Tree,
	ballotName ns.NS,
	closed bool,
) common.AdStrategyTally {

	ad, strat := load.LoadStrategy(ctx, govTree, ballotName, closed)
	var tally common.Tally
	must.Try(func() { tally = LoadTally(ctx, govTree, ballotName, closed) })
	return common.AdStrategyTally{Ad: ad, Strategy: strat, Tally: tally}
}
