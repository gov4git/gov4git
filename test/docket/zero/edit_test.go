package zero

import (
	"testing"

	"github.com/gov4git/gov4git/proto/docket/ops"
	"github.com/gov4git/gov4git/proto/docket/policies/zero"
	"github.com/gov4git/gov4git/proto/docket/schema"
	"github.com/gov4git/gov4git/runtime"
	"github.com/gov4git/gov4git/test"
	"github.com/gov4git/lib4git/base"
	"github.com/gov4git/lib4git/testutil"
)

func TestEdit(t *testing.T) {
	base.LogVerbosely()
	ctx := testutil.NewCtx(t, runtime.TestWithCache)
	cty := test.NewTestCommunity(t, ctx, 2)

	id := schema.MotionID("123")

	// open
	ops.OpenMotion(
		ctx,
		cty.Organizer(),
		id,
		schema.MotionConcernType,
		zero.ZeroPolicyName,
		cty.MemberUser(0),
		"concern #1",
		"description #1",
		"https://1",
		nil)

	// edit
	ops.EditMotion(ctx, cty.Organizer(), id, "concern #2", "description #2", "https://2", nil)

	// verify state changed
	m := ops.ShowMotion(ctx, cty.Gov(), id)
	if m.Motion.Title != "concern #2" {
		t.Errorf("expecting %v, got %v", "concern #2", m.Motion.Title)
	}
	if m.Motion.Body != "description #2" {
		t.Errorf("expecting %v, got %v", "description #2", m.Motion.Body)
	}
	if m.Motion.TrackerURL != "https://2" {
		t.Errorf("expecting %v, got %v", "https://2", m.Motion.TrackerURL)
	}

	// testutil.Hang()
}
