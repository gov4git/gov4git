package arb

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/petar/gov4git/proto"
	"github.com/petar/gov4git/services/gov/arb"
	"github.com/petar/gov4git/testutil"
)

func TestPoll(t *testing.T) {
	// create test community
	dir := filepath.Join(os.TempDir(), "gov4git_test") // t.TempDir()
	testCommunity, err := testutil.CreateTestCommunity(dir, 2)
	if err != nil {
		t.Fatal(err)
	}

	// invoke service
	svc := arb.GovArbService{GovConfig: testCommunity.CommunityGovConfig()}
	out, err := svc.Poll(
		testCommunity.WithWorkDir(context.Background(), "test"),
		&arb.GovArbPollIn{
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
