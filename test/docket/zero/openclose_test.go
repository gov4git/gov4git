package zero

import (
	"testing"

	"github.com/gov4git/gov4git/proto/account"
	"github.com/gov4git/gov4git/proto/docket/ops"
	"github.com/gov4git/gov4git/proto/docket/policies/zero"
	"github.com/gov4git/gov4git/proto/docket/schema"
	"github.com/gov4git/gov4git/runtime"
	"github.com/gov4git/gov4git/test"
	"github.com/gov4git/lib4git/base"
	"github.com/gov4git/lib4git/testutil"
)

func TestOpenClose(t *testing.T) {
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

	// list
	ms := ops.ListMotions(ctx, cty.Gov())
	if len(ms) != 1 {
		t.Errorf("expecting 1 motion, got %v", len(ms))
	}

	// give credits to user
	account.Deposit(ctx, cty.Gov(), cty.MemberAccountID(0), account.H(account.PluralAsset, 13.0), "test")

	// score
	ops.ScoreMotions(ctx, cty.Organizer())

	// close
	ops.CloseMotion(ctx, cty.Organizer(), id)

	// verify state changed
	m := ops.ShowMotion(ctx, cty.Gov(), id)
	if !m.Motion.Closed || m.Motion.Cancelled {
		t.Errorf("expecting closed and not cancelled")
	}

	// testutil.Hang()
}
