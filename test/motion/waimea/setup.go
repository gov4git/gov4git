package waimea

import (
	"context"
	"testing"

	"github.com/gov4git/gov4git/v2/proto/account"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotapi"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/motion/motionapi"
	"github.com/gov4git/gov4git/v2/proto/motion/motionpolicies/waimea"
	_ "github.com/gov4git/gov4git/v2/proto/motion/motionpolicies/waimea/use"
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

	// open concern
	motionapi.OpenMotion(
		ctx,
		cty.Organizer(),
		testConcernID,
		motionproto.MotionConcernType,
		waimea.ConcernPolicyName,
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
		waimea.ProposalPolicyName,
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
		waimea.ClaimsRefType,
	)

	motionapi.Pipeline(ctx, cty.Organizer())

	// list
	ms := motionapi.ListMotions(ctx, cty.Gov())
	if len(ms) != 2 {
		t.Errorf("expecting 2 motions, got %v", len(ms))
	}

	// give credits to users
	account.Issue(ctx, cty.Gov(), cty.MemberAccountID(0), account.H(account.PluralAsset, c.Voter0Credits), "test")
	account.Issue(ctx, cty.Gov(), cty.MemberAccountID(1), account.H(account.PluralAsset, c.Voter1Credits), "test")

	// cast votes
	conEls := func(amt float64) ballotproto.Elections {
		return ballotproto.OneElection(waimea.ConcernBallotChoice, amt)
	}
	propEls := func(amt float64) ballotproto.Elections {
		return ballotproto.OneElection(waimea.ProposalBallotChoice, amt)
	}

	// concern votes
	ballotapi.Vote(ctx, cty.MemberOwner(0), cty.Gov(), waimea.ConcernPollBallotName(testConcernID), conEls(c.Voter0ConcernStrength))
	ballotapi.Vote(ctx, cty.MemberOwner(1), cty.Gov(), waimea.ConcernPollBallotName(testConcernID), conEls(c.Voter1ConcernStrength))

	// proposal votes
	ballotapi.Vote(ctx, cty.MemberOwner(0), cty.Gov(), waimea.ProposalApprovalPollName(testProposalID), propEls(c.Voter0ProposalStrength))
	ballotapi.Vote(ctx, cty.MemberOwner(1), cty.Gov(), waimea.ProposalApprovalPollName(testProposalID), propEls(c.Voter1ProposalStrength))

	ballotapi.TallyAll(ctx, cty.Organizer(), 3)

	motionapi.Pipeline(ctx, cty.Organizer())

	return ctx, cty
}
