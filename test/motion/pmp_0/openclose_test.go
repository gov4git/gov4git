package pmp

import (
	"math"
	"testing"

	"github.com/gov4git/gov4git/v2/proto/account"
	"github.com/gov4git/gov4git/v2/proto/motion/motionapi"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
	"github.com/gov4git/gov4git/v2/runtime"
	"github.com/gov4git/gov4git/v2/test"
	"github.com/gov4git/lib4git/base"
	"github.com/gov4git/lib4git/must"
	"github.com/gov4git/lib4git/testutil"
)

func TestOpenCancelConcernCloseProposal(t *testing.T) {
	base.LogVerbosely()
	ctx := testutil.NewCtx(t, runtime.TestWithCache)
	cty := test.NewTestCommunity(t, ctx, 2)

	SetupTest(t, ctx, cty)

	motionapi.CancelMotion(ctx, cty.Organizer(), testConcernID)                     // issue
	motionapi.CloseMotion(ctx, cty.Organizer(), testProposalID, motionproto.Accept) // pr

	// uncomment to view and adjust notices
	// conNotices := ops.LoadMotionNotices(ctx, cty.Gov(), testConcernID)
	// propNotices := ops.LoadMotionNotices(ctx, cty.Gov(), testProposalID)
	// fmt.Println("CONCERN:", form.SprintJSON(conNotices))
	// fmt.Println("PROPOSAL:", form.SprintJSON(propNotices))

	// user accounts
	u0 := account.Get(ctx, cty.Gov(), cty.MemberAccountID(0)).Balance(account.PluralAsset)
	u1 := account.Get(ctx, cty.Gov(), cty.MemberAccountID(1)).Balance(account.PluralAsset)

	exp0 := testUser0Credits + math.Abs(testUser1ProposalStrength)
	if u0.Quantity != exp0 {
		t.Errorf("expecting %v, got %v", exp0, u0.Quantity)
	}

	exp1 := testUser1Credits + testUser1ProposalStrength
	if u1.Quantity != exp1 {
		t.Errorf("expecting %v, got %v", exp1, u1.Quantity)
	}
}

func TestOpenCancelConcernCancelProposal(t *testing.T) {
	base.LogVerbosely()
	ctx := testutil.NewCtx(t, runtime.TestWithCache)
	cty := test.NewTestCommunity(t, ctx, 2)

	SetupTest(t, ctx, cty)

	motionapi.CancelMotion(ctx, cty.Organizer(), testConcernID)        // issue
	motionapi.CancelMotion(ctx, cty.Organizer(), testProposalID, true) // pr

	// user accounts
	u0 := account.Get(ctx, cty.Gov(), cty.MemberAccountID(0)).Balance(account.PluralAsset)
	u1 := account.Get(ctx, cty.Gov(), cty.MemberAccountID(1)).Balance(account.PluralAsset)

	exp0 := testUser0Credits
	if u0.Quantity != exp0 {
		t.Errorf("expecting %v, got %v", exp0, u0.Quantity)
	}

	exp1 := testUser1Credits
	if u1.Quantity != exp1 {
		t.Errorf("expecting %v, got %v", exp1, u1.Quantity)
	}
}

func TestOpenCancelProposalCancelConcern(t *testing.T) {
	base.LogVerbosely()
	ctx := testutil.NewCtx(t, runtime.TestWithCache)
	cty := test.NewTestCommunity(t, ctx, 2)

	SetupTest(t, ctx, cty)

	motionapi.CancelMotion(ctx, cty.Organizer(), testProposalID, true) // pr
	motionapi.CancelMotion(ctx, cty.Organizer(), testConcernID)        // issue

	// user accounts
	u0 := account.Get(ctx, cty.Gov(), cty.MemberAccountID(0)).Balance(account.PluralAsset)
	u1 := account.Get(ctx, cty.Gov(), cty.MemberAccountID(1)).Balance(account.PluralAsset)

	exp0 := testUser0Credits
	if u0.Quantity != exp0 {
		t.Errorf("expecting %v, got %v", exp0, u0.Quantity)
	}

	exp1 := testUser1Credits
	if u1.Quantity != exp1 {
		t.Errorf("expecting %v, got %v", exp1, u1.Quantity)
	}
}

func TestOpenCloseProposalCancelConcern(t *testing.T) {
	base.LogVerbosely()
	ctx := testutil.NewCtx(t, runtime.TestWithCache)
	cty := test.NewTestCommunity(t, ctx, 2)

	SetupTest(t, ctx, cty)

	motionapi.CloseMotion(ctx, cty.Organizer(), testProposalID, motionproto.Accept) // pr

	err := must.Try(
		func() {
			motionapi.CancelMotion(ctx, cty.Organizer(), testConcernID)
		},
	) // issue
	if err == nil {
		t.Errorf("expecting error")
	}

	// uncomment to view and adjust notices
	// conNotices := ops.LoadMotionNotices(ctx, cty.Gov(), testConcernID)
	// propNotices := ops.LoadMotionNotices(ctx, cty.Gov(), testProposalID)
	// fmt.Println("CONCERN NOTICES:", form.SprintJSON(conNotices))
	// fmt.Println("PROPOSAL NOTICES:", form.SprintJSON(propNotices))

	// uncomment to view journal entries
	// h := history.List(ctx, cty.Gov())
	// fmt.Println("HISTORY:", form.SprintJSON(h))

	// user accounts
	u0 := account.Get(ctx, cty.Gov(), cty.MemberAccountID(0)).Balance(account.PluralAsset)
	u1 := account.Get(ctx, cty.Gov(), cty.MemberAccountID(1)).Balance(account.PluralAsset)

	exp0 := testUser0Credits - math.Abs(testUser0ConcernStrenth)
	if u0.Quantity <= exp0 {
		t.Errorf("expecting more than %v, got %v", exp0, u0.Quantity)
	}

	exp1 := testUser1Credits - math.Abs(testUser1ConcernStrength) - math.Abs(testUser1ProposalStrength)
	if u1.Quantity <= exp1 {
		t.Errorf("expecting more than %v, got %v", exp1, u1.Quantity)
	}
}
