package testutil

import (
	"context"
	"testing"

	"github.com/gov4git/gov4git/lib/base"
	"github.com/gov4git/gov4git/lib/git"
)

func InitRepo(t *testing.T, ctx context.Context) *git.Repository {
	repoDir := t.TempDir()
	base.Infof("repo %v", repoDir)
	return git.InitPlain(ctx, repoDir, false)
}

func Hang() {
	<-(chan int)(nil)
}
