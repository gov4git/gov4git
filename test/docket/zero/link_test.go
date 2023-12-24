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

func TestLink(t *testing.T) {
	base.LogVerbosely()
	ctx := testutil.NewCtx(t, runtime.TestWithCache)
	cty := test.NewTestCommunity(t, ctx, 2)

	id1 := motionproto.MotionID("123")
	id2 := motionproto.MotionID("456")

	// open
	motionapi.OpenMotion(
		ctx,
		cty.Organizer(),
		id1,
		motionproto.MotionConcernType,
		zero.ZeroPolicyName,
		cty.MemberUser(0),
		"concern #1",
		"description #1",
		"https://1",
		nil,
	)
	motionapi.OpenMotion(
		ctx,
		cty.Organizer(),
		id2,
		motionproto.MotionProposalType,
		zero.ZeroPolicyName,
		cty.MemberUser(1),
		"concern #2",
		"description #2",
		"https://2",
		nil,
	)

	// link
	motionapi.LinkMotions(ctx, cty.Organizer(), id1, id2, "linkType")

	// verify state changed
	m1 := motionapi.ShowMotion(ctx, cty.Gov(), id1)
	if !m1.Motion.RefersTo(id2, "linkType") {
		t.Errorf("to ref not found")
	}
	m2 := motionapi.ShowMotion(ctx, cty.Gov(), id2)
	if !m2.Motion.ReferredBy(id1, "linkType") {
		t.Errorf("from ref not found")
	}

	// unlink
	motionapi.UnlinkMotions(ctx, cty.Organizer(), id1, id2, "linkType")

	// verify state changed
	m1 = motionapi.ShowMotion(ctx, cty.Gov(), id1)
	if m1.Motion.RefersTo(id2, "linkType") {
		t.Errorf("to ref still found")
	}
	m2 = motionapi.ShowMotion(ctx, cty.Gov(), id2)
	if m2.Motion.ReferredBy(id1, "linkType") {
		t.Errorf("from ref still found")
	}

	// testutil.Hang()
}
