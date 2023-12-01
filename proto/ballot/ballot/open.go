package ballot

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/account"
	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/history"
	"github.com/gov4git/gov4git/proto/member"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func Open(
	ctx context.Context,
	strat common.Strategy,
	govAddr gov.OwnerAddress,
	name common.BallotName,
	title string,
	description string,
	choices []string,
	participants member.Group,
) git.Change[form.Map, common.BallotAddress] {

	govCloned := gov.CloneOwner(ctx, govAddr)
	chg := Open_StageOnly(ctx, strat, govCloned, name, title, description, choices, participants)
	proto.Commit(ctx, govCloned.Public.Tree(), chg)
	govCloned.Public.Push(ctx)
	return chg
}

func Open_StageOnly(
	ctx context.Context,
	strat common.Strategy,
	cloned gov.OwnerCloned,
	name common.BallotName,
	title string,
	description string,
	choices []string,
	participants member.Group,
) git.Change[form.Map, common.BallotAddress] {

	// check no open ballots by the same name
	openAdNS := common.BallotPath(name).Append(common.AdFilebase)
	if _, err := git.TreeStat(ctx, cloned.Public.Tree(), openAdNS); err == nil {
		must.Errorf(ctx, "ballot already exists: %v", openAdNS.GitPath())
	}

	// verify group exists
	if !member.IsGroup_Local(ctx, cloned.PublicClone(), participants) {
		must.Errorf(ctx, "participant group %v does not exist", participants)
	}

	// create escrow account
	account.Create_StageOnly(
		ctx, cloned.PublicClone(),
		common.BallotEscrowAccountID(name),
		common.BallotOwnerID(name),
	)

	// write ad
	ad := common.Advertisement{
		Gov:          cloned.GovAddress(),
		Name:         name,
		Title:        title,
		Description:  description,
		Choices:      choices,
		Strategy:     strat.Name(),
		Participants: participants,
		Frozen:       false,
		Closed:       false,
		Cancelled:    false,
		ParentCommit: git.Head(ctx, cloned.Public.Repo()),
	}
	git.ToFileStage(ctx, cloned.Public.Tree(), openAdNS, ad)

	// write initial tally
	tally := common.Tally{
		Ad:            ad,
		Scores:        map[string]float64{},
		VotesByUser:   map[member.User]map[string]common.StrengthAndScore{},
		AcceptedVotes: map[member.User]common.AcceptedElections{},
		RejectedVotes: map[member.User]common.RejectedElections{},
		Charges:       map[member.User]float64{},
	}
	openTallyNS := common.BallotPath(name).Append(common.TallyFilebase)
	git.ToFileStage(ctx, cloned.Public.Tree(), openTallyNS, tally)

	// write strategy
	openStratNS := common.BallotPath(name).Append(common.StrategyFilebase)
	git.ToFileStage(ctx, cloned.Public.Tree(), openStratNS, strat)

	// log
	history.Log_StageOnly(ctx, cloned.PublicClone(), &history.Event{
		Op: &history.Op{
			Op:     "ballot_open",
			Args:   history.M{"name": name},
			Result: history.M{"ad": ad},
		},
	})

	return git.NewChange(
		fmt.Sprintf("Create ballot of type %v", strat.Name()),
		"ballot_open",
		form.Map{
			"strategy":     strat,
			"name":         name,
			"participants": participants,
		},
		common.BallotAddress{Gov: cloned.GovAddress(), Name: name},
		nil,
	)
}
