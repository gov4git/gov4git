package arb

import (
	"context"
	"fmt"
	"testing"

	"github.com/gov4git/gov4git/lib/base"
	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/services/gov/arb"
	"github.com/gov4git/gov4git/testutil"
)

func TestVote(t *testing.T) {
	base.LogVerbosely()

	// create test community
	// dir := testutil.MakeStickyTestDir()
	dir := t.TempDir()
	testCommunity, err := testutil.CreateTestCommunity(dir, 1)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	// create poll
	arbService := arb.GovArbService{
		GovConfig:      testCommunity.CommunityGovConfig(),
		IdentityConfig: testCommunity.UserIdentityConfig(0),
	}
	pollOut, err := arbService.Poll(
		testCommunity.WithWorkDir(ctx, "test_poll"),
		&arb.PollIn{
			Path:            "test_poll",
			Choices:         []string{"a", "b", "c"},
			Group:           "participants",
			Strategy:        "prioritize",
			GoverningBranch: proto.MainBranch,
		})
	if err != nil {
		t.Fatal(err)
	}

	// cast a vote
	voteOut, err := arbService.Vote(
		testCommunity.WithWorkDir(ctx, "test_vote"),
		&arb.VoteIn{
			ReferendumBranch: pollOut.PollBranch,
			VoteChoice:       "a",
			VoteStrength:     1.0,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("poll: %v\nvote: %v\n", pollOut.Human(ctx), voteOut.Human(ctx))
}
