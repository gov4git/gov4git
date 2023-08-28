package collab

import (
	"context"
	"time"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func IsConcern(ctx context.Context, addr gov.GovAddress, name ConcernName) bool {
	return IsConcern_Local(ctx, gov.Clone(ctx, addr).Tree(), name)
}

func IsConcern_Local(ctx context.Context, t *git.Tree, name ConcernName) bool {
	err := must.Try(func() { concernKV.Get(ctx, concernNS, t, name) })
	return err == nil
}

func OpenConcern(ctx context.Context, addr gov.GovAddress, name ConcernName, title string, desc string, trackerURL string) {
	cloned := gov.Clone(ctx, addr)
	chg := OpenConcern_StageOnly(ctx, cloned.Tree(), name, title, desc, trackerURL)
	proto.Commit(ctx, cloned.Tree(), chg)
	cloned.Push(ctx)
}

func OpenConcern_StageOnly(ctx context.Context, t *git.Tree, name ConcernName, title string, desc string, trackerURL string) git.ChangeNoResult {
	must.Assert(ctx, !IsConcern_Local(ctx, t, name), ErrConcernAlreadyExists)
	state := Concern{
		TimeOpened: time.Now(),
		Name:       name,
		Title:      title,
		Desc:       desc,
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
	state.TimeClosed = time.Now()
	return concernKV.Set(ctx, concernNS, t, name, state)
}

func ListConcerns(ctx context.Context, addr gov.GovAddress) Concerns {
	return ListConcerns_Local(ctx, gov.Clone(ctx, addr).Tree())
}

func ListConcerns_Local(ctx context.Context, t *git.Tree) Concerns {
	_, concerns := concernKV.ListKeyValues(ctx, concernNS, t)
	return concerns
}
