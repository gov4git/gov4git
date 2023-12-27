package ballotapi

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/v2/proto"
	"github.com/gov4git/gov4git/v2/proto/account"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotio"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/history/trace"
	"github.com/gov4git/gov4git/v2/proto/member"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
	"github.com/gov4git/gov4git/v2/proto/purpose"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func Open(
	ctx context.Context,
	strat ballotproto.StrategyName,
	govAddr gov.OwnerAddress,
	name ballotproto.BallotName,
	owner account.AccountID,
	purpose purpose.Purpose,
	motionPolicy motionproto.PolicyName,
	title string,
	description string,
	choices []string,
	participants member.Group,

) git.Change[form.Map, ballotproto.BallotAddress] {

	govCloned := gov.CloneOwner(ctx, govAddr)
	chg := Open_StageOnly(ctx, strat, govCloned, name, owner, purpose, motionPolicy, title, description, choices, participants)
	proto.Commit(ctx, govCloned.Public.Tree(), chg)
	govCloned.Public.Push(ctx)
	return chg
}

func Open_StageOnly(
	ctx context.Context,
	strategyName ballotproto.StrategyName,
	cloned gov.OwnerCloned,
	name ballotproto.BallotName,
	owner account.AccountID,
	purpose purpose.Purpose,
	motionPolicy motionproto.PolicyName,
	title string,
	description string,
	choices []string,
	participants member.Group,

) git.Change[form.Map, ballotproto.BallotAddress] {

	// check no open ballots by the same name
	if _, err := git.TreeStat(ctx, cloned.Public.Tree(), name.AdNS()); err == nil {
		must.Errorf(ctx, "ballot already exists: %v", name.AdNS().GitPath())
	}

	// verify group exists
	if !member.IsGroup_Local(ctx, cloned.PublicClone(), participants) {
		must.Errorf(ctx, "participant group %v does not exist", participants)
	}

	// create escrow account
	account.Create_StageOnly(
		ctx, cloned.PublicClone(),
		ballotproto.BallotEscrowAccountID(name),
		account.NobodyAccountID,
		fmt.Sprintf("opening ballot %v", name),
	)

	// write ad
	ad := ballotproto.Advertisement{
		Gov:          cloned.GovAddress(),
		Name:         name,
		Owner:        owner,
		Purpose:      purpose,
		MotionPolicy: motionPolicy,
		//
		Title:       title,
		Description: description,
		//
		Choices:      choices,
		Strategy:     strategyName,
		Participants: participants,
		//
		Frozen:    false,
		Closed:    false,
		Cancelled: false,
		//
		ParentCommit: git.Head(ctx, cloned.Public.Repo()),
	}
	git.ToFileStage(ctx, cloned.Public.Tree(), name.AdNS(), ad)

	// initialize tally
	strategy := ballotio.LookupStrategy(ctx, strategyName)
	tally := strategy.Open(ctx, cloned, &ad)
	git.ToFileStage(ctx, cloned.Public.Tree(), name.TallyNS(), tally)

	// log
	trace.Log_StageOnly(ctx, cloned.PublicClone(), &trace.Event{
		Op:     "ballot_open",
		Args:   trace.M{"name": name},
		Result: trace.M{"ad": ad},
	})

	return git.NewChange(
		fmt.Sprintf("Create ballot of type %v", strategyName),
		"ballot_open",
		form.Map{
			"strategy":     strategyName,
			"name":         name,
			"participants": participants,
		},
		ballotproto.BallotAddress{Gov: cloned.GovAddress(), Name: name},
		nil,
	)
}
