package arb

import (
	"context"
	"fmt"
	"testing"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/services/gov/arb"
	"github.com/gov4git/gov4git/testutil"
)

func TestPoll(t *testing.T) {
	// create test community
	// dir := testutil.MakeStickyTestDir()
	dir := t.TempDir()
	testCommunity, err := testutil.CreateTestCommunity(dir, 2)
	if err != nil {
		t.Fatal(err)
	}

	// invoke service
	svc := arb.GovArbService{GovConfig: testCommunity.CommunityGovConfig()}
	out, err := svc.Poll(
		testCommunity.WithWorkDir(context.Background(), "test"),
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
	fmt.Printf("%#v\n", out)
}
