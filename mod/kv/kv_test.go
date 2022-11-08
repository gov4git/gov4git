package kv

import (
	"context"
	"testing"

	"github.com/gov4git/gov4git/lib/base"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/lib/ns"
	"github.com/gov4git/gov4git/lib/testutil"
)

func TestSetGet(t *testing.T) {
	base.LogVerbosely()
	ctx := context.Background()
	repo := testutil.InitPlain(t, ctx)

	m := ns.NS("ns")
	wt := git.Worktree(ctx, repo.Repo)

	x := KV[string, float64]{}
	key := "a"
	value := float64(3.14)
	x.Set(ctx, m, wt, key, value)

	got := x.Get(ctx, m, wt, key)
	if got != value {
		t.Errorf("expecting %v, got %v", value, got)
	}

	keys := x.ListKeys(ctx, m, wt)
	if len(keys) != 1 || keys[0] != key {
		t.Errorf("listing keys")
	}

	x.Remove(ctx, m, wt, key)
	keys = x.ListKeys(ctx, m, wt)
	if len(keys) != 0 {
		t.Errorf("list after remove")
	}

	// <-(chan int)(nil)
}
