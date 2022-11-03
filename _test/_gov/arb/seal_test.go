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

func TestSeal(t *testing.T) {
	// base.LogVerbosely()

	// create test community
	// dir := testutil.MakeStickyTestDir()
	dir := t.TempDir()
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

	// create ballot
	ballotOut, err := arbService0.CreateBallot(ctx,
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

	// cast two votes
	voteOut, err := arbService1.Vote(ctx,
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

	// tally the results
	tallyOut, err := arbService0.Tally(ctx,
		&arb.TallyIn{BallotBranch: ballotOut.BallotBranch, BallotPath: "test_ballot"})
	if err != nil {
		t.Fatal(err)
	}

	// seal the ballot
	sealOut, err := arbService0.Seal(ctx,
		&arb.SealIn{BallotBranch: ballotOut.BallotBranch, BallotPath: "test_ballot"})
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("ballot: %v\nvote: %v\ntally: %v\nseal: %v\n",
		form.Pretty(ballotOut), form.Pretty(voteOut), form.Pretty(tallyOut), form.Pretty(sealOut))
}
