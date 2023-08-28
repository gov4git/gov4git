package collab

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/ballot/ballot"
	"github.com/gov4git/gov4git/proto/ballot/qv"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/member"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
)

func FixConcernPriority(
	ctx context.Context,
	addr gov.GovAddress,
	concern Name,
	priority float64,
) git.Change[form.Map, Concern] {
	cloned := gov.Clone(ctx, addr)
	chg := FixConcernPriority_StageOnly(ctx, addr, cloned, concern, priority)
	proto.Commit(ctx, cloned.Tree(), chg)
	cloned.Push(ctx)
	return chg
}

func FixConcernPriority_StageOnly(
	ctx context.Context,
	govAddr gov.GovAddress,
	govCloned git.Cloned,
	concernName Name,
	priority float64,
) git.Change[form.Map, Concern] {

	concern := concernKV.Get(ctx, concernNS, govCloned.Tree(), concernName)
	concern.Priority.Fixed = form.Float64(priority)
	concernKV.Set(ctx, concernNS, govCloned.Tree(), concernName, concern)

	return git.NewChange(
		fmt.Sprintf("Fix concern %v priority to %v", concernName, priority),
		"collab_fix_concern_priority",
		form.Map{"concern": concern, "priority": priority},
		concern,
		nil,
	)
}

func PrioritizeConcernByBallot(
	ctx context.Context,
	addr gov.GovAddress,
	concern Name,
) git.Change[form.Map, Concern] {
	cloned := gov.Clone(ctx, addr)
	chg := PrioritizeConcernByBallot_StageOnly(ctx, addr, cloned, concern)
	proto.Commit(ctx, cloned.Tree(), chg)
	cloned.Push(ctx)
	return chg
}

func PrioritizeConcernByBallot_StageOnly(
	ctx context.Context,
	govAddr gov.GovAddress,
	govCloned git.Cloned,
	concernName Name,
) git.Change[form.Map, Concern] {

	concern := concernKV.Get(ctx, concernNS, govCloned.Tree(), concernName)

	// XXX: what if there is a prior priority set

	ballotName := ConcernPriorityBallotName(concernName)
	chg := ballot.Open_StageOnly(
		ctx,
		qv.QV{},
		govAddr,
		govCloned,
		ballotName,
		fmt.Sprintf("Priority poll for concern %v", concernName),            // title
		fmt.Sprintf("Up/down vote the priority of concern %v", concernName), // description
		[]string{PriorityBallotChoice},                                      // choices
		member.Everybody,
	)
	concern.Priority.Ballot = &ballotName

	concernKV.Set(ctx, concernNS, govCloned.Tree(), concernName, concern)

	return git.NewChange(
		fmt.Sprintf("Prioritize concern %v by ballot %v", concernName, ballotName),
		"collab_prioritize_concern_by_ballot",
		form.Map{"concern": concern},
		concern,
		form.Forms{chg},
	)
}
