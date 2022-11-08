package id

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/lib/must"
)

func TestInit(t *testing.T) {
	ctx := context.Background()
	publicDir := filepath.Join(t.TempDir(), "public")
	privateDir := filepath.Join(t.TempDir(), "private")
	fmt.Printf("public_dir=%v private_dir=%v\n", publicDir, privateDir)
	git.InitPlain(ctx, publicDir, true)
	git.InitPlain(ctx, privateDir, true)

	publicAddr := git.NewAddress(git.URL(publicDir), git.MainBranch)
	privateAddr := git.NewAddress(git.URL(privateDir), git.MainBranch)
	Init(ctx, publicAddr, privateAddr)

	if err := must.Try(func() { Init(ctx, publicAddr, privateAddr) }); err == nil {
		t.Fatal("second init must fail")
	}
}
