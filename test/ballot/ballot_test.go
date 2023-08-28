package ballot

import (
	"fmt"
	"math"
	"testing"

	"github.com/gov4git/gov4git/proto/balance"
	"github.com/gov4git/gov4git/proto/ballot/ballot"
	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/ballot/qv"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/member"
	"github.com/gov4git/gov4git/test"
	"github.com/gov4git/lib4git/base"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/ns"
	"github.com/gov4git/lib4git/testutil"
)

func TestOpenClose(t *testing.T) {
	base.LogVerbosely()
	ctx := testutil.NewCtx()
	cty := test.NewTestCommunity(t, ctx, 2)

	ballotName := ns.NS{"a", "b", "c"}
	choices := []string{"x", "y", "z"}

	// open
	strat := qv.QV{}
	ballot.Open(ctx, strat, cty.Gov(), ballotName, "ballot_name", "ballot description", choices, member.Everybody)

	// list
	ads := ballot.List(ctx, cty.Gov())
	if len(ads) != 1 {
		t.Errorf("expecting 1 ad, got %v", len(ads))
	}

	// give credits to user
	balance.Set(ctx, cty.Gov(), cty.MemberUser(0), qv.VotingCredits, 1.0)

	// vote
	elections := common.Elections{
		common.NewElection(choices[0], 1.0),
	}
	ballot.Vote(ctx, cty.MemberOwner(0), cty.Gov(), ballotName, elections)

	// tally
	tallyChg := ballot.Tally(
		ctx,
		cty.Organizer(),
		ballotName,
	)
	if tallyChg.Result.Scores[choices[0]] != 1.0 {
		t.Errorf("expecting %v vote, got %v", 1.0, tallyChg.Result.Scores[choices[0]])
	}

	// close
	ballot.Close(ctx, cty.Organizer(), ballotName, false)

	// verify state changed
	ast := ballot.Show(ctx, gov.GovAddress(cty.Organizer().Public), ballotName)
	if !ast.Ad.Closed {
		t.Errorf("expecting closed flag")
	}

	// verify no credits left
	credits := balance.Get(ctx, cty.Gov(), cty.MemberUser(0), qv.VotingCredits)
	if credits != 0.0 {
		t.Errorf("expecting %v, got %v", 0.0, credits)
	}

	// testutil.Hang()
}

func TestOpenCancel(t *testing.T) {
	ctx := testutil.NewCtx()
	cty := test.NewTestCommunity(t, ctx, 2)

	ballotName := ns.NS{"a", "b", "c"}
	choices := []string{"x", "y", "z"}

	// open
	strat := qv.QV{}
	ballot.Open(ctx, strat, cty.Gov(), ballotName, "ballot_name", "ballot description", choices, member.Everybody)

	// list
	ads := ballot.List(ctx, cty.Gov())
	if len(ads) != 1 {
		t.Errorf("expecting 1 ad, got %v", len(ads))
	}

	// give credits to user
	balance.Set(ctx, cty.Gov(), cty.MemberUser(0), qv.VotingCredits, 1.0)

	// vote
	elections := common.Elections{
		common.NewElection(choices[0], 1.0),
	}
	ballot.Vote(ctx, cty.MemberOwner(0), cty.Gov(), ballotName, elections)

	// tally
	tallyChg := ballot.Tally(
		ctx,
		cty.Organizer(),
		ballotName,
	)
	if tallyChg.Result.Scores[choices[0]] != 1.0 {
		t.Errorf("expecting %v vote, got %v", 1.0, tallyChg.Result.Scores[choices[0]])
	}

	// close
	ballot.Close(ctx, cty.Organizer(), ballotName, true)

	// verify state changed
	ast := ballot.Show(ctx, gov.GovAddress(cty.Organizer().Public), ballotName)
	if !ast.Ad.Closed || !ast.Ad.Cancelled {
		t.Errorf("expecting closed and cancelled")
	}

	// verify no credits left
	credits := balance.Get(ctx, cty.Gov(), cty.MemberUser(0), qv.VotingCredits)
	if credits != 1.0 {
		t.Errorf("expecting %v, got %v", 1.0, credits)
	}

	// testutil.Hang()
}

func TestTallyAll(t *testing.T) {
	ctx := testutil.NewCtx()
	cty := test.NewTestCommunity(t, ctx, 2)

	ballotName0 := ns.NS{"a", "b", "c"}
	ballotName1 := ns.NS{"d", "e", "f"}
	choices := []string{"x", "y", "z"}

	// open two ballots
	strat := qv.QV{}
	openChg0 := ballot.Open(ctx, strat, cty.Gov(), ballotName0, "ballot_0", "ballot 0", choices, member.Everybody)
	fmt.Println("open 0: ", form.SprintJSON(openChg0))
	openChg1 := ballot.Open(ctx, strat, cty.Gov(), ballotName1, "ballot_1", "ballot 1", choices, member.Everybody)
	fmt.Println("open 1: ", form.SprintJSON(openChg1))

	// give credits to users
	balance.Set(ctx, cty.Gov(), cty.MemberUser(0), qv.VotingCredits, 5.0)
	balance.Set(ctx, cty.Gov(), cty.MemberUser(1), qv.VotingCredits, 5.0)

	// vote
	elections0 := common.Elections{common.NewElection(choices[0], 5.0)}
	elections1 := common.Elections{common.NewElection(choices[0], -5.0)}
	voteChg0 := ballot.Vote(ctx, cty.MemberOwner(0), cty.Gov(), ballotName0, elections0)
	fmt.Println("vote 0: ", form.SprintJSON(voteChg0))
	voteChg1 := ballot.Vote(ctx, cty.MemberOwner(1), cty.Gov(), ballotName1, elections1)
	fmt.Println("vote 1: ", form.SprintJSON(voteChg1))

	// tally
	tallyChg := ballot.TallyAll(ctx, cty.Organizer(), 2)
	fmt.Println("tally: ", form.SprintJSON(tallyChg))

	// verify tallies are correct
	ast0 := ballot.Show(ctx, cty.Gov(), ballotName0)
	if ast0.Tally.Scores[choices[0]] != math.Sqrt(5.0) {
		t.Errorf("expecting %v, got %v", math.Sqrt(5.0), ast0.Tally.Scores[choices[0]])
	}
	ast1 := ballot.Show(ctx, cty.Gov(), ballotName1)
	if ast1.Tally.Scores[choices[0]] != -math.Sqrt(5.0) {
		t.Errorf("expecting %v, got %v", -math.Sqrt(5.0), ast1.Tally.Scores[choices[0]])
	}

	// testutil.Hang()
}
