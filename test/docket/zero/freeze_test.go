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

func TestFreeze(t *testing.T) {
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

	// freeze
	ops.FreezeMotion(ctx, cty.Organizer(), id)

	// score
	ops.ScoreMotions(ctx, cty.Organizer())

	// verify state changed
	m := ops.ShowMotion(ctx, cty.Gov(), id)
	if !m.Motion.Frozen {
		t.Errorf("expecting frozen")
	}

	// unfreeze
	ops.UnfreezeMotion(ctx, cty.Organizer(), id)

	m = ops.ShowMotion(ctx, cty.Gov(), id)
	if m.Motion.Frozen {
		t.Errorf("expecting not frozen")
	}

	// testutil.Hang()
}
