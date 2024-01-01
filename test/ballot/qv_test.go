package ballot

import (
	"fmt"
	"testing"

	"github.com/gov4git/gov4git/v2/proto/account"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotapi"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotio"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/member"
	"github.com/gov4git/gov4git/v2/proto/purpose"
	"github.com/gov4git/gov4git/v2/runtime"
	"github.com/gov4git/gov4git/v2/test"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/testutil"
)

func TestQV(t *testing.T) {
	ctx := testutil.NewCtx(t, runtime.TestWithCache)
	cty := test.NewTestCommunity(t, ctx, 2)

	ballotName := ballotproto.ParseBallotID("a/b/c")
	choices := []string{"x", "y", "z"}

	// give voter credits
	account.Issue(ctx, cty.Gov(), cty.MemberAccountID(0), account.H(account.PluralAsset, 100.0), "test")
	account.Issue(ctx, cty.Gov(), cty.MemberAccountID(1), account.H(account.PluralAsset, 100.0), "test")

	// open
	strat := ballotio.QVStrategyName
	openChg := ballotapi.Open(
		ctx,
		strat,
		cty.Organizer(),
		ballotName,
		account.NobodyAccountID,
		purpose.Unspecified,
		"",
		"ballot_id",
		"ballot description",
		choices,
		member.Everybody,
	)
	fmt.Println("open: ", form.SprintJSON(openChg))

	// first round of votes
	var elections0 [2]ballotproto.Elections
	elections0[0] = ballotproto.Elections{
		ballotproto.NewElection(choices[0], 2.0),
		ballotproto.NewElection(choices[1], 3.0),
		ballotproto.NewElection(choices[2], 5.0),
	}
	elections0[1] = ballotproto.Elections{
		ballotproto.NewElection(choices[0], -6.0),
		ballotproto.NewElection(choices[1], -4.0),
		ballotproto.NewElection(choices[2], -2.0),
	}
	voteChg00 := ballotapi.Vote(ctx, cty.MemberOwner(0), cty.Gov(), ballotName, elections0[0])
	fmt.Println("vote 0/0: ", form.SprintJSON(voteChg00))
	voteChg01 := ballotapi.Vote(ctx, cty.MemberOwner(1), cty.Gov(), ballotName, elections0[1])
	fmt.Println("vote 0/1: ", form.SprintJSON(voteChg01))

	// first tally
	tallyChg0 := ballotapi.Tally(ctx, cty.Organizer(), ballotName, testMaxPar)
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
	var elections1 [2]ballotproto.Elections
	elections1[0] = ballotproto.Elections{
		{VoteChoice: choices[0], VoteStrengthChange: -1.0},
		{VoteChoice: choices[1], VoteStrengthChange: -2.0},
		{VoteChoice: choices[2], VoteStrengthChange: -4.0},
	}
	elections1[1] = ballotproto.Elections{
		{VoteChoice: choices[0], VoteStrengthChange: 5.0},
		{VoteChoice: choices[1], VoteStrengthChange: 3.0},
		{VoteChoice: choices[2], VoteStrengthChange: 1.0},
	}
	voteChg10 := ballotapi.Vote(ctx, cty.MemberOwner(0), cty.Gov(), ballotName, elections1[0])
	fmt.Println("vote 1/0: ", form.SprintJSON(voteChg10))
	voteChg11 := ballotapi.Vote(ctx, cty.MemberOwner(1), cty.Gov(), ballotName, elections1[1])
	fmt.Println("vote 1/1: ", form.SprintJSON(voteChg11))

	// second tally
	tallyChg1 := ballotapi.Tally(ctx, cty.Organizer(), ballotName, testMaxPar)
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
	closeChg := ballotapi.Close(ctx, cty.Organizer(), ballotName, account.BurnAccountID)
	fmt.Println("close: ", form.SprintJSON(closeChg))

	// check the balances
	c0 := account.Get(ctx, cty.Gov(), cty.MemberAccountID(0)).Balance(account.PluralAsset).Quantity
	c1 := account.Get(ctx, cty.Gov(), cty.MemberAccountID(1)).Balance(account.PluralAsset).Quantity

	xb0 := 100.0 - 3.0
	xb1 := 100.0 - 3.0

	if c0 != xb0 {
		t.Errorf("expecting %v, got %v", xb0, c0)
	}
	if c1 != xb1 {
		t.Errorf("expecting %v, got %v", xb1, c1)
	}

	// testutil.Hang()
}
