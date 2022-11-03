package id

import (
	"context"
	"testing"

	"github.com/gov4git/gov4git/lib/base"
	"github.com/gov4git/gov4git/lib/files"
	"github.com/gov4git/gov4git/lib/git"
)

func TestInit(t *testing.T) {
	ctx := context.Background()

	base.Infof("using temp dir %v", t.TempDir())
	testPrivDir := files.PathDir(t.TempDir()).Subdir("private")
	testPubDir := files.PathDir(t.TempDir()).Subdir("public")
	testPriv := git.LocalInDir(testPrivDir)
	testPub := git.LocalInDir(testPubDir)

	// make bare test public and private repos
	if err := testPriv.InitBare(ctx); err != nil {
		t.Fatal(err)
	}
	if err := testPub.InitBare(ctx); err != nil {
		t.Fatal(err)
	}

	// init soul
	apiCtx := files.WithWorkDir(ctx, files.PathDir(t.TempDir()).Subdir("soul_api"))
	publicOrigin := git.Origin{
		Repo:   git.URL(testPubDir.Path),
		Branch: git.MainBranch,
	}
	privateOrigin := git.Origin{
		Repo:   git.URL(testPrivDir.Path),
		Branch: git.MainBranch,
	}
	api := IdentityPrivateService{
		PublicAddress:  publicOrigin,
		PrivateAddress: privateOrigin,
	}
	if _, err := api.Init(apiCtx); err != nil {
		t.Fatal(err)
	}

	// re-init should return error
	if _, err := api.Init(apiCtx); err == nil {
		t.Fatal("re-initializing a soul should fail")
	}
}
