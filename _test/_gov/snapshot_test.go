package gov

import (
	"context"
	"fmt"
	"testing"

	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/govproto"
	"github.com/gov4git/gov4git/proto/idproto"
	"github.com/gov4git/gov4git/services/gov"
	"github.com/gov4git/gov4git/testutil"
)

func TestSnapshot(t *testing.T) {
	// base.LogVerbosely()

	// create test community
	// dir := testutil.MakeStickyTestDir()
	dir := t.TempDir()
	testCommunity, err := testutil.CreateTestCommunity(dir, 1)
	if err != nil {
		t.Fatal(err)
	}
	ctx := testCommunity.WithWorkDir(context.Background())

	govService := gov.GovService{
		GovConfig: govproto.GovConfig{CommunityURL: testCommunity.CommunityRepoURL()},
	}

	// create a local clone of the community
	communityClone, err := git.MakeLocalInCtx(ctx, "community_clone")
	fmt.Println("community clone in", communityClone.Path)
	if err != nil {
		t.Fatal(err)
	}
	if err := communityClone.CloneOrInitBranch(ctx, testCommunity.CommunityRepoURL(), proto.MainBranch); err != nil {
		t.Fatal(err)
	}

	// snapshot a user's public identity repo
	_, err = govService.SnapshotBranchLatest(ctx, &gov.SnapshotBranchLatestIn{
		SourceRepo:   testCommunity.UserPublicRepoURL(0),
		SourceBranch: idproto.IdentityBranch,
		Community:    communityClone,
	})
	if err != nil {
		t.Fatal(err)
	}

	// snapshot again
	_, err = govService.SnapshotBranchLatest(ctx, &gov.SnapshotBranchLatestIn{
		SourceRepo:   testCommunity.UserPublicRepoURL(0),
		SourceBranch: idproto.IdentityBranch,
		Community:    communityClone,
	})
	if err != nil {
		t.Fatal(err)
	}

}
