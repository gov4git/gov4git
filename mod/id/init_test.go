package id

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/lib/must"
	"github.com/gov4git/gov4git/mod"
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
	m := PrivateMod{
		NS:      mod.NS(""),
		Public:  publicAddr,
		Private: privateAddr,
	}

	m.Init(ctx)

	if err := must.Try0(func() { m.Init(ctx) }); err == nil {
		t.Fatal("second init must fail")
	}
}
