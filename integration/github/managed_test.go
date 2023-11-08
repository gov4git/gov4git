//go:build integration
// +build integration

package github_test

import (
	"testing"

	govgh "github.com/gov4git/gov4git/github"
	"github.com/gov4git/gov4git/proto/docket/ops"
	"github.com/gov4git/gov4git/runtime"
	"github.com/gov4git/gov4git/test"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/testutil"
)

func TestSyncManagedIssues(t *testing.T) {

	ghRepo := TestRepo
	ghClient := client

	// init governance
	ctx := testutil.NewCtx(t, runtime.TestWithCache)
	cty := test.NewTestCommunity(t, ctx, 2)

	// import issues
	govgh.SyncManagedIssues(ctx, ghRepo, ghClient, cty.Organizer())

	// list motions
	ms := ops.ListMotions(ctx, cty.Gov())
	if len(ms) < 4 {
		t.Errorf("expecting at least 4, got %v,\n%v", len(ms), form.SprintJSON(ms))
	}

	// issue-1: open, not-frozen
	if ms[0].ID != "1" {
		t.Errorf("expecting 1, got %v", ms[0].ID)
	}
	if ms[0].Closed || ms[0].Frozen {
		t.Errorf("expecting open, not-frozen")
	}

	// issue-2: open, frozen
	if ms[1].ID != "2" {
		t.Errorf("expecting 2, got %v", ms[1].ID)
	}
	if ms[1].Closed || !ms[1].Frozen {
		t.Errorf("expecting open, frozen")
	}

	// issue-3: closed, frozen
	if ms[2].ID != "3" {
		t.Errorf("expecting 3, got %v", ms[2].ID)
	}
	if !ms[2].Closed || !ms[2].Frozen {
		t.Errorf("expecting closed, frozen")
	}
}
