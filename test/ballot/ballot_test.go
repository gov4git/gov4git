package ballot

import (
	"fmt"
	"testing"

	"github.com/gov4git/gov4git/proto/balance"
	"github.com/gov4git/gov4git/proto/ballot/ballot"
	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/ballot/qv"
	"github.com/gov4git/gov4git/proto/member"
	"github.com/gov4git/gov4git/test"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/ns"
	"github.com/gov4git/lib4git/testutil"
)

func TestBallotClose(t *testing.T) {
	ctx := testutil.NewCtx()
	cty := test.NewTestCommunity(t, ctx, 2)

	ballotName := ns.NS{"a", "b", "c"}
	choices := []string{"x", "y", "z"}

	// open
	strat := qv.QV{}
	openChg := ballot.Open(
		ctx,
		strat,
		cty.Gov(),
		ballotName,
		"ballot_name",
		"ballot description",
		choices,
		member.Everybody,
	)
	fmt.Println("open: ", openChg)

	// list
	ads := ballot.List(ctx, cty.Gov(), false)
	if len(ads) != 1 {
		t.Errorf("expecting 1 ad, got %v", len(ads))
	}
	fmt.Println("ads: ", form.SprintJSON(ads))

	// give credits to user
	balance.Set(ctx, cty.Gov(), cty.MemberUser(0), qv.VotingCredits, 1.0)

	// vote
	elections := common.Elections{
		common.NewElection(choices[0], 1.0),
	}
	voteChg := ballot.Vote(
		ctx,
		cty.MemberOwner(0),
		cty.Gov(),
		ballotName,
		elections,
	)
	fmt.Println("vote: ", form.SprintJSON(voteChg))

	// tally
	tallyChg := ballot.Tally(
		ctx,
		cty.Organizer(),
		ballotName,
	)
	fmt.Println("tally: ", form.SprintJSON(tallyChg))
	if tallyChg.Result.Scores[choices[0]] != 1.0 {
		t.Errorf("expecting %v vote, got %v", 1.0, tallyChg.Result.Scores[choices[0]])
	}

	// close
	closeChg := ballot.Close(ctx, cty.Organizer(), ballotName, false)
	fmt.Println("close: ", form.SprintJSON(closeChg))

	// verify no credits left
	credits := balance.Get(ctx, cty.Gov(), cty.MemberUser(0), qv.VotingCredits)
	if credits != 0.0 {
		t.Errorf("expecting %v, got %v", 0.0, credits)
	}

	// testutil.Hang()
}

func TestBallotCancel(t *testing.T) {
	ctx := testutil.NewCtx()
	cty := test.NewTestCommunity(t, ctx, 2)

	ballotName := ns.NS{"a", "b", "c"}
	choices := []string{"x", "y", "z"}

	// open
	strat := qv.QV{}
	openChg := ballot.Open(
		ctx,
		strat,
		cty.Gov(),
		ballotName,
		"ballot_name",
		"ballot description",
		choices,
		member.Everybody,
	)
	fmt.Println("open: ", openChg)

	// list
	ads := ballot.List(ctx, cty.Gov(), false)
	if len(ads) != 1 {
		t.Errorf("expecting 1 ad, got %v", len(ads))
	}
	fmt.Println("ads: ", form.SprintJSON(ads))

	// give credits to user
	balance.Set(ctx, cty.Gov(), cty.MemberUser(0), qv.VotingCredits, 1.0)

	// vote
	elections := common.Elections{
		common.NewElection(choices[0], 1.0),
	}
	voteChg := ballot.Vote(
		ctx,
		cty.MemberOwner(0),
		cty.Gov(),
		ballotName,
		elections,
	)
	fmt.Println("vote: ", form.SprintJSON(voteChg))

	// tally
	tallyChg := ballot.Tally(
		ctx,
		cty.Organizer(),
		ballotName,
	)
	fmt.Println("tally: ", form.SprintJSON(tallyChg))
	if tallyChg.Result.Scores[choices[0]] != 1.0 {
		t.Errorf("expecting %v vote, got %v", 1.0, tallyChg.Result.Scores[choices[0]])
	}

	// close
	closeChg := ballot.Close(ctx, cty.Organizer(), ballotName, true)
	fmt.Println("close: ", form.SprintJSON(closeChg))

	// verify no credits left
	credits := balance.Get(ctx, cty.Gov(), cty.MemberUser(0), qv.VotingCredits)
	if credits != 1.0 {
		t.Errorf("expecting %v, got %v", 1.0, credits)
	}

	// testutil.Hang()
}
