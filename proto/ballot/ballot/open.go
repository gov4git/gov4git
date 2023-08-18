package ballot

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/member"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
	"github.com/gov4git/lib4git/ns"
)

func Open(
	ctx context.Context,
	strat common.Strategy,
	govAddr gov.GovAddress,
	name ns.NS,
	title string,
	description string,
	choices []string,
	participants member.Group,
) git.Change[form.Map, common.BallotAddress] {

	govCloned := git.CloneOne(ctx, git.Address(govAddr))
	chg := OpenStageOnly(ctx, strat, govAddr, govCloned, name, title, description, choices, participants)
	proto.Commit(ctx, govCloned.Tree(), chg)
	govCloned.Push(ctx)
	return chg
}

func OpenStageOnly(
	ctx context.Context,
	strat common.Strategy,
	govAddr gov.GovAddress,
	govCloned git.Cloned,
	name ns.NS,
	title string,
	description string,
	choices []string,
	participants member.Group,
) git.Change[form.Map, common.BallotAddress] {

	// check no open ballots by the same name
	openAdNS := common.BallotPath(name).Sub(common.AdFilebase)
	if _, err := govCloned.Tree().Filesystem.Stat(openAdNS.Path()); err == nil {
		must.Errorf(ctx, "ballot already exists")
	}

	// verify group exists
	if !member.IsGroupLocal(ctx, govCloned.Tree(), participants) {
		must.Errorf(ctx, "participant group does not exist")
	}

	// write ad
	ad := common.Advertisement{
		Gov:          govAddr,
		Name:         name,
		Title:        title,
		Description:  description,
		Choices:      choices,
		Strategy:     strat.Name(),
		Participants: participants,
		Frozen:       false,
		Closed:       false,
		Cancelled:    false,
		ParentCommit: git.Head(ctx, govCloned.Repo()),
	}
	git.ToFileStage(ctx, govCloned.Tree(), openAdNS.Path(), ad)

	// write initial tally
	tally := common.Tally{
		Ad:            ad,
		Scores:        map[string]float64{},
		VotesByUser:   map[member.User]map[string]common.StrengthAndScore{},
		AcceptedVotes: map[member.User]common.AcceptedElections{},
		RejectedVotes: map[member.User]common.RejectedElections{},
		Charges:       map[member.User]float64{},
	}
	openTallyNS := common.BallotPath(name).Sub(common.TallyFilebase)
	git.ToFileStage(ctx, govCloned.Tree(), openTallyNS.Path(), tally)

	// write strategy
	openStratNS := common.BallotPath(name).Sub(common.StrategyFilebase)
	git.ToFileStage(ctx, govCloned.Tree(), openStratNS.Path(), strat)

	return git.NewChange(
		fmt.Sprintf("Create ballot of type %v", strat.Name()),
		"ballot_open",
		form.Map{
			"strategy":     strat,
			"name":         name,
			"participants": participants,
		},
		common.BallotAddress{Gov: govAddr, Name: name},
		nil,
	)
}
