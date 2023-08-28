package collab

import (
	"context"
	"time"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func IsProposal(ctx context.Context, addr gov.GovAddress, name ProposalName) bool {
	return IsProposal_Local(ctx, gov.Clone(ctx, addr).Tree(), name)
}

func IsProposal_Local(ctx context.Context, t *git.Tree, name ProposalName) bool {
	err := must.Try(func() { proposalKV.Get(ctx, proposalNS, t, name) })
	return err == nil
}

func OpenProposal(ctx context.Context, addr gov.GovAddress, name ProposalName, title string, desc string, trackerURL string) {
	cloned := gov.Clone(ctx, addr)
	chg := OpenProposal_StageOnly(ctx, cloned.Tree(), name, title, desc, trackerURL)
	proto.Commit(ctx, cloned.Tree(), chg)
	cloned.Push(ctx)
}

func OpenProposal_StageOnly(ctx context.Context, t *git.Tree, name ProposalName, title string, desc string, trackerURL string) git.ChangeNoResult {
	must.Assert(ctx, !IsProposal_Local(ctx, t, name), ErrProposalAlreadyExists)
	state := Proposal{
		TimeOpened: time.Now(),
		Name:       name,
		Desc:       desc,
		Title:      title,
		TrackerURL: trackerURL,
		Closed:     false,
	}
	return proposalKV.Set(ctx, proposalNS, t, name, state)
}

func CloseProposal(ctx context.Context, addr gov.GovAddress, name ProposalName) {
	cloned := gov.Clone(ctx, addr)
	chg := CloseProposal_StageOnly(ctx, cloned.Tree(), name)
	proto.Commit(ctx, cloned.Tree(), chg)
	cloned.Push(ctx)
}

func CloseProposal_StageOnly(ctx context.Context, t *git.Tree, name ProposalName) git.ChangeNoResult {
	state := proposalKV.Get(ctx, proposalNS, t, name)
	must.Assert(ctx, !state.Closed, ErrProposalAlreadyClosed)
	state.Closed = true
	state.TimeClosed = time.Now()
	return proposalKV.Set(ctx, proposalNS, t, name, state)
}

func ListProposals(ctx context.Context, addr gov.GovAddress) Proposals {
	return ListProposals_Local(ctx, gov.Clone(ctx, addr).Tree())
}

func ListProposals_Local(ctx context.Context, t *git.Tree) Proposals {
	_, proposals := proposalKV.ListKeyValues(ctx, proposalNS, t)
	return proposals
}
