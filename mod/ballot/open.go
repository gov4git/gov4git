package ballot

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/lib/must"
	"github.com/gov4git/gov4git/lib/ns"
	"github.com/gov4git/gov4git/mod/gov"
	"github.com/gov4git/gov4git/mod/member"
)

func Open[S Strategy](
	ctx context.Context,
	govAddr gov.CommunityAddress,
	name ns.NS,
	title string,
	description string,
	choices []string,
	participants member.Group,
) git.Change[BallotAddress[S]] {

	govRepo, govTree := git.CloneBranchTree(ctx, git.Address(govAddr))
	chg := OpenStageOnly[S](ctx, govAddr, govRepo, govTree, name, title, description, choices, participants)
	git.Commit(ctx, govTree, chg.Msg)
	git.Push(ctx, govRepo)
	return chg
}

func OpenStageOnly[S Strategy](
	ctx context.Context,
	govAddr gov.CommunityAddress,
	govRepo *git.Repository,
	govTree *git.Tree,
	name ns.NS,
	title string,
	description string,
	choices []string,
	participants member.Group,
) git.Change[BallotAddress[S]] {

	// check no open ballots by the same name
	openAdNS := OpenBallotNS[S](name).Sub(adFilebase)
	if _, err := govTree.Filesystem.Stat(openAdNS.Path()); err == nil {
		must.Errorf(ctx, "ballot already exists")
	}

	// check no closed ballots by the same name
	closedAdNS := ClosedBallotNS[S](name).Sub(adFilebase)
	if _, err := govTree.Filesystem.Stat(closedAdNS.Path()); err == nil {
		must.Errorf(ctx, "closed ballot with same name exists")
	}

	// verify group exists
	if !member.IsGroup(ctx, govTree, participants) {
		must.Errorf(ctx, "participant group does not exist")
	}

	// write ad
	var s S
	ad := AdForm{
		Community:    govAddr,
		Name:         name,
		Title:        title,
		Description:  description,
		Choices:      choices,
		Strategy:     s.StrategyName(),
		Participants: participants,
		ParentCommit: git.Head(ctx, govRepo),
	}
	git.ToFileStage(ctx, govTree, openAdNS.Path(), ad)

	return git.Change[BallotAddress[S]]{
		Result: BallotAddress[S]{Gov: govAddr, Name: name},
		Msg:    fmt.Sprintf("Create ballot of type %v", s.StrategyName()),
	}
}
