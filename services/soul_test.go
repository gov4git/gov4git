package services

import (
	"context"
	"testing"

	"github.com/petar/gitty/lib/base"
	"github.com/petar/gitty/lib/files"
	"github.com/petar/gitty/lib/git"
	"github.com/petar/gitty/proto"
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
	api := SoulService{SoulConfig: proto.SoulConfig{PublicURL: testPubDir.Path, PrivateURL: testPrivDir.Path}}
	if r := api.Init(apiCtx); r.Err() != nil {
		t.Fatal(r.Err())
	}
}
