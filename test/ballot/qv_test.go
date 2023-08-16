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

func TestQV(t *testing.T) {
	ctx := testutil.NewCtx()
	cty := test.NewTestCommunity(t, ctx, 2)

	ballotName := ns.NS{"a", "b", "c"}
	choices := []string{"x", "y", "z"}

	// give voter credits
	balance.Set(ctx, cty.Gov(), cty.MemberUser(0), qv.VotingCredits, 100.0)
	balance.Set(ctx, cty.Gov(), cty.MemberUser(1), qv.VotingCredits, 100.0)

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
	fmt.Println("open: ", form.SprintJSON(openChg))

	// first round of votes
	var elections0 [2]common.Elections
	elections0[0] = common.Elections{
		common.NewElection(choices[0], 2.0),
		common.NewElection(choices[1], 3.0),
		common.NewElection(choices[2], 5.0),
	}
	elections0[1] = common.Elections{
		common.NewElection(choices[0], -6.0),
		common.NewElection(choices[1], -4.0),
		common.NewElection(choices[2], -2.0),
	}
	voteChg00 := ballot.Vote(ctx, cty.MemberOwner(0), cty.Gov(), ballotName, elections0[0])
	fmt.Println("vote 0/0: ", form.SprintJSON(voteChg00))
	voteChg01 := ballot.Vote(ctx, cty.MemberOwner(1), cty.Gov(), ballotName, elections0[1])
	fmt.Println("vote 0/1: ", form.SprintJSON(voteChg01))

	// first tally
	tallyChg0 := ballot.Tally(ctx, cty.Organizer(), ballotName)
	fmt.Println("tally 0: ", form.SprintJSON(tallyChg0))
	expScores0 := map[string]float64{
		choices[0]: -1.0352761804100827,
		choices[1]: -0.2679491924311228,
		choices[2]: 0.8218544151266947,
	}
	for k, v := range expScores0 {
		got := tallyChg0.Result.Scores[k]
		if got != v {
			t.Errorf("expecting %v, got %v", v, got)
		}
	}

	// second round of votes
	var elections1 [2]common.Elections
	elections1[0] = common.Elections{
		{VoteChoice: choices[0], VoteStrengthChange: -1.0},
		{VoteChoice: choices[1], VoteStrengthChange: -2.0},
		{VoteChoice: choices[2], VoteStrengthChange: -4.0},
	}
	elections1[1] = common.Elections{
		{VoteChoice: choices[0], VoteStrengthChange: 5.0},
		{VoteChoice: choices[1], VoteStrengthChange: 3.0},
		{VoteChoice: choices[2], VoteStrengthChange: 1.0},
	}
	voteChg10 := ballot.Vote(ctx, cty.MemberOwner(0), cty.Gov(), ballotName, elections1[0])
	fmt.Println("vote 1/0: ", form.SprintJSON(voteChg10))
	voteChg11 := ballot.Vote(ctx, cty.MemberOwner(1), cty.Gov(), ballotName, elections1[1])
	fmt.Println("vote 1/1: ", form.SprintJSON(voteChg11))

	// second tally
	tallyChg1 := ballot.Tally(ctx, cty.Organizer(), ballotName)
	fmt.Println("tally 1: ", form.SprintJSON(tallyChg1))
	expScores1 := map[string]float64{
		choices[0]: 0.0,
		choices[1]: 0.0,
		choices[2]: 0.0,
	}
	for k, v := range expScores1 {
		got := tallyChg1.Result.Scores[k]
		if got != v {
			t.Errorf("expecting %v, got %v", v, got)
		}
	}

	// close
	closeChg := ballot.Close(ctx, cty.Organizer(), ballotName, false)
	fmt.Println("close: ", form.SprintJSON(closeChg))

	// check the balances
	b0 := balance.Get(ctx, cty.Gov(), cty.MemberUser(0), qv.VotingCredits)
	b1 := balance.Get(ctx, cty.Gov(), cty.MemberUser(1), qv.VotingCredits)
	xb0 := 100.0 - 3.0
	xb1 := 100.0 - 3.0
	if b0 != xb0 {
		t.Errorf("expecting %v, got %v", xb0, b0)
	}
	if b1 != xb1 {
		t.Errorf("expecting %v, got %v", xb1, b1)
	}

	// testutil.Hang()
}
