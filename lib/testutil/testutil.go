package testutil

import (
	"context"
	"testing"

	"github.com/gov4git/gov4git/lib/base"
	"github.com/gov4git/gov4git/lib/git"
)

type TestPlainRepo struct {
	Dir  string
	Repo *git.Repository
}

func InitPlain(t *testing.T, ctx context.Context) TestPlainRepo {
	repoDir := t.TempDir()
	base.Infof("repo %v", repoDir)
	repo := git.InitPlain(ctx, repoDir, false)
	return TestPlainRepo{Dir: repoDir, Repo: repo}
}

func Hang() {
	<-(chan int)(nil)
}
