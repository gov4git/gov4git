package pmp

import (
	"context"
	"testing"

	"github.com/gov4git/gov4git/v2/proto/account"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotapi"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/motion/motionapi"
	"github.com/gov4git/gov4git/v2/proto/motion/motionpolicies/pmp_0"
	_ "github.com/gov4git/gov4git/v2/proto/motion/motionpolicies/pmp_0/concern"
	_ "github.com/gov4git/gov4git/v2/proto/motion/motionpolicies/pmp_0/proposal"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
	"github.com/gov4git/gov4git/v2/test"
)

var (
	testConcernID  = motionproto.MotionID("123")
	testProposalID = motionproto.MotionID("456")
)

const (
	testUser0Credits          = 101.0
	testUser1Credits          = 103.0
	testUser0ConcernStrenth   = 30.0
	testUser0ProposalStrength = 70.0
	testUser1ConcernStrength  = -20.0
	testUser1ProposalStrength = -10.0
)

func SetupTest(
	t *testing.T,
	ctx context.Context,
	cty *test.TestCommunity,
) {

	// open concern
	motionapi.OpenMotion(
		ctx,
		cty.Organizer(),
		testConcernID,
		motionproto.MotionConcernType,
		pmp_0.ConcernPolicyName,
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
		pmp_0.ProposalPolicyName,
		cty.MemberUser(1),
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
		pmp_0.ClaimsRefType,
	)

	// list
	ms := motionapi.ListMotions(ctx, cty.Gov())
	if len(ms) != 2 {
		t.Errorf("expecting 2 motions, got %v", len(ms))
	}

	// give credits to users
	account.Issue(ctx, cty.Gov(), cty.MemberAccountID(0), account.H(account.PluralAsset, testUser0Credits), "test")
	account.Issue(ctx, cty.Gov(), cty.MemberAccountID(1), account.H(account.PluralAsset, testUser1Credits), "test")

	// cast votes
	conEls := func(amt float64) ballotproto.Elections {
		return ballotproto.OneElection(pmp_0.ConcernBallotChoice, amt)
	}
	propEls := func(amt float64) ballotproto.Elections {
		return ballotproto.OneElection(pmp_0.ProposalBallotChoice, amt)
	}

	// concern votes
	ballotapi.Vote(ctx, cty.MemberOwner(0), cty.Gov(), pmp_0.ConcernPollBallotName(testConcernID), conEls(testUser0ConcernStrenth))
	ballotapi.Vote(ctx, cty.MemberOwner(1), cty.Gov(), pmp_0.ConcernPollBallotName(testConcernID), conEls(testUser1ConcernStrength))

	// proposal votes
	ballotapi.Vote(ctx, cty.MemberOwner(0), cty.Gov(), pmp_0.ProposalApprovalPollName(testProposalID), propEls(testUser0ProposalStrength))
	ballotapi.Vote(ctx, cty.MemberOwner(1), cty.Gov(), pmp_0.ProposalApprovalPollName(testProposalID), propEls(testUser1ProposalStrength))

	ballotapi.TallyAll(ctx, cty.Organizer(), 3)

	motionapi.ScoreMotions(ctx, cty.Organizer())
	motionapi.UpdateMotions(ctx, cty.Organizer())
}
