package kv

import (
	"context"
	"testing"

	"github.com/gov4git/gov4git/lib/base"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/mod"
)

func testInitRepo(t *testing.T, ctx context.Context) *git.Repository {
	repoDir := t.TempDir()
	base.Infof("repo %v", repoDir)
	return git.InitPlain(ctx, repoDir, false)
}

func TestSetGet(t *testing.T) {
	base.LogVerbosely()
	ctx := context.Background()
	repo := testInitRepo(t, ctx)

	m := mod.NS("ns")
	wt := git.Worktree(ctx, repo)

	key := git.URL("a")
	value := float64(3.14)
	Set(ctx, m, wt, key, value)

	got := Get[float64](ctx, m, wt, key)
	if got != value {
		t.Errorf("expecting %v, got %v", value, got)
	}

	keys := ListKeys(ctx, m, wt)
	if len(keys) != 1 || keys[0] != key {
		t.Errorf("listing keys")
	}

	Remove(ctx, m, wt, key)
	keys = ListKeys(ctx, m, wt)
	if len(keys) != 0 {
		t.Errorf("list after remove")
	}

	// <-(chan int)(nil)
}
