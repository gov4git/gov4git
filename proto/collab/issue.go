package collab

import (
	"context"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func IsIssue(ctx context.Context, addr gov.GovAddress, name IssueName) bool {
	return IsIssueLocal(ctx, gov.Clone(ctx, addr).Tree(), name)
}

func IsIssueLocal(ctx context.Context, t *git.Tree, name IssueName) bool {
	err := must.Try(func() { issuesKV.Get(ctx, issuesNS, t, name) })
	return err == nil
}

func OpenIssue(ctx context.Context, addr gov.GovAddress, name IssueName, trackerURL string) {
	cloned := gov.Clone(ctx, addr)
	chg := OpenIssueStageOnly(ctx, cloned.Tree(), name, trackerURL)
	proto.Commit(ctx, cloned.Tree(), chg.Msg)
	cloned.Push(ctx)
}

func OpenIssueStageOnly(ctx context.Context, t *git.Tree, name IssueName, trackerURL string) git.ChangeNoResult {
	if IsIssueLocal(ctx, t, name) {
		must.Panic(ctx, ErrIssueAlreadyExists)
	}
	state := IssueState{
		Name:       name,
		TrackerURL: trackerURL,
		Closed:     false,
	}
	return issuesKV.Set(ctx, issuesNS, t, name, state)
}

func CloseIssue(ctx context.Context, addr gov.GovAddress, name IssueName) {
	cloned := gov.Clone(ctx, addr)
	chg := CloseIssueStageOnly(ctx, cloned.Tree(), name)
	proto.Commit(ctx, cloned.Tree(), chg.Msg)
	cloned.Push(ctx)
}

func CloseIssueStageOnly(ctx context.Context, t *git.Tree, name IssueName) git.ChangeNoResult {
	state := issuesKV.Get(ctx, issuesNS, t, name)
	if state.Closed {
		must.Panic(ctx, ErrIssueAlreadyClosed)
	}
	state.Closed = true
	return issuesKV.Set(ctx, issuesNS, t, name, state)
}
