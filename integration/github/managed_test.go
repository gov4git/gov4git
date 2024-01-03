//go:build integration
// +build integration

package github_test

import (
	"testing"

	govgh "github.com/gov4git/gov4git/v2/github"
	"github.com/gov4git/gov4git/v2/proto/motion/motionapi"
	"github.com/gov4git/gov4git/v2/proto/motion/motionpolicies/pmp"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
	"github.com/gov4git/gov4git/v2/runtime"
	"github.com/gov4git/gov4git/v2/test"
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
	ms := motionapi.ListMotions(ctx, cty.Gov())
	if len(ms) < 4 {
		t.Errorf("expecting at least 4 got %v,\n%v", len(ms), form.SprintJSON(ms))
	}

	// issue-1: open, not-frozen
	m, ok := ms.FindID("1")
	if !ok {
		t.Errorf("expecting 1")
	}
	if m.Closed || m.Frozen {
		t.Errorf("expecting open, not-frozen")
	}
	// refs
	if len(m.RefBy) < 1 {
		t.Errorf("expecting %v, got %v", 1, m.RefBy)
	}
	ms0RefBy0 := motionproto.Ref{Type: pmp.ClaimsRefType, From: motionproto.MotionID("5"), To: motionproto.MotionID("1")}
	if !m.RefBy.Contains(ms0RefBy0) {
		t.Errorf("expecting %v, got %v", ms0RefBy0, m.RefBy)
	}

	// issue-2: open, frozen
	m, ok = ms.FindID("2")
	if !ok {
		t.Errorf("expecting 2")
	}
	if m.Closed || m.Frozen {
		t.Errorf("expecting (got): open (%v), frozen (%v)", m.Closed, m.Frozen)
	}
	// refs
	if len(m.RefBy) < 2 {
		t.Errorf("expecting %v, got %v", 2, m.RefBy)
	}
	ms1RefBy0 := motionproto.Ref{Type: pmp.ClaimsRefType, From: motionproto.MotionID("5"), To: motionproto.MotionID("2")}
	ms1RefBy1 := motionproto.Ref{Type: pmp.ClaimsRefType, From: motionproto.MotionID("7"), To: motionproto.MotionID("2")}
	if !m.RefBy.Contains(ms1RefBy0) {
		t.Errorf("expecting %v, got %v", ms1RefBy0, m.RefBy)
	}
	if !m.RefBy.Contains(ms1RefBy1) {
		t.Errorf("expecting %v, got %v", ms1RefBy1, m.RefBy)
	}

	// issue-5: closed, frozen
	m, ok = ms.FindID("5")
	if !ok {
		t.Errorf("expecting 5")
	}
	if m.Closed || m.Frozen {
		t.Errorf("expecting open, not frozen")
	}
	// refs
	if len(m.RefTo) < 2 {
		t.Errorf("expecting %v, got %v", 2, m.RefTo)
	}
	ms2RefTo0 := motionproto.Ref{Type: pmp.ClaimsRefType, From: motionproto.MotionID("5"), To: motionproto.MotionID("1")}
	ms2RefTo1 := motionproto.Ref{Type: pmp.ClaimsRefType, From: motionproto.MotionID("5"), To: motionproto.MotionID("2")}
	if !m.RefTo.Contains(ms2RefTo0) {
		t.Errorf("expecting %v, got %v", ms2RefTo0, m.RefTo)
	}
	if !m.RefTo.Contains(ms2RefTo1) {
		t.Errorf("expecting %v, got %v", ms2RefTo1, m.RefTo)
	}

	// issue-7: closed, frozen
	m, ok = ms.FindID("7")
	if !ok {
		t.Errorf("expecting 7")
	}
	if m.Closed || m.Frozen {
		t.Errorf("expecting open, not frozen")
	}
	// refs
	if len(m.RefTo) < 1 {
		t.Errorf("expecting %v, got %v", 1, m.RefTo)
	}
	ms3RefTo0 := motionproto.Ref{Type: pmp.ClaimsRefType, From: motionproto.MotionID("7"), To: motionproto.MotionID("2")}
	if !m.RefTo.Contains(ms3RefTo0) {
		t.Errorf("expecting %v, got %v", ms3RefTo0, m.RefTo)
	}
}
