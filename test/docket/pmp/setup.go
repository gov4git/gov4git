package pmp

import (
	"context"
	"testing"

	"github.com/gov4git/gov4git/v2/proto/account"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballot"
	"github.com/gov4git/gov4git/v2/proto/ballot/common"
	"github.com/gov4git/gov4git/v2/proto/docket/ops"
	"github.com/gov4git/gov4git/v2/proto/docket/policies/pmp"
	"github.com/gov4git/gov4git/v2/proto/docket/policies/pmp/concern"
	"github.com/gov4git/gov4git/v2/proto/docket/policies/pmp/proposal"
	"github.com/gov4git/gov4git/v2/proto/docket/schema"
	"github.com/gov4git/gov4git/v2/test"
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

func SetupTest(
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
