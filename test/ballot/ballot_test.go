package ballot

import (
	"fmt"
	"math"
	"testing"

	"github.com/gov4git/gov4git/v2/proto/account"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotapi"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotio"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/member"
	"github.com/gov4git/gov4git/v2/proto/purpose"
	"github.com/gov4git/gov4git/v2/runtime"
	"github.com/gov4git/gov4git/v2/test"
	"github.com/gov4git/lib4git/base"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/testutil"
)

const testMaxPar = 3

func TestOpenClose(t *testing.T) {
	base.LogVerbosely()
	ctx := testutil.NewCtx(t, runtime.TestWithCache)
	cty := test.NewTestCommunity(t, ctx, 2)

	ballotName := ballotproto.ParseBallotID("a/b/c")
	choices := []string{"x", "y", "z"}

	// open
	strat := ballotio.QVPolicyName
	ballotapi.Open(ctx, strat, cty.Organizer(), ballotName, account.NobodyAccountID, purpose.Unspecified, "", "ballot_id", "ballot description", choices, member.Everybody)

	// list
	ads := ballotapi.List(ctx, cty.Gov())
	if len(ads) != 1 {
		t.Errorf("expecting 1 ad, got %v", len(ads))
	}

	// give credits to user
	account.Issue(ctx, cty.Gov(), cty.MemberAccountID(0), account.H(account.PluralAsset, 1.0), "test")

	// vote
	elections := ballotproto.Elections{
		ballotproto.NewElection(choices[0], 1.0),
	}
	ballotapi.Vote(ctx, cty.MemberOwner(0), cty.Gov(), ballotName, elections)

	// tally
	tallyChg := ballotapi.Tally(ctx, cty.Organizer(), ballotName, testMaxPar)
	if tallyChg.Result.Scores[choices[0]] != 1.0 {
		t.Errorf("expecting %v vote, got %v", 1.0, tallyChg.Result.Scores[choices[0]])
	}

	// close
	ballotapi.Close(ctx, cty.Organizer(), ballotName, account.BurnAccountID)

	// verify state changed
	ast := ballotapi.Show(ctx, gov.Address(cty.Organizer().Public), ballotName)
	if !ast.Ad.Closed {
		t.Errorf("expecting closed flag")
	}

	// verify no credits left
	credits := account.Get(ctx, cty.Gov(), cty.MemberAccountID(0)).Balance(account.PluralAsset)
	if credits.Quantity != 0.0 {
		t.Errorf("expecting %v, got %v", 0.0, credits)
	}

	// testutil.Hang()
}

func TestOpenCancel(t *testing.T) {
	ctx := testutil.NewCtx(t, runtime.TestWithCache)
	cty := test.NewTestCommunity(t, ctx, 2)

	ballotName := ballotproto.ParseBallotID("a/b/c")
	choices := []string{"x", "y", "z"}

	// open
	strat := ballotio.QVPolicyName
	ballotapi.Open(ctx, strat, cty.Organizer(), ballotName, account.NobodyAccountID, purpose.Unspecified, "", "ballot_id", "ballot description", choices, member.Everybody)

	// list
	ads := ballotapi.List(ctx, cty.Gov())
	if len(ads) != 1 {
		t.Errorf("expecting 1 ad, got %v", len(ads))
	}

	// give credits to user
	account.Issue(ctx, cty.Gov(), cty.MemberAccountID(0), account.H(account.PluralAsset, 1.0), "test")

	// vote
	elections := ballotproto.Elections{
		ballotproto.NewElection(choices[0], 1.0),
	}
	ballotapi.Vote(ctx, cty.MemberOwner(0), cty.Gov(), ballotName, elections)

	// tally
	tallyChg := ballotapi.Tally(ctx, cty.Organizer(), ballotName, testMaxPar)
	if tallyChg.Result.Scores[choices[0]] != 1.0 {
		t.Errorf("expecting %v vote, got %v", 1.0, tallyChg.Result.Scores[choices[0]])
	}

	// close
	ballotapi.Cancel(ctx, cty.Organizer(), ballotName)

	// verify state changed
	ast := ballotapi.Show(ctx, gov.Address(cty.Organizer().Public), ballotName)
	if !ast.Ad.Closed || !ast.Ad.Cancelled {
		t.Errorf("expecting closed and cancelled")
	}

	// verify no credits left
	credits := account.Get(ctx, cty.Gov(), cty.MemberAccountID(0)).Balance(account.PluralAsset)
	if credits.Quantity != 1.0 {
		t.Errorf("expecting %v, got %v", 1.0, credits)
	}

	// testutil.Hang()
}

func TestTallyAll(t *testing.T) {
	ctx := testutil.NewCtx(t, runtime.TestWithCache)
	cty := test.NewTestCommunity(t, ctx, 2)

	ballotName0 := ballotproto.ParseBallotID("a/b/c")
	ballotName1 := ballotproto.ParseBallotID("d/e/f")
	choices := []string{"x", "y", "z"}

	// open two ballots
	strat := ballotio.QVPolicyName
	openChg0 := ballotapi.Open(ctx, strat, cty.Organizer(), ballotName0, account.NobodyAccountID, purpose.Unspecified, "", "ballot_0", "ballot 0", choices, member.Everybody)
	fmt.Println("open 0: ", form.SprintJSON(openChg0))
	openChg1 := ballotapi.Open(ctx, strat, cty.Organizer(), ballotName1, account.NobodyAccountID, purpose.Unspecified, "", "ballot_1", "ballot 1", choices, member.Everybody)
	fmt.Println("open 1: ", form.SprintJSON(openChg1))

	// give credits to users
	account.Issue(ctx, cty.Gov(), cty.MemberAccountID(0), account.H(account.PluralAsset, 5.0), "test")
	account.Issue(ctx, cty.Gov(), cty.MemberAccountID(1), account.H(account.PluralAsset, 5.0), "test")

	// vote
	elections0 := ballotproto.Elections{ballotproto.NewElection(choices[0], 5.0)}
	elections1 := ballotproto.Elections{ballotproto.NewElection(choices[0], -5.0)}
	voteChg0 := ballotapi.Vote(ctx, cty.MemberOwner(0), cty.Gov(), ballotName0, elections0)
	fmt.Println("vote 0: ", form.SprintJSON(voteChg0))
	voteChg1 := ballotapi.Vote(ctx, cty.MemberOwner(1), cty.Gov(), ballotName1, elections1)
	fmt.Println("vote 1: ", form.SprintJSON(voteChg1))

	// tally
	tallyChg := ballotapi.TallyAll(ctx, cty.Organizer(), 2)
	fmt.Println("tally: ", form.SprintJSON(tallyChg))

	// verify tallies are correct
	ast0 := ballotapi.Show(ctx, cty.Gov(), ballotName0)
	if ast0.Tally.Scores[choices[0]] != math.Sqrt(5.0) {
		t.Errorf("expecting %v, got %v", math.Sqrt(5.0), ast0.Tally.Scores[choices[0]])
	}
	ast1 := ballotapi.Show(ctx, cty.Gov(), ballotName1)
	if ast1.Tally.Scores[choices[0]] != -math.Sqrt(5.0) {
		t.Errorf("expecting %v, got %v", -math.Sqrt(5.0), ast1.Tally.Scores[choices[0]])
	}

	// testutil.Hang()
}
