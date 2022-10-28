package arb

import (
	"fmt"
	"testing"

	"github.com/gov4git/gov4git/lib/base"
	"github.com/gov4git/gov4git/lib/form"
	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/govproto"
	"github.com/gov4git/gov4git/services/gov/arb"
	"github.com/gov4git/gov4git/testutil"
)

func TestBallot(t *testing.T) {
	base.LogVerbosely()

	// create test community
	// dir := testutil.MakeStickyTestDir()
	dir := t.TempDir()
	testCommunity, err := testutil.CreateTestCommunity(dir, 2)
	if err != nil {
		t.Fatal(err)
	}
	ctx := testCommunity.Background()

	// create a ballot
	svc := arb.GovArbService{GovConfig: testCommunity.CommunityGovConfig()}
	createOut, err := svc.CreateBallot(ctx,
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

	// list open ballots
	listOut, err := svc.List(ctx, &arb.ListIn{BallotBranch: proto.MainBranch})
	if err != nil {
		t.Fatal(err)
	}

	if len(listOut.OpenBallots) != 1 || listOut.OpenBallots[0] != "test_ballot" {
		t.Errorf("expecting one open ballot, got %v", form.Pretty(listOut))
	}

	fmt.Printf("%#v\n", createOut)
}
