package collab

import (
	"context"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func IsConcern(ctx context.Context, addr gov.GovAddress, name ConcernName) bool {
	return IsConcernLocal(ctx, gov.Clone(ctx, addr).Tree(), name)
}

func IsConcernLocal(ctx context.Context, t *git.Tree, name ConcernName) bool {
	err := must.Try(func() { concernKV.Get(ctx, concernNS, t, name) })
	return err == nil
}

func OpenConcern(ctx context.Context, addr gov.GovAddress, name ConcernName, trackerURL string) {
	cloned := gov.Clone(ctx, addr)
	chg := OpenConcern_StageOnly(ctx, cloned.Tree(), name, trackerURL)
	proto.Commit(ctx, cloned.Tree(), chg)
	cloned.Push(ctx)
}

func OpenConcern_StageOnly(ctx context.Context, t *git.Tree, name ConcernName, trackerURL string) git.ChangeNoResult {
	must.Assert(ctx, !IsConcernLocal(ctx, t, name), ErrConcernAlreadyExists)
	state := ConcernState{
		Name:       name,
		TrackerURL: trackerURL,
		Closed:     false,
	}
	return concernKV.Set(ctx, concernNS, t, name, state)
}

func CloseConcern(ctx context.Context, addr gov.GovAddress, name ConcernName) {
	cloned := gov.Clone(ctx, addr)
	chg := CloseConcern_StageOnly(ctx, cloned.Tree(), name)
	proto.Commit(ctx, cloned.Tree(), chg)
	cloned.Push(ctx)
}

func CloseConcern_StageOnly(ctx context.Context, t *git.Tree, name ConcernName) git.ChangeNoResult {
	state := concernKV.Get(ctx, concernNS, t, name)
	must.Assert(ctx, !state.Closed, ErrConcernAlreadyClosed)
	state.Closed = true
	return concernKV.Set(ctx, concernNS, t, name, state)
}
