package pmp

import (
	"context"
	"testing"

	"github.com/gov4git/gov4git/v2/proto/account"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotapi"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/motion/motionapi"
	"github.com/gov4git/gov4git/v2/proto/motion/motionpolicies/pmp_0"
	"github.com/gov4git/gov4git/v2/proto/motion/motionpolicies/pmp_1"
	_ "github.com/gov4git/gov4git/v2/proto/motion/motionpolicies/pmp_1/use"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
	"github.com/gov4git/gov4git/v2/runtime"
	"github.com/gov4git/gov4git/v2/test"
	"github.com/gov4git/lib4git/base"
	"github.com/gov4git/lib4git/testutil"
)

var (
	testConcernID  = motionproto.MotionID("123")
	testProposalID = motionproto.MotionID("456")
)

func SetupTest(
	t *testing.T,
	c *testCase,

) (context.Context, *test.TestCommunity) {

	base.LogVerbosely()
	ctx := testutil.NewCtx(t, runtime.TestWithCache)
	cty := test.NewTestCommunity(t, ctx, 3) // 3 members

	account.Issue(
		ctx,
		cty.Gov(),
		pmp_0.MatchingPoolAccountID,
		account.H(account.PluralAsset, c.MatchCredits),
		"match donation",
	)

	// open concern
	motionapi.OpenMotion(
		ctx,
		cty.Organizer(),
		testConcernID,
		motionproto.MotionConcernType,
		pmp_1.ConcernPolicyName,
		cty.MemberUser(0),
		"concern #1",
		"body #1",
		"https://1",
		nil)

	// open proposal
	motionapi.OpenMotion(
		ctx,
		cty.Organizer(),
		testProposalID,
		motionproto.MotionProposalType,
		pmp_1.ProposalPolicyName,
		cty.MemberUser(2),
		"proposal #2",
		"body #2",
		"https://2",
		nil)

	// link
	motionapi.LinkMotions(
		ctx,
		cty.Organizer(),
		testProposalID,
		testConcernID,
		pmp_1.ClaimsRefType,
	)

	motionapi.Pipeline(ctx, cty.Organizer())

	// list
	ms := motionapi.ListMotions(ctx, cty.Gov())
	if len(ms) != 2 {
		t.Errorf("expecting 2 motions, got %v", len(ms))
	}

	// give credits to users
	account.Issue(ctx, cty.Gov(), cty.MemberAccountID(0), account.H(account.PluralAsset, c.User0Credits), "test")
	account.Issue(ctx, cty.Gov(), cty.MemberAccountID(1), account.H(account.PluralAsset, c.User1Credits), "test")

	// cast votes
	conEls := func(amt float64) ballotproto.Elections {
		return ballotproto.OneElection(pmp_1.ConcernBallotChoice, amt)
	}
	propEls := func(amt float64) ballotproto.Elections {
		return ballotproto.OneElection(pmp_1.ProposalBallotChoice, amt)
	}

	// concern votes
	ballotapi.Vote(ctx, cty.MemberOwner(0), cty.Gov(), pmp_1.ConcernPollBallotName(testConcernID), conEls(c.User0ConcernStrength))
	ballotapi.Vote(ctx, cty.MemberOwner(1), cty.Gov(), pmp_1.ConcernPollBallotName(testConcernID), conEls(c.User1ConcernStrength))

	// proposal votes
	ballotapi.Vote(ctx, cty.MemberOwner(0), cty.Gov(), pmp_1.ProposalApprovalPollName(testProposalID), propEls(c.User0ProposalStrength))
	ballotapi.Vote(ctx, cty.MemberOwner(1), cty.Gov(), pmp_1.ProposalApprovalPollName(testProposalID), propEls(c.User1ProposalStrength))

	ballotapi.TallyAll(ctx, cty.Organizer(), 3)

	motionapi.Pipeline(ctx, cty.Organizer())

	return ctx, cty
}
