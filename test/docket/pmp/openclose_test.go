package pmp

import (
	"context"
	"testing"

	"github.com/gov4git/gov4git/proto/docket/ops"
	"github.com/gov4git/gov4git/proto/docket/policies/pmp"
	"github.com/gov4git/gov4git/proto/docket/policies/pmp/concern"
	"github.com/gov4git/gov4git/proto/docket/policies/pmp/proposal"
	"github.com/gov4git/gov4git/proto/docket/schema"
	"github.com/gov4git/gov4git/runtime"
	"github.com/gov4git/gov4git/test"
	"github.com/gov4git/lib4git/base"
	"github.com/gov4git/lib4git/must"
	"github.com/gov4git/lib4git/testutil"
)

var (
	testConcernID  = schema.MotionID("123")
	testProposalID = schema.MotionID("456")
)

func testCreateMotions(
	t *testing.T,
	ctx context.Context,
	cty *test.TestCommunity,
) {

	// open concern
	ops.OpenMotion(
		ctx,
		cty.Organizer(),
		testConcernID,
		schema.MotionConcernType,
		concern.ConcernPolicyName,
		cty.MemberUser(0),
		"concern #1",
		"body #1",
		"https://1",
		nil)

	// open proposal
	ops.OpenMotion(
		ctx,
		cty.Organizer(),
		testProposalID,
		schema.MotionProposalType,
		proposal.ProposalPolicyName,
		cty.MemberUser(1),
		"proposal #2",
		"body #2",
		"https://2",
		nil)

	// link
	ops.LinkMotions(
		ctx,
		cty.Organizer(),
		testProposalID,
		testConcernID,
		pmp.ResolvesRefType,
	)

	// list
	ms := ops.ListMotions(ctx, cty.Gov())
	if len(ms) != 2 {
		t.Errorf("expecting 2 motions, got %v", len(ms))
	}

}

func TestOpenCancelConcernCloseProposal(t *testing.T) {
	base.LogVerbosely()
	ctx := testutil.NewCtx(t, runtime.TestWithCache)
	cty := test.NewTestCommunity(t, ctx, 2)

	testCreateMotions(t, ctx, cty)

	ops.ScoreMotions(ctx, cty.Organizer())

	ops.CancelMotion(ctx, cty.Organizer(), testConcernID)       // issue
	ops.CloseMotion(ctx, cty.Organizer(), testProposalID, true) // pr
}

func TestOpenCancelConcernCancelProposal(t *testing.T) {
	base.LogVerbosely()
	ctx := testutil.NewCtx(t, runtime.TestWithCache)
	cty := test.NewTestCommunity(t, ctx, 2)

	testCreateMotions(t, ctx, cty)

	ops.ScoreMotions(ctx, cty.Organizer())

	ops.CancelMotion(ctx, cty.Organizer(), testConcernID)        // issue
	ops.CancelMotion(ctx, cty.Organizer(), testProposalID, true) // pr
}

func TestOpenCancelProposalCancelConcern(t *testing.T) {
	base.LogVerbosely()
	ctx := testutil.NewCtx(t, runtime.TestWithCache)
	cty := test.NewTestCommunity(t, ctx, 2)

	testCreateMotions(t, ctx, cty)

	ops.ScoreMotions(ctx, cty.Organizer())

	ops.CancelMotion(ctx, cty.Organizer(), testProposalID, true) // pr
	ops.CancelMotion(ctx, cty.Organizer(), testConcernID)        // issue
}

func TestOpenCloseProposalCancelConcern(t *testing.T) {
	base.LogVerbosely()
	ctx := testutil.NewCtx(t, runtime.TestWithCache)
	cty := test.NewTestCommunity(t, ctx, 2)

	testCreateMotions(t, ctx, cty)

	ops.ScoreMotions(ctx, cty.Organizer())

	ops.CloseMotion(ctx, cty.Organizer(), testProposalID, true) // pr
	err := must.Try(
		func() {
			ops.CancelMotion(ctx, cty.Organizer(), testConcernID)
		},
	) // issue
	if err == nil {
		t.Errorf("expecting error")
	}
}
