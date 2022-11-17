package core

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/mod"
	"github.com/gov4git/gov4git/mod/ballot/load"
	"github.com/gov4git/gov4git/mod/ballot/proto"
	"github.com/gov4git/gov4git/mod/gov"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/ns"
)

func Close(
	ctx context.Context,
	govAddr gov.CommunityAddress,
	ballotName ns.NS,
) git.ChangeNoResult {

	govRepo, govTree := git.Clone(ctx, git.Address(govAddr))
	chg := CloseStageOnly(ctx, govAddr, govRepo, govTree, ballotName)
	mod.Commit(ctx, govTree, chg.Msg)
	git.Push(ctx, govRepo)
	return chg
}

func CloseStageOnly(
	ctx context.Context,
	govAddr gov.CommunityAddress,
	govRepo *git.Repository,
	govTree *git.Tree,
	ballotName ns.NS,
) git.ChangeNoResult {

	// verify ad and strategy are present
	load.LoadStrategy(ctx, govTree, ballotName)

	openNS := proto.OpenBallotNS(ballotName)
	closedNS := proto.ClosedBallotNS(ballotName)
	git.RenameStage(ctx, govTree, openNS.Path(), closedNS.Path())

	return git.ChangeNoResult{
		Msg: fmt.Sprintf("closed ballot %v", ballotName),
	}
}
