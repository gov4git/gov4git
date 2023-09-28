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

func Show(ctx context.Context, govAddr gov.GovAddress, ballotName ns.NS) common.AdStrategyTally {

	return Show_Local(ctx, govAddr, git.CloneOne(ctx, git.Address(govAddr)).Tree(), ballotName)
}

func Show_Local(
	ctx context.Context,
	govAddr gov.GovAddress,
	govTree *git.Tree,
	ballotName ns.NS,
) common.AdStrategyTally {

	ad, strat := load.LoadStrategy(ctx, govTree, ballotName)
	var tally common.Tally
	must.Try(func() { tally = LoadTally(ctx, govTree, ballotName) })
	return common.AdStrategyTally{Ad: ad, Strategy: strat, Tally: tally}
}
