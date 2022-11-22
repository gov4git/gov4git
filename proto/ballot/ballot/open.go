package ballot

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/member"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
	"github.com/gov4git/lib4git/ns"
)

func Open(
	ctx context.Context,
	strat common.Strategy,
	govAddr gov.PublicAddress,
	name ns.NS,
	title string,
	description string,
	choices []string,
	participants member.Group,
) git.Change[common.BallotAddress] {

	govRepo, govTree := git.Clone(ctx, git.Address(govAddr))
	chg := OpenStageOnly(ctx, strat, govAddr, govRepo, govTree, name, title, description, choices, participants)
	proto.Commit(ctx, govTree, chg.Msg)
	git.Push(ctx, govRepo)
	return chg
}

func OpenStageOnly(
	ctx context.Context,
	strat common.Strategy,
	govAddr gov.PublicAddress,
	govRepo *git.Repository,
	govTree *git.Tree,
	name ns.NS,
	title string,
	description string,
	choices []string,
	participants member.Group,
) git.Change[common.BallotAddress] {

	// check no open ballots by the same name
	openAdNS := common.OpenBallotNS(name).Sub(common.AdFilebase)
	if _, err := govTree.Filesystem.Stat(openAdNS.Path()); err == nil {
		must.Errorf(ctx, "ballot already exists")
	}

	// check no closed ballots by the same name
	closedAdNS := common.ClosedBallotNS(name).Sub(common.AdFilebase)
	if _, err := govTree.Filesystem.Stat(closedAdNS.Path()); err == nil {
		must.Errorf(ctx, "closed ballot with same name exists")
	}

	// verify group exists
	if !member.IsGroupLocal(ctx, govTree, participants) {
		must.Errorf(ctx, "participant group does not exist")
	}

	// write ad
	ad := common.Advertisement{
		Community:    govAddr,
		Name:         name,
		Title:        title,
		Description:  description,
		Choices:      choices,
		Strategy:     strat.Name(),
		Participants: participants,
		ParentCommit: git.Head(ctx, govRepo),
	}
	git.ToFileStage(ctx, govTree, openAdNS.Path(), ad)

	// write strategy
	openStratNS := common.OpenBallotNS(name).Sub(common.StrategyFilebase)
	git.ToFileStage(ctx, govTree, openStratNS.Path(), strat)

	return git.Change[common.BallotAddress]{
		Result: common.BallotAddress{Gov: govAddr, Name: name},
		Msg:    fmt.Sprintf("Create ballot of type %v", strat.Name()),
	}
}
