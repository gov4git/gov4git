package zero

import (
	"testing"

	"github.com/gov4git/gov4git/v2/proto/motion/motionapi"
	"github.com/gov4git/gov4git/v2/proto/motion/motionpolicies/zero"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
	"github.com/gov4git/gov4git/v2/runtime"
	"github.com/gov4git/gov4git/v2/test"
	"github.com/gov4git/lib4git/base"
	"github.com/gov4git/lib4git/testutil"
)

func TestEdit(t *testing.T) {
	base.LogVerbosely()
	ctx := testutil.NewCtx(t, runtime.TestWithCache)
	cty := test.NewTestCommunity(t, ctx, 2)

	id := motionproto.MotionID("123")

	// open
	motionapi.OpenMotion(
		ctx,
		cty.Organizer(),
		id,
		motionproto.MotionConcernType,
		zero.ZeroPolicyName,
		cty.MemberUser(0),
		"concern #1",
		"description #1",
		"https://1",
		nil)

	// edit
	motionapi.EditMotion(ctx, cty.Organizer(), id, cty.MemberUser(0), "concern #2", "description #2", "https://2", nil)

	// verify state changed
	m := motionapi.ShowMotion(ctx, cty.Gov(), id)
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
