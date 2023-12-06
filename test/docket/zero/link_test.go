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

func TestLink(t *testing.T) {
	base.LogVerbosely()
	ctx := testutil.NewCtx(t, runtime.TestWithCache)
	cty := test.NewTestCommunity(t, ctx, 2)

	id1 := schema.MotionID("123")
	id2 := schema.MotionID("456")

	// open
	ops.OpenMotion(
		ctx,
		cty.Organizer(),
		id1,
		zero.ZeroPolicyName,
		cty.MemberUser(0),
		"concern #1",
		"description #1",
		schema.MotionConcernType,
		"https://1",
		nil,
	)
	ops.OpenMotion(
		ctx,
		cty.Organizer(),
		id2,
		zero.ZeroPolicyName,
		cty.MemberUser(1),
		"concern #2",
		"description #2",
		schema.MotionProposalType,
		"https://2",
		nil,
	)

	// link
	ops.LinkMotions(ctx, cty.Organizer(), id1, id2, "linkType")

	// verify state changed
	m1 := ops.ShowMotion(ctx, cty.Gov(), id1)
	if !m1.Motion.RefersTo(id2, "linkType") {
		t.Errorf("to ref not found")
	}
	m2 := ops.ShowMotion(ctx, cty.Gov(), id2)
	if !m2.Motion.ReferredBy(id1, "linkType") {
		t.Errorf("from ref not found")
	}

	// unlink
	ops.UnlinkMotions(ctx, cty.Organizer(), id1, id2, "linkType")

	// verify state changed
	m1 = ops.ShowMotion(ctx, cty.Gov(), id1)
	if m1.Motion.RefersTo(id2, "linkType") {
		t.Errorf("to ref still found")
	}
	m2 = ops.ShowMotion(ctx, cty.Gov(), id2)
	if m2.Motion.ReferredBy(id1, "linkType") {
		t.Errorf("from ref still found")
	}

	// testutil.Hang()
}
