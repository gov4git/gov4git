package pmp

import (
	"context"
	"math"
	"testing"

	"github.com/gov4git/gov4git/proto/account"
	"github.com/gov4git/gov4git/proto/ballot/ballot"
	"github.com/gov4git/gov4git/proto/ballot/common"
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

const (
	testUser0Credits          = 101.0
	testUser1Credits          = 103.0
	testUser0ConcernStrenth   = 30.0
	testUser0ProposalStrength = 70.0
	testUser1ConcernStrength  = -20.0
	testUser1ProposalStrength = -10.0
)

func testSetup(
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
		pmp.ClaimsRefType,
	)

	// list
	ms := ops.ListMotions(ctx, cty.Gov())
	if len(ms) != 2 {
		t.Errorf("expecting 2 motions, got %v", len(ms))
	}

	// give credits to users
	account.Issue(ctx, cty.Gov(), cty.MemberAccountID(0), account.H(account.PluralAsset, testUser0Credits), "test")
	account.Issue(ctx, cty.Gov(), cty.MemberAccountID(1), account.H(account.PluralAsset, testUser1Credits), "test")

	// cast votes
	conEls := func(amt float64) common.Elections {
		return common.OneElection(pmp.ConcernBallotChoice, amt)
	}
	propEls := func(amt float64) common.Elections {
		return common.OneElection(pmp.ProposalBallotChoice, amt)
	}

	// concern votes
	ballot.Vote(ctx, cty.MemberOwner(0), cty.Gov(), pmp.ConcernPollBallotName(testConcernID), conEls(testUser0ConcernStrenth))
	ballot.Vote(ctx, cty.MemberOwner(1), cty.Gov(), pmp.ConcernPollBallotName(testConcernID), conEls(testUser1ConcernStrength))

	// proposal votes
	ballot.Vote(ctx, cty.MemberOwner(0), cty.Gov(), pmp.ProposalApprovalPollName(testProposalID), propEls(testUser0ProposalStrength))
	ballot.Vote(ctx, cty.MemberOwner(1), cty.Gov(), pmp.ProposalApprovalPollName(testProposalID), propEls(testUser1ProposalStrength))

	ballot.TallyAll(ctx, cty.Organizer(), 3)

	ops.ScoreMotions(ctx, cty.Organizer())
	ops.UpdateMotions(ctx, cty.Organizer())
}

func TestOpenCancelConcernCloseProposal(t *testing.T) {
	base.LogVerbosely()
	ctx := testutil.NewCtx(t, runtime.TestWithCache)
	cty := test.NewTestCommunity(t, ctx, 2)

	testSetup(t, ctx, cty)

	ops.CancelMotion(ctx, cty.Organizer(), testConcernID)                // issue
	ops.CloseMotion(ctx, cty.Organizer(), testProposalID, schema.Accept) // pr

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

	testSetup(t, ctx, cty)

	ops.CancelMotion(ctx, cty.Organizer(), testConcernID)        // issue
	ops.CancelMotion(ctx, cty.Organizer(), testProposalID, true) // pr

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

	testSetup(t, ctx, cty)

	ops.CancelMotion(ctx, cty.Organizer(), testProposalID, true) // pr
	ops.CancelMotion(ctx, cty.Organizer(), testConcernID)        // issue

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

	testSetup(t, ctx, cty)

	ops.CloseMotion(ctx, cty.Organizer(), testProposalID, schema.Accept) // pr

	err := must.Try(
		func() {
			ops.CancelMotion(ctx, cty.Organizer(), testConcernID)
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
