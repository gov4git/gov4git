package collab

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/util"
)

func ProposalAddressesConcern(ctx context.Context, addr gov.GovAddress, proposal Name, concern Name) git.Change[form.Map, ProposalConcernPair] {
	cloned := gov.Clone(ctx, addr)
	chg := ProposalAddressesConcern_StageOnly(ctx, cloned.Tree(), proposal, concern)
	proto.Commit(ctx, cloned.Tree(), chg)
	cloned.Push(ctx)
	return chg
}

func ProposalAddressesConcern_StageOnly(ctx context.Context, t *git.Tree, proposalName Name, concernName Name) git.Change[form.Map, ProposalConcernPair] {

	// read state
	proposal := proposalKV.Get(ctx, proposalNS, t, proposalName)
	concern := concernKV.Get(ctx, concernNS, t, concernName)

	// update state
	if !util.IsIn(concernName, proposal.Addresses...) {
		proposal.Addresses = append(proposal.Addresses, concernName)
	}
	if !util.IsIn(proposalName, concern.AddressedBy...) {
		concern.AddressedBy = append(concern.AddressedBy, proposalName)
	}

	// write state
	concernKV.Set(ctx, concernNS, t, concernName, concern)
	proposalKV.Set(ctx, proposalNS, t, proposalName, proposal)

	return git.NewChange(
		fmt.Sprintf("Proposal %v addresses concern %v", proposalName, concernName),
		"collab_proposal_addresses_concern",
		form.Map{"proposal": proposal, "concern": concern},
		ProposalConcernPair{Proposal: proposal, Concern: concern},
		nil,
	)
}
