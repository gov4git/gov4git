package collab

import (
	"context"

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

func OpenProposal(ctx context.Context, addr gov.GovAddress, name ProposalName, trackerURL string) {
	cloned := gov.Clone(ctx, addr)
	chg := OpenProposal_StageOnly(ctx, cloned.Tree(), name, trackerURL)
	proto.Commit(ctx, cloned.Tree(), chg)
	cloned.Push(ctx)
}

func OpenProposal_StageOnly(ctx context.Context, t *git.Tree, name ProposalName, trackerURL string) git.ChangeNoResult {
	must.Assert(ctx, !IsProposal_Local(ctx, t, name), ErrProposalAlreadyExists)
	state := ProposalState{
		Name:       name,
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
	return proposalKV.Set(ctx, proposalNS, t, name, state)
}
