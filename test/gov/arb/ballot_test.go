package arb

import (
	"fmt"
	"testing"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/services/gov/arb"
	"github.com/gov4git/gov4git/testutil"
)

func TestBallot(t *testing.T) {
	// create test community
	// dir := testutil.MakeStickyTestDir()
	dir := t.TempDir()
	testCommunity, err := testutil.CreateTestCommunity(dir, 2)
	if err != nil {
		t.Fatal(err)
	}
	ctx := testCommunity.Background()

	// invoke service
	svc := arb.GovArbService{GovConfig: testCommunity.CommunityGovConfig()}
	out, err := svc.CreateBallot(ctx,
		&arb.CreateBallotIn{
			Path:            "test_ballot",
			Choices:         []string{"a", "b", "c"},
			Group:           "all",
			Strategy:        proto.PriorityPollStrategyName,
			GoverningBranch: proto.MainBranch,
		})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%#v\n", out)
}
