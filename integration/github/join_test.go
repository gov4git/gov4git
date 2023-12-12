//go:build integration
// +build integration

package github_test

import (
	"fmt"
	"testing"

	govgh "github.com/gov4git/gov4git/github"
	"github.com/gov4git/gov4git/runtime"
	"github.com/gov4git/gov4git/test"
	"github.com/gov4git/lib4git/base"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/testutil"
)

func TestProcessJoinRequestIssues(t *testing.T) {
	base.LogVerbosely()

	ghRepo := TestRepo
	ghClient := client

	// init governance
	ctx := testutil.NewCtx(t, runtime.TestWithCache)
	cty := test.NewTestCommunity(t, ctx, 2)

	// import issues
	chg := govgh.ProcessJoinRequestIssuesApprovedByMaintainer(ctx, ghRepo, ghClient, cty.Organizer(), false)
	fmt.Println("REPORT", form.SprintJSON(chg.Result))

	if len(chg.Result.Joined) != 0 {
		t.Errorf("expecting no joins")
	}
	if len(chg.Result.NotJoined) != 1 {
		t.Fatalf("expecting 1 non-join")
	}
	if chg.Result.NotJoined[0] != testJoinRequestAuthor {
		t.Errorf("expecting %v, got %v", testJoinRequestAuthor, chg.Result.NotJoined[0])
	}
}

var (
	testJoinRequestAuthor = "petar"
)
