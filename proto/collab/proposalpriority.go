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

func FixProposalPriority(
	ctx context.Context,
	addr gov.GovAddress,
	proposal Name,
	priority float64,
) git.Change[form.Map, Proposal] {
	cloned := gov.Clone(ctx, addr)
	chg := FixProposalPriority_StageOnly(ctx, addr, cloned, proposal, priority)
	proto.Commit(ctx, cloned.Tree(), chg)
	cloned.Push(ctx)
	return chg
}

func FixProposalPriority_StageOnly(
	ctx context.Context,
	govAddr gov.GovAddress,
	govCloned git.Cloned,
	proposalName Name,
	priority float64,
) git.Change[form.Map, Proposal] {

	proposal := proposalKV.Get(ctx, proposalNS, govCloned.Tree(), proposalName)
	proposal.Priority.Fixed = form.Float64(priority)
	proposalKV.Set(ctx, proposalNS, govCloned.Tree(), proposalName, proposal)

	return git.NewChange(
		fmt.Sprintf("Fix proposal %v priority to %v", proposalName, priority),
		"collab_fix_proposal_priority",
		form.Map{"proposal": proposal, "priority": priority},
		proposal,
		nil,
	)
}

func PrioritizeProposalByBallot(
	ctx context.Context,
	addr gov.GovAddress,
	proposal Name,
) git.Change[form.Map, Proposal] {
	cloned := gov.Clone(ctx, addr)
	chg := PrioritizeProposalByBallot_StageOnly(ctx, addr, cloned, proposal)
	proto.Commit(ctx, cloned.Tree(), chg)
	cloned.Push(ctx)
	return chg
}

func PrioritizeProposalByBallot_StageOnly(
	ctx context.Context,
	govAddr gov.GovAddress,
	govCloned git.Cloned,
	proposalName Name,
) git.Change[form.Map, Proposal] {

	proposal := proposalKV.Get(ctx, proposalNS, govCloned.Tree(), proposalName)

	// XXX: what if there is a prior priority set

	ballotName := ProposalPriorityBallotName(proposalName)
	chg := ballot.Open_StageOnly(
		ctx,
		qv.QV{},
		govAddr,
		govCloned,
		ballotName,
		fmt.Sprintf("Priority poll for proposal %v", proposalName),            // title
		fmt.Sprintf("Up/down vote the priority of proposal %v", proposalName), // description
		[]string{PriorityBallotChoice},                                        // choices
		member.Everybody,
	)
	proposal.Priority.Ballot = &ballotName

	proposalKV.Set(ctx, proposalNS, govCloned.Tree(), proposalName, proposal)

	return git.NewChange(
		fmt.Sprintf("Prioritize proposal %v by ballot %v", proposalName, ballotName),
		"collab_prioritize_proposal_by_ballot",
		form.Map{"proposal": proposal},
		proposal,
		form.Forms{chg},
	)
}
