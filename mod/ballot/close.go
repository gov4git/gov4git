package ballot

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/mod"
	"github.com/gov4git/gov4git/mod/gov"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/ns"
)

func Close[S Strategy](
	ctx context.Context,
	govAddr gov.CommunityAddress,
	ballotName ns.NS,
) git.ChangeNoResult {

	govRepo, govTree := git.Clone(ctx, git.Address(govAddr))
	chg := CloseStageOnly[S](ctx, govAddr, govRepo, govTree, ballotName)
	mod.Commit(ctx, govTree, chg.Msg)
	git.Push(ctx, govRepo)
	return chg
}

func CloseStageOnly[S Strategy](
	ctx context.Context,
	govAddr gov.CommunityAddress,
	govRepo *git.Repository,
	govTree *git.Tree,
	ballotName ns.NS,
) git.ChangeNoResult {

	openNS := OpenBallotNS[S](ballotName)
	closedNS := ClosedBallotNS[S](ballotName)

	// verify ad is present
	git.FromFile[Advertisement](ctx, govTree, openNS.Sub(adFilebase).Path())

	git.RenameStage(ctx, govTree, openNS.Path(), closedNS.Path())

	return git.ChangeNoResult{
		Msg: fmt.Sprintf("closed ballot %v", ballotName),
	}
}
