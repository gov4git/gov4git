package testutil

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/gov4git/gov4git/lib/base"
	"github.com/gov4git/gov4git/lib/git"
)

type PlainRepo struct {
	Dir  string
	Repo *git.Repository
}

func InitPlainRepo(t *testing.T, ctx context.Context) PlainRepo {
	repoDir := t.TempDir()
	base.Infof("repo %v", repoDir)
	repo := git.InitPlain(ctx, repoDir, false) // not bare
	return PlainRepo{Dir: repoDir, Repo: repo}
}

func Hang() {
	<-(chan int)(nil)
}

type LocalAddress struct {
	Dir     string
	Repo    *git.Repository
	Tree    *git.Tree
	Address git.Address
}

func (x LocalAddress) String() string {
	return fmt.Sprintf("test address, dir=%v\n", x.Dir)
}

func NewLocalAddress(ctx context.Context, t *testing.T, branch git.Branch, isBare bool) LocalAddress {
	dir := filepath.Join(t.TempDir(), UniqueString(ctx))
	repo := git.InitPlain(ctx, dir, isBare)
	addr := git.NewAddress(git.URL(dir), branch)
	var tree *git.Tree
	if !isBare {
		tree = git.Worktree(ctx, repo)
	}
	return LocalAddress{Dir: dir, Repo: repo, Address: addr, Tree: tree}
}
