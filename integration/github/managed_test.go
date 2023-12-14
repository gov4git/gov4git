//go:build integration
// +build integration

package github_test

import (
	"testing"

	govgh "github.com/gov4git/gov4git/github"
	"github.com/gov4git/gov4git/proto/docket/ops"
	"github.com/gov4git/gov4git/proto/docket/policies/pmp"
	"github.com/gov4git/gov4git/proto/docket/schema"
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
	syncChanges := govgh.SyncManagedIssues(ctx, ghRepo, ghClient, cty.Organizer())
	_ = syncChanges
	// fmt.Println(form.SprintJSON(syncChanges))

	// list motions
	ms := ops.ListMotions(ctx, cty.Gov())
	if len(ms) < 4 {
		t.Errorf("expecting at least 4 got %v,\n%v", len(ms), form.SprintJSON(ms))
	}
	// fmt.Println(form.SprintJSON(ms))

	// issue-1: open, not-frozen
	if ms[0].ID != "1" {
		t.Errorf("expecting 1, got %v", ms[0].ID)
	}
	if ms[0].Closed || ms[0].Frozen {
		t.Errorf("expecting open, not-frozen")
	}
	// refs
	if len(ms[0].RefBy) != 1 {
		t.Errorf("expecting %v, got %v", 1, ms[0].RefBy)
	}
	ms0RefBy0 := schema.Ref{Type: pmp.ResolvesRefType, From: schema.MotionID("5"), To: schema.MotionID("1")}
	if ms[0].RefBy[0] != ms0RefBy0 {
		t.Errorf("expecting %v, got %v", ms0RefBy0, ms[0].RefBy[0])
	}

	// issue-2: open, frozen
	if ms[1].ID != "2" {
		t.Errorf("expecting 2, got %v", ms[1].ID)
	}
	if ms[1].Closed || !ms[1].Frozen {
		t.Errorf("expecting open, frozen")
	}
	// refs
	if len(ms[1].RefBy) != 2 {
		t.Errorf("expecting %v, got %v", 2, ms[1].RefBy)
	}
	ms1RefBy0 := schema.Ref{Type: pmp.ResolvesRefType, From: schema.MotionID("5"), To: schema.MotionID("2")}
	ms1RefBy1 := schema.Ref{Type: pmp.ResolvesRefType, From: schema.MotionID("7"), To: schema.MotionID("2")}
	if ms[1].RefBy[0] != ms1RefBy0 {
		t.Errorf("expecting %v, got %v", ms1RefBy0, ms[1].RefBy[0])
	}
	if ms[1].RefBy[1] != ms1RefBy1 {
		t.Errorf("expecting %v, got %v", ms1RefBy1, ms[1].RefBy[1])
	}

	// issue-5: closed, frozen
	if ms[2].ID != "5" {
		t.Errorf("expecting 5, got %v", ms[2].ID)
	}
	if ms[2].Closed || ms[2].Frozen {
		t.Errorf("expecting open, not frozen")
	}
	// refs
	if len(ms[2].RefTo) != 2 {
		t.Errorf("expecting %v, got %v", 2, ms[2].RefTo)
	}
	ms2RefTo0 := schema.Ref{Type: pmp.ResolvesRefType, From: schema.MotionID("5"), To: schema.MotionID("1")}
	ms2RefTo1 := schema.Ref{Type: pmp.ResolvesRefType, From: schema.MotionID("5"), To: schema.MotionID("2")}
	if ms[2].RefTo[0] != ms2RefTo0 {
		t.Errorf("expecting %v, got %v", ms2RefTo0, ms[2].RefTo[0])
	}
	if ms[2].RefTo[1] != ms2RefTo1 {
		t.Errorf("expecting %v, got %v", ms2RefTo1, ms[2].RefTo[1])
	}

	// issue-7: closed, frozen
	if ms[3].ID != "7" {
		t.Errorf("expecting 7, got %v", ms[3].ID)
	}
	if ms[3].Closed || ms[3].Frozen {
		t.Errorf("expecting open, not frozen")
	}
	// refs
	if len(ms[3].RefTo) != 1 {
		t.Errorf("expecting %v, got %v", 1, ms[3].RefTo)
	}
	ms3RefTo0 := schema.Ref{Type: pmp.ResolvesRefType, From: schema.MotionID("7"), To: schema.MotionID("2")}
	if ms[3].RefTo[0] != ms3RefTo0 {
		t.Errorf("expecting %v, got %v", ms3RefTo0, ms[3].RefTo[0])
	}
}
