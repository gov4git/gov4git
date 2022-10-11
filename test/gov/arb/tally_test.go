package arb

import (
	"fmt"
	"testing"

	"github.com/gov4git/gov4git/lib/base"
	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/services/gov/arb"
	"github.com/gov4git/gov4git/testutil"
)

func TestTally(t *testing.T) {
	base.LogVerbosely()

	// create test community
	dir := testutil.MakeStickyTestDir()
	// dir := t.TempDir()
	testCommunity, err := testutil.CreateTestCommunity(dir, 2)
	if err != nil {
		t.Fatal(err)
	}
	ctx := testCommunity.Background()

	// create arb services for both users
	arbService0 := arb.GovArbService{
		GovConfig:      testCommunity.CommunityGovConfig(),
		IdentityConfig: testCommunity.UserIdentityConfig(0),
	}
	arbService1 := arb.GovArbService{
		GovConfig:      testCommunity.CommunityGovConfig(),
		IdentityConfig: testCommunity.UserIdentityConfig(1),
	}

	// create poll
	pollOut, err := arbService0.Poll(ctx,
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

	// cast two votes
	voteOut, err := arbService1.Vote(ctx,
		&arb.VoteIn{
			ReferendumBranch: pollOut.PollBranch,
			VoteChoice:       "a",
			VoteStrength:     1.0,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	// tally the results
	base.Infof("tallying")
	tallyOut, err := arbService0.Tally(ctx, &arb.TallyIn{ReferendumBranch: pollOut.PollBranch})
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("poll: %v\nvote: %v\ntally: %v\n", pollOut.Human(ctx), voteOut.Human(ctx), tallyOut.Human(ctx))
}
