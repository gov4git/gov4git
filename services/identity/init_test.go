package identity

import (
	"context"
	"testing"

	"github.com/petar/gov4git/lib/base"
	"github.com/petar/gov4git/lib/files"
	"github.com/petar/gov4git/lib/git"
	"github.com/petar/gov4git/proto"
)

func TestSoulInit(t *testing.T) {
	ctx := context.Background()

	base.Infof("using temp dir %v", t.TempDir())
	testPrivDir := files.PathDir(t.TempDir()).Subdir("private")
	testPubDir := files.PathDir(t.TempDir()).Subdir("public")
	testPriv := git.LocalFromDir(testPrivDir)
	testPub := git.LocalFromDir(testPubDir)

	// make bare test public and private repos
	if err := testPriv.InitBare(ctx); err != nil {
		t.Fatal(err)
	}
	if err := testPub.InitBare(ctx); err != nil {
		t.Fatal(err)
	}

	// init soul
	apiCtx := files.WithWorkDir(ctx, files.PathDir(t.TempDir()).Subdir("soul_api"))
	api := IdentityService{IdentityConfig: proto.IdentityConfig{PublicURL: testPubDir.Path, PrivateURL: testPrivDir.Path}}
	if _, err := api.Init(apiCtx, &IdentityInitIn{}); err != nil {
		t.Fatal(err)
	}

	// re-init should return error
	if _, err := api.Init(apiCtx, &IdentityInitIn{}); err == nil {
		t.Fatal("re-initializing a soul should fail")
	}
}
