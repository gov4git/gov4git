package arb

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/petar/gov4git/lib/base"
	"github.com/petar/gov4git/proto"
	"github.com/petar/gov4git/services/gov/arb"
	"github.com/petar/gov4git/testutil"
)

func TestVote(t *testing.T) {
	base.LogVerbosely()

	// create test community
	dir := filepath.Join(os.TempDir(), "gov4git_test") // t.TempDir()
	testCommunity, err := testutil.CreateTestCommunity(dir, 1)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	// create poll
	govService := arb.GovArbService{
		GovConfig:      testCommunity.CommunityGovConfig(),
		IdentityConfig: testCommunity.UserIdentityConfig(0),
	}
	pollOut, err := govService.Poll(
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
	voteOut, err := govService.Vote(
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
