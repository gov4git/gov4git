//go:build integration
// +build integration

package github_test

import (
	"testing"

	govgh "github.com/gov4git/gov4git/github"
	"github.com/gov4git/gov4git/runtime"
	"github.com/gov4git/gov4git/test"
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

	/*
		// list ballots
		ads1 := ballot.List(ctx, cty.Gov())
		if len(ads1) < 4 {
			t.Errorf("expecting at least 3, got %v,\n%v", len(ads1), form.SprintJSON(ads1))
		}
		// issue-1: open, not-frozen
		if ads1[0].Name.GitPath() != "github/issues/1" {
			t.Errorf("expecting github/issues/1, got %v", ads1[0].Name.GitPath())
		}
		if ads1[0].Closed || ads1[0].Frozen {
			t.Errorf("expecting open, not-frozen")
		}
		// issue-2: open, frozen
		if ads1[1].Name.GitPath() != "github/issues/2" {
			t.Errorf("expecting github/issues/2, got %v", ads1[1].Name.GitPath())
		}
		if ads1[1].Closed || !ads1[1].Frozen {
			t.Errorf("expecting open, frozen")
		}
		// issue-3: closed, frozen
		if ads1[2].Name.GitPath() != "github/issues/3" {
			t.Errorf("expecting github/issues/3, got %v", ads1[2].Name.GitPath())
		}
		if !ads1[2].Closed || !ads1[2].Frozen {
			t.Errorf("expecting closed, frozen")
		}
	*/
}
