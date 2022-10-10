package gov

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/gov4git/gov4git/lib/base"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/services/gov"
	"github.com/gov4git/gov4git/testutil"
)

// XXX: test taking snapshot twice in the same community repo

func TestSnapshot(t *testing.T) {
	base.LogVerbosely()

	// create test community
	dir := filepath.Join(os.TempDir(), "gov4git_test") // t.TempDir()
	testCommunity, err := testutil.CreateTestCommunity(dir, 1)
	if err != nil {
		t.Fatal(err)
	}
	ctx := testCommunity.WithWorkDir(context.Background())

	govService := gov.GovService{
		GovConfig: proto.GovConfig{CommunityURL: testCommunity.CommunityRepoURL()},
	}

	// create a local clone of the community
	communityClone, err := git.MakeLocalInCtx(ctx, "community_clone")
	if err != nil {
		t.Fatal(err)
	}
	if err := communityClone.CloneOrInitBranch(ctx, testCommunity.CommunityRepoURL(), proto.MainBranch); err != nil {
		t.Fatal(err)
	}

	// snapshot a user's public identity repo
	_, err = govService.SnapshotBranchLatest(ctx, &gov.SnapshotBranchLatestIn{
		SourceRepo:   testCommunity.UserPublicRepoURL(0),
		SourceBranch: proto.IdentityBranch,
		Community:    communityClone,
	})
	if err != nil {
		t.Fatal(err)
	}

	// snapshot again
	_, err = govService.SnapshotBranchLatest(ctx, &gov.SnapshotBranchLatestIn{
		SourceRepo:   testCommunity.UserPublicRepoURL(0),
		SourceBranch: proto.IdentityBranch,
		Community:    communityClone,
	})
	if err != nil {
		t.Fatal(err)
	}

}
