package core

import (
	"context"

	"github.com/gov4git/gov4git/mod/ballot/load"
	"github.com/gov4git/gov4git/mod/ballot/proto"
	"github.com/gov4git/gov4git/mod/gov"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/ns"
)

func Show(
	ctx context.Context,
	govAddr gov.CommunityAddress,
	ballotName ns.NS,
) proto.AdStrategyTally {

	govRepo := git.CloneRepo(ctx, git.Address(govAddr))
	chg := ShowLocal(ctx, govAddr, git.Worktree(ctx, govRepo), ballotName)
	return chg
}

func ShowLocal(
	ctx context.Context,
	govAddr gov.CommunityAddress,
	govTree *git.Tree,
	ballotName ns.NS,
) proto.AdStrategyTally {

	ad, strat := load.LoadStrategy(ctx, govTree, ballotName)
	tally := LoadTally(ctx, govTree, ballotName)
	return proto.AdStrategyTally{Ad: ad, Strategy: strat, Tally: tally}
}
