package ballot

import (
	"context"

	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/ballot/load"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/ns"
)

func Show(
	ctx context.Context,
	govAddr gov.CommunityAddress,
	ballotName ns.NS,
) common.AdStrategyTally {

	govRepo := git.CloneRepo(ctx, git.Address(govAddr))
	chg := ShowLocal(ctx, govAddr, git.Worktree(ctx, govRepo), ballotName)
	return chg
}

func ShowLocal(
	ctx context.Context,
	govAddr gov.CommunityAddress,
	govTree *git.Tree,
	ballotName ns.NS,
) common.AdStrategyTally {

	ad, strat := load.LoadStrategy(ctx, govTree, ballotName)
	tally := LoadTally(ctx, govTree, ballotName)
	return common.AdStrategyTally{Ad: ad, Strategy: strat, Tally: tally}
}
