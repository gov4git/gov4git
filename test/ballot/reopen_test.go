package ballot

import (
	"math"
	"testing"

	"github.com/gov4git/gov4git/proto/account"
	"github.com/gov4git/gov4git/proto/ballot/ballot"
	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/ballot/load"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/member"
	"github.com/gov4git/gov4git/runtime"
	"github.com/gov4git/gov4git/test"
	"github.com/gov4git/lib4git/base"
	"github.com/gov4git/lib4git/testutil"
)

func TestReopen(t *testing.T) {
	base.LogVerbosely()
	ctx := testutil.NewCtx(t, runtime.TestWithCache)
	cty := test.NewTestCommunity(t, ctx, 2)

	ballotName := common.BallotName{"a", "b", "c"}
	choices := []string{"x", "y", "z"}

	// open
	strat := load.QVStrategyName
	ballot.Open(ctx, strat, cty.Organizer(), ballotName, "ballot_name", "ballot description", choices, member.Everybody)

	// give credits to user
	account.Deposit(ctx, cty.Gov(), cty.MemberAccountID(0), account.H(account.PluralAsset, 4.0))

	// vote#1
	elections := common.Elections{
		common.NewElection(choices[0], 2.0),
	}
	ballot.Vote(ctx, cty.MemberOwner(0), cty.Gov(), ballotName, elections)

	// tally#1
	tallyChg := ballot.Tally(ctx, cty.Organizer(), ballotName, testMaxPar)
	if tallyChg.Result.Scores[choices[0]] != math.Sqrt(2.0) {
		t.Errorf("expecting %v vote, got %v", math.Sqrt(2.0), tallyChg.Result.Scores[choices[0]])
	}

	// close
	ballot.Close(ctx, cty.Organizer(), ballotName, false)

	// verify state changed
	ast := ballot.Show(ctx, gov.Address(cty.Organizer().Public), ballotName)
	if !ast.Ad.Closed {
		t.Errorf("expecting closed flag")
	}

	// reopen
	ballot.Reopen(ctx, cty.Organizer(), ballotName)

	// verify state changed
	ast = ballot.Show(ctx, gov.Address(cty.Organizer().Public), ballotName)
	if ast.Ad.Closed || ast.Ad.Cancelled {
		t.Errorf("expecting not closed and not cancelled")
	}

	// vote#2
	ballot.Vote(ctx, cty.MemberOwner(0), cty.Gov(), ballotName, elections)

	// tally#2
	tallyChg = ballot.Tally(ctx, cty.Organizer(), ballotName, testMaxPar)
	if tallyChg.Result.Scores[choices[0]] != 2.0 {
		t.Errorf("expecting %v vote, got %v", 2.0, tallyChg.Result.Scores[choices[0]])
	}

	// testutil.Hang()
}
