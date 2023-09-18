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
	"github.com/gov4git/lib4git/must"
	"github.com/gov4git/lib4git/ns"
	"github.com/gov4git/lib4git/testutil"
)

func TestInsufficientCredits(t *testing.T) {
	ctx := testutil.NewCtx(t, false)
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
		"ballot title",
		"ballot description",
		choices,
		member.Everybody,
	)
	fmt.Println("open: ", form.SprintJSON(openChg))

	// give voter credits
	balance.Set(ctx, cty.Gov(), cty.MemberUser(0), qv.VotingCredits, 1.0)

	// vote
	elections := common.Elections{
		common.NewElection(choices[0], 2.0),
	}
	if err := must.Try(
		func() { ballot.Vote(ctx, cty.MemberOwner(0), cty.Gov(), ballotName, elections) },
	); err != nil {
		fmt.Println("vote rejected: ", err.Error())
	} else {
		t.Fatalf("vote must fail")
	}

	// tally
	tallyChg := ballot.Tally(ctx, cty.Organizer(), ballotName)
	fmt.Println("tally: ", form.SprintJSON(tallyChg))
	if tallyChg.Result.Scores[choices[0]] != 0.0 {
		t.Errorf("expecting %v, got %v", 0.0, tallyChg.Result.Scores[choices[0]])
	}

	// close
	closeChg := ballot.Close(ctx, cty.Organizer(), ballotName, false)
	fmt.Println("close: ", form.SprintJSON(closeChg))

	// testutil.Hang()
}
