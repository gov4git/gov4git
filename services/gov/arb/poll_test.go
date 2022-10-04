package arb

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/petar/gov4git/lib/files"
	"github.com/petar/gov4git/lib/git"
	"github.com/petar/gov4git/proto"
)

func TestPoll(t *testing.T) {
	// prep test directory
	dir := filepath.Join(os.TempDir(), "gov4git_test") // t.TempDir()
	fmt.Printf("working in directory %v\n", dir)
	ctx := files.WithWorkDir(context.Background(), files.Dir{Path: filepath.Join(dir, "work")})
	if err := os.RemoveAll(dir); err != nil {
		t.Fatal(err)
	}

	// init community repo
	govRepoDir := filepath.Join(dir, "community_repo")
	govRepo := git.Local{Path: govRepoDir}
	if err := govRepo.InitBare(ctx); err != nil {
		t.Fatal(err)
	}
	// clone and make first commit
	clonedGovRepoDir := filepath.Join(dir, "community_clone")
	clonedGovRepo := git.Local{Path: clonedGovRepoDir}

	if err := clonedGovRepo.CloneOrInitBranch(ctx, govRepoDir, "main"); err != nil {
		t.Fatal(err)
	}
	if err := clonedGovRepo.Dir().WriteByteFile("empty", nil); err != nil {
		t.Fatal(err)
	}
	if err := clonedGovRepo.Add(ctx, []string{"empty"}); err != nil {
		t.Fatal(err)
	}
	if err := clonedGovRepo.Commit(ctx, "first"); err != nil {
		t.Fatal(err)
	}
	if err := clonedGovRepo.PushUpstream(ctx); err != nil {
		t.Fatal(err)
	}

	// invoke service
	svc := GovArbService{GovConfig: proto.GovConfig{CommunityURL: govRepoDir}}
	out, err := svc.Poll(ctx, &GovArbPollIn{
		Path:            "test_poll",
		Choices:         []string{"a", "b", "c"},
		Group:           "participants",
		Strategy:        "prioritize",
		GoverningBranch: "main",
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%#v\n", out)
}
