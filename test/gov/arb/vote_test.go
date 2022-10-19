package arb

import (
	"fmt"
	"testing"

	"github.com/gov4git/gov4git/lib/form"
	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/govproto"
	"github.com/gov4git/gov4git/services/gov/arb"
	"github.com/gov4git/gov4git/testutil"
)

func TestVote(t *testing.T) {
	// base.LogVerbosely()

	// create test community
	// dir := testutil.MakeStickyTestDir()
	dir := t.TempDir()
	testCommunity, err := testutil.CreateTestCommunity(dir, 1)
	if err != nil {
		t.Fatal(err)
	}
	ctx := testCommunity.Background()

	// create ballot
	arbService := arb.GovArbService{
		GovConfig:      testCommunity.CommunityGovConfig(),
		IdentityConfig: testCommunity.UserIdentityConfig(0),
	}
	ballotOut, err := arbService.CreateBallot(ctx,
		&arb.CreateBallotIn{
			Path:            "test_ballot",
			Choices:         []string{"a", "b", "c"},
			Group:           "all",
			Strategy:        govproto.PriorityPollStrategyName,
			GoverningBranch: proto.MainBranch,
		})
	if err != nil {
		t.Fatal(err)
	}

	// cast a vote
	voteOut, err := arbService.Vote(ctx,
		&arb.VoteIn{
			BallotBranch: ballotOut.BallotBranch,
			BallotPath:   "test_ballot",
			Votes: []govproto.Election{
				{Choice: "a", Strength: 1.0},
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("ballot: %v\nvote: %v\n", form.Pretty(ballotOut), form.Pretty(voteOut))
}
